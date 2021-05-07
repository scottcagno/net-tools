package data

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	defaultMinBufferSize = 16
	defaultBufferSize    = 4096
	nullByte             = byte(0x00)
	maxUint16            = uint16(65535)                // 65535 bytes is ~64KB
	maxUint32            = uint32(4294967295)           // 4294967295 bytes is ~4096MB
	maxUint64            = uint64(18446744073709551615) // 18446744073709551615 bytes is ~16GB
)

type Store struct {
	sync.Mutex
	fd *os.File
	r  *DataReader
	w  *DataWriter
}

func OpenStore(path string) (*Store, error) {
	fd, err := OpenFile(path)
	if err != nil {
		return nil, err
	}
	return &Store{
		fd: fd,
		r:  NewDataReaderSize(fd, defaultBufferSize),
		w:  NewDataWriterSize(fd, defaultBufferSize),
	}, nil
}

func (s *Store) WriteData(b []byte) error {
	s.Lock()
	defer s.Unlock()
	return s.w.WriteBinary(b)
}

func (s *Store) ReadData() ([]byte, error) {
	s.Lock()
	defer s.Unlock()
	return s.r.ReadBinary()
}

func (s *Store) DeleteData() ([]byte, error) {
	s.Lock()
	defer s.Unlock()
	// peek at the header
	n, err := s.r.PeekUint64()
	if err != nil {
		return nil, err
	}
	// peek at the data
	b, err := s.r.br.Peek(int(n))
	if err != nil {
		return nil, err
	}
	// write over the record with an empty one
	rec := make([]byte, n, n)
	log.Printf("\nold (%d bytes): %s\nnew (%d bytes): %s\n", len(b), b, len(rec), rec)
	err = s.w.WriteBinary(rec)
	if err != nil {
		return nil, err
	}
	err = s.w.bw.Flush()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Store) SkipData() error {
	s.Lock()
	defer s.Unlock()
	n, err := s.r.ReadUint64()
	if err != nil {
		return err
	}
	_, err = s.r.br.Discard(int(n))
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetAllEntries() ([]int64, error) {
	s.Lock()
	defer s.Unlock()
	// go back to the beginning
	if _, err := s.seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	var offset int64
	var entries []int64
	for {
		n, err := s.r.ReadUint64()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		d, err := s.r.br.Discard(int(n))
		if err != nil {
			return nil, err
		}
		entries = append(entries, offset)
		offset += int64(8 + d)
	}
	return entries, nil
}

func (s *Store) GetEntry(n int) ([]byte, error) {
	s.Lock()
	defer s.Unlock()
	// go back to the beginning
	if _, err := s.seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	var entry int64
	for entry != int64(n) {
		n, err := s.r.ReadUint64()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		_, err = s.r.br.Discard(int(n))
		if err != nil {
			return nil, err
		}
		entry++
	}
	return s.r.ReadBinary()
}

func (s *Store) seek(offset int64, whence int) (int64, error) {
	off, err := s.fd.Seek(offset, whence)
	if err != nil {
		return off, err
	}
	s.r = NewDataReaderSize(s.fd, defaultBufferSize)
	return off, err
}

func (s *Store) GetEntryOffset(n int) (int64, error) {
	s.Lock()
	defer s.Unlock()
	// go back to the beginning
	if _, err := s.seek(0, io.SeekStart); err != nil {
		return -1, err
	}
	var entry int64
	var offset int64
	for entry != int64(n) {
		n, err := s.r.ReadUint64()
		if err == io.EOF {
			break
		}
		if err != nil {
			return -1, err
		}
		d, err := s.r.br.Discard(int(n))
		if err != nil {
			return -1, err
		}
		entry++
		offset += int64(8 + d)
	}
	return offset, nil
}

func (s *Store) Close() error {
	s.Lock()
	defer s.Unlock()
	err := s.w.bw.Flush()
	if err != nil {
		return err
	}
	err = s.fd.Sync()
	if err != nil {
		return err
	}
	err = s.fd.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) String() string {
	fi, err := s.fd.Stat()
	if err != nil {
		log.Panic(err)
	}
	var sb strings.Builder
	timef := "01/02/2006 " + time.Kitchen
	sb.WriteString(fmt.Sprintf("===[ %s ]===\n", time.Now().Format(timef)))
	sb.WriteString("\nFile Info\n==========\n")
	sb.WriteString(fmt.Sprintf("Name: %s\n", fi.Name()))
	sb.WriteString(fmt.Sprintf("Size: %d\n", fi.Size()))
	sb.WriteString(fmt.Sprintf("Modified: %s\n", fi.ModTime().Format(timef)))
	sb.WriteString("\nBuffer Info\n==========\n")
	sb.WriteString(fmt.Sprintf("Size: %d\n", s.w.Size()))
	sb.WriteString(fmt.Sprintf("Buffered: %d\n", s.w.Buffered()))
	sb.WriteString(fmt.Sprintf("Available: %d\n", s.w.Available()))
	return sb.String()
}

func OpenFile(path string) (*os.File, error) {
	var fd *os.File
	var err error
	if _, err = os.Stat(path); os.IsNotExist(err) {
		dir, file := filepath.Split(path)
		err = os.MkdirAll(dir, os.ModeDir)
		if err != nil {
			return nil, err
		}
		fd, err = os.Create(dir + file)
		if err != nil {
			return nil, err
		}
		err = fd.Close()
		if err != nil {
			return fd, err
		}
	}
	// removing os.O_APPEND for now
	fd, err = os.OpenFile(path, os.O_RDWR, os.ModeSticky)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

func hasRoom(data []byte) bool {
	for i := 0; i < len(data); i++ {
		if data[i] != nullByte {
			return false
		}
	}
	return true
}

func (s *Store) grow() error {
	// grab file name and size
	info, err := s.fd.Stat()
	if err != nil {
		return err
	}
	path := info.Name()
	size := info.Size()
	// close file
	err = s.fd.Close()
	if err != nil {
		return err
	}
	// truncate file
	err = os.Truncate(path, size*2)
	if err != nil {
		return err
	}
	// re-open file
	s.fd, err = OpenFile(path)
	if err != nil {
		return err
	}
	// re-initialize buffers
	s.r = NewDataReaderSize(s.fd, defaultBufferSize)
	s.w = NewDataWriterSize(s.fd, defaultBufferSize)
	return nil
}
