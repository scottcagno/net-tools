package ngin

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"unsafe"
)

// https://play.golang.org/p/4Dc1YH4A-AX

const (
	KB = 1 << 10 // 1 KB
	MB = 1 << 20 // 1 MB
	GB = 1 << 30 // 1 GB
)

var (
	KB_STR = fmt.Sprintf("kb: %db\n", KB)
	MB_STR = fmt.Sprintf("mb: %dkb, %db\n", MB/KB, MB)
	GB_STR = fmt.Sprintf("gb: %dmb, %dkb, %db\n", GB/MB, GB/KB, GB)
)

const (
	blockSZ = 16
	metaSZ  = int(unsafe.Sizeof(metaData{}))
)

func align(n int) uint64 {
	return uint64(((n + metaSZ) + blockSZ - 1) &^ (blockSZ - 1))
}

type metaData struct {
	header      uint16
	blockSize   uint16
	recordCount uint32
}

func (m *metaData) decode(r io.Reader) error {
	var err error
	dec := NewBinaryReader(r)
	m.header, err = dec.ReadUint16()
	if err != nil {
		return err
	}
	m.blockSize, err = dec.ReadUint16()
	if err != nil {
		return err
	}
	m.recordCount, err = dec.ReadUint32()
	if err != nil {
		return err
	}
	return nil
}

func (m *metaData) encode(w io.Writer) error {
	var err error
	enc := NewBinaryWriter(w)
	err = enc.WriteUint16(m.header)
	if err != nil {
		return err
	}
	err = enc.WriteUint16(m.blockSize)
	if err != nil {
		return err
	}
	err = enc.WriteUint32(m.recordCount)
	if err != nil {
		return err
	}
	err = enc.Flush()
	if err != nil {
		return err
	}
	return nil
}

type Engine struct {
	*metaData
	sync.Mutex
	data   *os.File
	cursor uint64
}

func (e *Engine) readMeta() error {
	e.Lock()
	defer e.Unlock()
	// seek back to start of datafile
	_, err := e.data.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	// create a new meta, decode into it
	meta := new(metaData)
	err = meta.decode(e.data)
	if err != nil {
		return err
	}
	// unset and set the new meta
	e.metaData = nil
	e.metaData = meta
	// seek back to where we were before
	_, err = e.data.Seek(int64(e.cursor), io.SeekStart)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) writeMeta() error {
	if e.metaData == nil {
		return fmt.Errorf("error: metafile is empty!\n")
	}
	e.Lock()
	defer e.Unlock()
	// seek back to start of datafile
	_, err := e.data.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	// write contents of metafile to disk
	err = e.metaData.encode(e.data)
	if err != nil {
		return err
	}
	// seek back to where we were before
	_, err = e.data.Seek(int64(e.cursor), io.SeekStart)
	if err != nil {
		return err
	}
	return nil
}

func OpenEngine(path string) (*Engine, error) {
	path += `.db`
	// check to see if we need to create a new file
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		// sanitize the filepath
		dirs, _ := filepath.Split(path)
		// create any directories
		if err := os.MkdirAll(dirs, os.ModeDir); err != nil {
			return nil, err
		}
		// create the new file
		fd, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		// init metadata
		meta := &metaData{
			header:      uint16(0),
			blockSize:   uint16(blockSZ),
			recordCount: uint32(0),
		}
		err = meta.encode(fd)
		if err != nil {
			return nil, err
		}
		// initally size it to 4MB
		if err = fd.Truncate(4 * MB); err != nil {
			return nil, err
		}
		// close the file
		if err = fd.Close(); err != nil {
			return nil, err
		}
	}
	// already existing data file, open
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, os.ModeSticky)
	if err != nil {
		return nil, err
	}
	// decode metadata off disk
	meta := new(metaData)
	err = meta.decode(fd)
	if err != nil {
		return nil, err
	}
	// seek ahead one block and set offset
	offset := align(0)
	_, err = fd.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return nil, err
	}
	// init and return engine
	en := &Engine{
		metaData: meta,
		data:     fd,
		cursor:   offset,
	}
	return en, nil
}

// findNEmptyBlocks finds n count of consecutive empty blocks and returns
// the offset and if it was successful in finding n consecutive blocks or not
func (e *Engine) findNEmptyBlocks(n uint64) (int64, error) {
	nblock, position := uint64(0), uint16(1)
	offset := int64(position * e.metaData.blockSize)
	rd := bufio.NewReader(e.data)
	for {
		_, err := e.data.Seek(offset, io.SeekStart)
		if err != nil {
			return -1, err
		}
		b, err := rd.ReadByte()
		if err != nil {
			return -1, err
		}
		if b == 0x00 {
			nblock++
		}
		err = rd.UnreadByte()
		if nblock == n {
			break
		}
		offset = int64(position * e.metaData.blockSize)
	}
	// seek back to where we were before
	_, err := e.data.Seek(int64(e.cursor), io.SeekStart)
	if err != nil {
		return -1, err
	}
	return offset / int64(e.metaData.blockSize), nil
}

// Write writes a data record to the engine at the first available slot.
// It returns a non-nil error if there is an issue growing the file size.
func (e *Engine) Write(b []byte) (int, error) {
	blocks := align(len(b))
	pos, err := e.findNEmptyBlocks(blocks)
	if err == io.EOF {
		err = e.grow()
		if err != nil {
			return -1, err
		}
	} else if err != nil {
		return -1, err
	}
	e.write(pos*int64(e.metaData.blockSize), blocks, b)
	return int(pos), nil
}

func (e *Engine) write(offset int64, blocks uint64, data []byte) {
	e.Lock()
	defer e.Unlock()
	// write record
}

func (e *Engine) grow() error {
	return nil
}

func (e *Engine) Close() error {
	err := e.writeMeta()
	if err != nil {
		return err
	}
	err = e.data.Sync()
	if err != nil {
		return err
	}
	err = e.data.Close()
	if err != nil {
		return err
	}
	e.metaData, e.cursor = nil, 0
	runtime.GC()
	return nil
}
