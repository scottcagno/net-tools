package db

import (
	"os"
	"path/filepath"
)

type DataReader interface {
	Buffered() int
	Size() int
	Discard(n int) (discarded int, err error)
	Peek(n int) ([]byte, error)

	Read(p []byte) (n int, err error)
	ReadAt(p []byte, off int64) (n int, err error)

	ReadByte() (byte, error)
	ReadBytes(delim byte) ([]byte, error)

	ReadLine() (line []byte, isPrefix bool, err error)
	ReadRune() (r rune, size int, err error)
	ReadSlice(delim byte) (line []byte, err error)
	ReadString(delim byte) (string, error)

	UnreadByte() error
	UnreadRune() error
}

type DataWriter interface {
	Available() int
	Buffered() int
	Flush() error
	Size() int

	Write(p []byte) (nn int, err error)
	WriteAt(p []byte, off int64) (n int, err error)
	WriteByte(c byte) error

	WriteString(s string) (int, error)
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
