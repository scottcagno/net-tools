package file

import (
	"fmt"
	"os"
	"path"
)

func Open(filepath string) (*os.File, error) {
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
		// directory exists
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
		// file exists
	}
	// all directories and files should exist, or have been created by now
	fd, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v\n", err)
	}
	if err := fd.Sync(); err != nil {
		return nil, fmt.Errorf("error calling sync: %v\n", err)
	}
	// return file descriptor
	return fd, nil
}
