package engine

import (
	"fmt"
	"log"
	"os"
	"path"
)

type Engine struct {
	file *os.File
}

func Open(filepath string) *Engine {
	f, err := createOrOpen(filepath)
	if err != nil {
		log.Fatalf("encountered error: %v\n", err)
	}
	return &Engine{
		file: f,
	}
}

func createOrOpen(filepath string) (*os.File, error) {
	// split path and file from filepath
	dir, file := path.Split(filepath)
	// if there are any nested directories...
	if dir != "" {
		// check to see if directories already exist...
		if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
			// if not, make all nested directories...
			if err := os.MkdirAll(dir, 0644); err != nil {
				return nil, fmt.Errorf("error creating directories: %v\n", err)
			}
		}
		// case: directory exists, continue
	}
	// if the file param isn't empty...
	if file != "" {
		// check to see if file already exists...
		if _, err := os.Stat(filepath); err != nil && os.IsNotExist(err) {
			// if not, create file...
			fd, err := os.Create(filepath)
			if err != nil {
				return nil, fmt.Errorf("error creating file: %v\n", err)
			}
			// 1<<21 is 2MB
			if err := grow(fd); err != nil {
				return nil, fmt.Errorf("error truncating file: %v\n", err)
			}
			if err = fd.Close(); err != nil {
				return nil, err
			}
		}
		// case: file exists, continue
	}
	// all directories and files should exist, or have been created by now
	fd, err := os.OpenFile(filepath, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v\n", err)
	}
	if err := fd.Sync(); err != nil {
		return nil, fmt.Errorf("error calling sync: %v\n", err)
	}
	// return file descriptor
	return fd, nil
}

func grow(fd *os.File) error {
	fi, err := os.Stat(fd.Name())
	if err != nil {
		return fmt.Errorf("error getting file stats: %v\n", err)
	}
	fdsize := fi.Size()
	if fi.Size() == 0 {
		fdsize = 1 << 20 // 1MB
	}
	page := int64(os.Getpagesize())
	size := ((fdsize * 2) + page - 1) &^ (page - 1)
	if err := fd.Truncate(size); err != nil {
		return fmt.Errorf("error truncating file: %v\n", err)
	}
	return nil
}

func (e *Engine) Read(p []byte) (int, error) {
	return e.file.Read(p)
}

func (e *Engine) Write(p []byte) (int, error) {
	return e.file.Write(p)
}

func (e *Engine) Seek(offset int64, whence int) (int64, error) {
	return e.file.Seek(offset, whence)
}

func (e *Engine) ReadAt(p []byte, off int64) (int, error) {
	return e.file.ReadAt(p, off)
}

func (e *Engine) WriteAt(p []byte, off int64) (int, error) {
	return e.file.WriteAt(p, off)
}

func (e *Engine) Close() {
	if err := e.file.Sync(); err != nil {
		log.Fatalf("error calling sync: %v\n", err)
	}
	if err := e.file.Close(); err != nil {
		log.Fatalf("error calling close: %v\n", err)
	}
}
