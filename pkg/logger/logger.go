package logger

import (
	"log"
	"os"
)

var (
	Warning *log.Logger
	Info    *log.Logger
	Error   *log.Logger
)

func init() {
	fd, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	Info = log.New(fd, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(fd, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(fd, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
