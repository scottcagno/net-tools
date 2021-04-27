package file

import (
	"bufio"
	"io"
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

func (r *Reader) Read() ([]byte, error) {
	line, err := r.readLine()
	return line, err
}

func (r *Reader) ReadAll() ([][]byte, error) {
	var lines [][]byte
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
