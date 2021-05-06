package db

import (
	"bufio"
	"os"
)

type MetaFile struct {
	file *os.File
	*bufio.ReadWriter
}

func OpenMetaFile(path string) (*MetaFile, error) {
	file, err := OpenFile(path)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReaderSize(file, blockSZ)
	w := bufio.NewWriterSize(file, blockSZ)
	return &MetaFile{
		file,
		bufio.NewReadWriter(r, w),
	}, nil
}

func Close() error {
	return nil
}
