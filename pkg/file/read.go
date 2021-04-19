package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Reader struct {
	Delim     byte
	r         *bufio.Reader
	numLine   int
	rawBuffer []byte
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		Delim: '\n',
		r:     bufio.NewReader(r),
	}
}

func (r *Reader) LineNumber() int {
	return r.numLine
}

func (r *Reader) Read2() (line []byte, err error) {
	line, err = r.readLine2()
	return line, err
}

func (r *Reader) Read() (line []byte, err error) {
	line, err = r.readLine()
	return line, err
}

func (r *Reader) ReadAll() (lines [][]byte, err error) {
	for {
		line, err := r.readLine()
		if err == io.EOF {
			return lines, nil
		}
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
}

func (r *Reader) readLine2() ([]byte, error) {
	line, pref, err := make([]byte, 0), true, error(nil)
	for pref && err == nil {
		line, pref, err = r.r.ReadLine()
		line = append(line, line...)
	}
	r.numLine++
	return line, err
}

func (r *Reader) readLine() ([]byte, error) {
	line, err := r.r.ReadSlice('\n')
	if err == bufio.ErrBufferFull {
		r.rawBuffer = append(r.rawBuffer[:0], line...)
		for err == bufio.ErrBufferFull {
			line, err = r.r.ReadSlice('\n')
			r.rawBuffer = append(r.rawBuffer, line...)
		}
		line = r.rawBuffer
	}
	if len(line) > 0 && err == io.EOF {
		err = nil
		// For backwards compatibility, drop trailing \r before EOF.
		if line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
	}
	r.numLine++
	// Normalize \r\n to \n on all input lines.
	if n := len(line); n >= 2 && line[n-2] == '\r' && line[n-1] == '\n' {
		line[n-2] = '\n'
		line = line[:n-1]
	}
	return line, err
}

const defaultBufferSize = 1 << 12 // 4KB

func Read(path string) error {
	return read(path, -1)
}

func ReadLine(path string, line int) error {
	return read(path, line)
}

func read(path string, line int) error {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error opening file: %v\n", err)
	}
	defer f.Close()
	rd := bufio.NewReaderSize(f, defaultBufferSize) // 4KB
	for ln := 1; ln != line+1; ln++ {
		b, err := rd.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("(EOF) %d: %s\n", ln, b)
				break
			}
			return fmt.Errorf("error reading file line: %v\n", err)
		}
		if line == -1 || line == ln {
			_ = b
			fmt.Printf("%d: %s\n", ln, b)
		}
	}
	return nil
}
