package io

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

type File struct {
	*os.File
}

func NewFile(filepath string) (*File, error) {
	// split path and file from filepath
	dir, file := path.Split(filepath)
	// if there are any nested directories...
	if dir != "" {
		// check to see if directories already exist...
		if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
			// if not, make all nested directories...
			if err := os.MkdirAll(dir, 0755); err != nil {
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
			if _, err := os.Create(filepath); err != nil {
				return nil, fmt.Errorf("error creating file: %v\n", err)
			}
		}
		// case: file exists, continue
	}
	// all directories and files should exist, or have been created by now
	// was: os.OpenFile(filepath, os.O_APPEND|os.O_RDWR, 0644), lets drop the append for now
	fd, err := os.OpenFile(filepath, os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v\n", err)
	}
	if err := fd.Sync(); err != nil {
		return nil, fmt.Errorf("error calling sync: %v\n", err)
	}
	// return file descriptor
	return &File{fd}, nil
}

func OpenFile(filepath string) (*File, error) {
	return NewFile(filepath)
}

func (f *File) ReadFromLine(lineNum int) ([]byte, error) {
	br := bufio.NewReader(f)
	for i := 1; i < lineNum; i++ {
		_, _ = br.ReadSlice('\n')
	}
	line, err := br.ReadSlice('\n')
	if err != nil && err != io.EOF {
		log.Printf("Encountered error readline line %d: %v\n", lineNum, err)
	}
	return line, nil
}

func FileIndex(r io.Reader, delim byte) []int {
	br := bufio.NewReader(r)
	var line []byte
	var err error
	var idx []int
	for i := 0; err != io.EOF; i++ {
		line, err = br.ReadSlice(delim)
		idx = append(idx, len(line))
	}
	return idx
}
