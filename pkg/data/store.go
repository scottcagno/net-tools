package data

import (
	"os"
	"path/filepath"
	"sync"
)

const (
	defaultMinBufferSize = 16
	defaultBufferSize    = 4096
	nullByte             = byte(0x00)
	maxKeyLen            = 65535                // uint16 max (65535 bytes is ~64KB)
	maxValueLen          = 4294967295           // uint32 max (4294967295 bytes is ~4096MB)
	maxTotalLen          = 18446744073709551615 // uint64 max (18446744073709551615 bytes is ~16GB)
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

func (s *Store) NextData() error {
	s.Lock()
	defer s.Unlock()
	n, err := s.r.ReadUint64()
	if err != nil {
		return err
	}
	_, err = s.r.Discard(int(n))
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Close() error {
	s.Lock()
	defer s.Unlock()
	err := s.fd.Sync()
	if err != nil {
		return err
	}
	err = s.fd.Close()
	if err != nil {
		return err
	}
	return nil
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
