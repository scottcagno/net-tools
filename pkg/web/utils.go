package web

import (
	"io"
	"log"
	"os"
)

const (
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
)

func NewLogger(writer io.Writer) *log.Logger {
	if writer == nil {
		writer = os.Stderr
	}
	return log.New(writer, "", log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix)
}

func NewLoggerWithPrefix(writer io.Writer, prefix string) *log.Logger {
	logger := NewLogger(writer)
	logger.SetPrefix(prefix)
	return logger
}
