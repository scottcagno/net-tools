package file

import (
	"bufio"
	"io"
)

type Writer struct {
	Delim byte
	w     *bufio.Writer
}

func NewWriter(w io.WriteSeeker) *Writer {
	return &Writer{
		Delim: '\n',
		w:     bufio.NewWriter(w),
	}
}

func (w *Writer) Write(line []byte) error {
	if line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}
	if _, err := w.w.Write(line); err != nil {
		return err
	}
	return w.w.WriteByte('\n')
}

func (w *Writer) Flush() {
	w.w.Flush()
}

func (w *Writer) Error() error {
	_, err := w.w.Write(nil)
	return err
}

func (w *Writer) WriteAll(lines [][]byte) error {
	for _, line := range lines {
		err := w.Write(line)
		if err != nil {
			return err
		}
	}
	return w.w.Flush()
}
