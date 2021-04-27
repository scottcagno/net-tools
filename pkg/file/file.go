package file

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

const defaultBufferSize = 1 << 12 // 4KB

type FileHandler struct {
	f  *os.File
	r  *Reader
	w  *Writer
	lc []int64 // line cache
}

func NewFileHandler(filepath string) (*FileHandler, error) {
	fd, err := CreateOrOpen(filepath)
	if err != nil {
		return nil, err
	}
	fh := &FileHandler{
		f:  fd,
		r:  NewReader(fd),
		w:  NewWriter(fd),
		lc: make([]int64, 1),
	}
	err = fh.initLineCache()
	if err != nil {
		return nil, err
	}
	return fh, nil
}

func (f *FileHandler) LineCache() []int64 {
	return f.lc
}

func (f *FileHandler) initLineCache() error {
	for {
		_, err := f.r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		pos, err := f.f.Seek(0, io.SeekCurrent)
		log.Printf(">>>>> pos->%d\n", pos)
		if err != nil {
			return err
		}
		f.lc = append(f.lc, pos)
	}
	_, err := f.f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileHandler) Seek(offset int64, whence int) (int64, error) {
	return f.f.Seek(offset, whence)
}

func (f *FileHandler) Truncate(size int64) error {
	return f.f.Truncate(size)
}

func (f *FileHandler) Read(b []byte) (int, error) {
	return f.f.Read(b)
}

func (f *FileHandler) Write(b []byte) (int, error) {
	return f.f.Write(b)
}

func (f *FileHandler) SeekLine(ln int64) error {
	// bounds check
	if ln > int64(len(f.lc)) {
		return fmt.Errorf("line %d cannot be found, or is out of bounds!\n", ln)
	}
	// seek to line cache byte offset
	_, err := f.f.Seek(f.lc[ln], io.SeekStart)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileHandler) CurrentLine() int64 {
	off, err := f.f.Seek(0, io.SeekCurrent)
	if err != nil {
		return -1
	}
	return off
}

func (f *FileHandler) ReadLine() ([]byte, error) {
	return f.r.Read()
}

func (f *FileHandler) ReadLineN(n int) ([]byte, error) {
	// bounds check
	if n > len(f.lc) {
		return nil, fmt.Errorf("line %d cannot be found, or is out of bounds!\n", n)
	}
	// get current cursor offset
	cur := f.CurrentLine()
	if cur == -1 {
		return nil, fmt.Errorf("error obtaining current line\n")
	}
	// seek to line cache offset
	err := f.SeekLine(f.lc[n])
	if err != nil {
		return nil, err
	}
	// read line
	line, err := f.ReadLine()
	if err != nil {
		return line, err
	}
	// set cursor back to original offset
	_, err = f.f.Seek(cur, io.SeekStart)
	if err != nil {
		return nil, err
	}
	return line, nil
}

func (f *FileHandler) WriteLine(line []byte) error {
	return f.w.Write(line)
}

func (f *FileHandler) getLineN(n int) int64 {
	if n-1 > len(f.lc) {
		return -1
	}
	return f.lc[n-1]
}

func (f *FileHandler) Close() error {
	return f.f.Close()
}

func (f *FileHandler) Flush() {
	f.w.Flush()
}

func (f *FileHandler) Error() error {
	return f.w.Error()
}

func CreateOrOpen(filepath string) (*os.File, error) {
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
	return fd, nil
}
