package db

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"unsafe"
)

const blockSZ = 64
const headerSZ = int(unsafe.Sizeof(RecordHeader{}))

func align(n int) uint {
	return uint(((n + headerSZ) + blockSZ - 1) &^ (blockSZ - 1))
}

type RecordHeader struct {
	Status uint16 // 0xC1
	Blocks uint16
	Length uint32
	Time   uint32
}

type Record struct {
	RecordHeader
	Data []byte
}

type DataFile struct {
	file *os.File
	rw   *bufio.ReadWriter
}

func OpenDataFile(path string) (*DataFile, error) {
	file, err := OpenFile(path)
	if err != nil {
		return nil, err
	}
	if blockSZ < 64 {
		return nil, fmt.Errorf("block size cannot be smaller than 64 bytes!\n")
	}
	r := bufio.NewReaderSize(file, blockSZ)
	w := bufio.NewWriterSize(file, blockSZ)
	return &DataFile{
		file: file,
		rw:   bufio.NewReadWriter(r, w),
	}, nil
}

func (df *DataFile) Seek(offset int64, whence int) (int64, error) {
	return df.file.Seek(offset, whence)
}

func (df *DataFile) SeekBlock(block int64, whence int) (int64, error) {
	offset := block * blockSZ
	return df.file.Seek(offset, whence)
}

func (df *DataFile) ReadRecord(r *Record) error {
	// do something
	return nil
}

func (df *DataFile) WriteRecord(b []byte) error {
	// do something
	_ = &Record{
		RecordHeader: RecordHeader{
			Status: uint16(0),
			Blocks: uint16(0),
			Length: uint32(0),
			Time:   uint32(0),
		},
	}
	return nil
}

func (df *DataFile) Stat() (os.FileInfo, error) {
	return df.file.Stat()
}

func (df *DataFile) Read(b []byte) (int, error) {
	return df.rw.Read(b)
}

func (df *DataFile) ReadByte() (byte, error) {
	return df.rw.ReadByte()
}

func (df *DataFile) ReadBytes() ([]byte, error) {
	n, err := df.ReadUvarint()
	if err != nil {
		return nil, err
	}
	b := make([]byte, n)
	if _, err := io.ReadFull(df.rw, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (df *DataFile) ReadString() (string, error) {
	b, err := df.ReadBytes()
	if err != nil {
		return "", err
	}
	return *(*string)(unsafe.Pointer(&b)), nil
}

func (df *DataFile) ReadUint16() (uint16, error) {
	var buf [2]byte
	if _, err := io.ReadFull(df.rw, buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(buf[:]), nil
}

func (df *DataFile) ReadUint32() (uint32, error) {
	var buf [4]byte
	if _, err := io.ReadFull(df.rw, buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf[:]), nil
}

func (df *DataFile) ReadUint64() (uint64, error) {
	var buf [8]byte
	if _, err := io.ReadFull(df.rw, buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf[:]), nil
}

func (df *DataFile) ReadUvarint() (uint64, error) {
	return binary.ReadUvarint(df.rw)
}

func (df *DataFile) ReaderBuffered() int {
	return df.rw.Reader.Buffered()
}

func (df *DataFile) ReaderSize() int {
	return df.rw.Reader.Size()
}

func (df *DataFile) UnreadByte() error {
	return df.rw.UnreadByte()
}

func (df *DataFile) Write(b []byte) (int, error) {
	return df.rw.Write(b)
}

func (df *DataFile) WriteByte(b byte) error {
	return df.rw.WriteByte(b)
}

func (df *DataFile) WriteBytes(b []byte) error {
	return df.WriteString(*(*string)(unsafe.Pointer(&b)))
}

func (df *DataFile) WriteString(s string) error {
	if err := df.WriteUvarint(uint64(len(s))); err != nil {
		return err
	}
	_, err := df.rw.WriteString(s)
	return err
}

func (df *DataFile) WriteData(b []byte) error {
	// check current position
	pos, err := df.file.Seek(0, io.SeekCurrent)
	debug("current position: %d, pos(mod)blocksize: %d", pos, pos%blockSZ)
	if pos%blockSZ != 0 {
		return fmt.Errorf("error: bad offset! (%d)\n", pos)
	}
	// block align data, len of b + 6 bytes for "header"
	bc := align(len(b)+6) / blockSZ
	debug("writing %d blocks (%d bytes)", bc, bc*blockSZ)
	// write block count (bc)
	err = df.WriteUint16(uint16(bc))
	if err != nil {
		return err
	}
	// write data length (dl)
	err = df.WriteUint32(uint32(len(b)))
	if err != nil {
		return err
	}
	debug("writing %d bytes of data", len(b))
	// write data
	_, err = df.rw.Write(b)
	if err != nil {
		return err
	}
	pad := bc*blockSZ - uint(len(b))
	debug("writing %d bytes of padding", pad)
	// write leftover, aka padding
	_, err = df.rw.Write(make([]byte, pad))
	if err != nil {
		return err
	}
	return err
}

func (df *DataFile) ReadData() ([]byte, error) {
	// check current position
	pos, err := df.file.Seek(0, io.SeekCurrent)
	debug("current position: %d, pos(mod)blocksize: %d", pos, pos%blockSZ)
	if pos%blockSZ != 0 {
		return nil, fmt.Errorf("error: bad offset! (%d)\n", pos)
	}
	// read block count (bc) and calculate next offset
	bc, err := df.ReadUint16()
	if err != nil {
		return nil, err
	}
	debug("reading %d blocks (%d bytes)", bc, bc*blockSZ)

	// skip to next block if block count is zero and return
	if bc == 0 {
		debug("block count is %d at offset %d, skipping!", bc, pos)
		_, err = df.rw.Discard(blockSZ - 2)
		return nil, err
	}

	// read data length (dl)
	dl, err := df.ReadUint32()
	if err != nil {
		return nil, err
	}
	debug("reading %d bytes of data", dl)
	// make buffer and read data off disk
	b := make([]byte, dl)
	_, err = io.ReadFull(df.rw, b)
	if err != nil {
		return nil, err
	}
	// calculate padding to the next record offset,
	// making sure to seek past to the next record
	pad := int(uint32(bc)*blockSZ - dl)
	debug("reading %d bytes of padding, seek %d more bytes (%d)\n", pad, pad, pad+int(dl))
	//_, err = df.file.Seek(pad, io.SeekStart)
	_, err = df.rw.Discard(pad)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (df *DataFile) WriteUint16(x uint16) error {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], x)
	_, err := df.rw.Write(buf[:])
	return err
}

func (df *DataFile) WriteUint32(x uint32) error {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], x)
	_, err := df.rw.Write(buf[:])
	return err
}

func (df *DataFile) WriteUint64(x uint64) error {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], x)
	_, err := df.rw.Write(buf[:])
	return err
}

func (df *DataFile) WriteUvarint(x uint64) error {
	var buf [10]byte
	n := binary.PutUvarint(buf[:], x)
	_, err := df.rw.Write(buf[:n])
	return err
}

func (df *DataFile) Available() int {
	return df.rw.Available()
}

func (df *DataFile) WriterBuffered() int {
	return df.rw.Writer.Buffered()
}

func (df *DataFile) WriterSize() int {
	return df.rw.Writer.Size()
}

func (df *DataFile) Peek(n int) ([]byte, error) {
	return df.rw.Peek(n)
}

func (df *DataFile) Discard(n int) (int, error) {
	return df.rw.Discard(n)
}

func (df *DataFile) Flush() error {
	return df.rw.Flush()
}

func (df *DataFile) Sync() error {
	err := df.rw.Flush()
	if err != nil {
		return err
	}
	return df.file.Sync()
}

func debug(format string, v ...interface{}) {
	log.Printf("[DEBUG] >> %s\n", fmt.Sprintf(format, v...))
}

func (df *DataFile) Truncate(size int64) error {
	// flush buffer
	err := df.rw.Flush()
	if err != nil {
		return err
	}
	// commit drive
	err = df.file.Sync()
	if err != nil {
		return err
	}
	// grab file name for later
	path := df.file.Name()
	// unset buffer and close file
	df.rw = nil
	err = df.file.Close()
	if err != nil {
		return err
	}
	// truncate file
	err = os.Truncate(path, size)
	if err != nil {
		return err
	}
	// re-open file
	df.file, err = OpenFile(path)
	if err != nil {
		return err
	}
	// re-initialize buffers
	r := bufio.NewReaderSize(df.file, blockSZ)
	w := bufio.NewWriterSize(df.file, blockSZ)
	df.rw = bufio.NewReadWriter(r, w)

	return nil
}

func (df *DataFile) Close() error {
	err := df.rw.Flush()
	if err != nil {
		return err
	}
	err = df.file.Sync()
	if err != nil {
		return err
	}
	err = df.file.Close()
	if err != nil {
		return err
	}
	return nil
}
