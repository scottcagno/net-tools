package file

import (
	"log"
	"os"
)

func Pwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
