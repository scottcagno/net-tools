package ngin

import (
	"os"
	"path/filepath"
)

func openFile(path string) (*os.File, error) {
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
		// close the file
		if err = fd.Close(); err != nil {
			return nil, err
		}
	}
	// already existing
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, os.ModeSticky)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

func truncateFile(fd *os.File, size int64) error {
	// initally size it to 4MB
	if err := fd.Truncate(size); err != nil {
		return err
	}
	// close the file
	if err := fd.Close(); err != nil {
		return nil
	}
	return nil
}
