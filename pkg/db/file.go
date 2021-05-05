package db

import (
	"log"
	"os"
	"path/filepath"
)

type DataFile struct {
	file *os.File
	info os.FileInfo
}

func OpenDataFile(path string) *DataFile {
	file, err := openFile(path)
	if err != nil {
		log.Panicf("error: %v\n", err)
	}
	info, err := file.Stat()
	if err != nil {
		log.Panicf("error: %v\n", err)
	}
	return &DataFile{
		file: file,
		info: info,
	}
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
