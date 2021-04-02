package httputil

import (
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	fd, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	InfoLogger = log.New(fd, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(fd, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(fd, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
