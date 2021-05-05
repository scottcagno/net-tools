package db

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
	"unsafe"
)

const blockSZ = 16

type DataFile struct {
	file *os.File
	rw   *bufio.ReadWriter
}

func OpenDataFile(path string) (*DataFile, error) {
	file, err := openFile(path)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReaderSize(file, blockSZ)
	w := bufio.NewWriterSize(file, blockSZ)
	return &DataFile{
		file: file,
		rw:   bufio.NewReadWriter(r, w),
	}, nil
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

func (df *DataFile) Truncate(size int64) error {
	err := df.rw.Flush()
	if err != nil {
		return err
	}
	err = df.file.Sync()
	if err != nil {
		return err
	}
	df.rw = nil
	err = df.file.Truncate(size)
	if err != nil {
		return err
	}
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

func openFile(path string) (*os.File, error) {
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
	fd, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND, os.ModeSticky)
	if err != nil {
		return nil, err
	}
	return fd, nil
}
