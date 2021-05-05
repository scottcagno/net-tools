package main

import (
	"fmt"
	"github.com/scottcagno/net-tools/pkg/db"
	"log"
)

func main() {

	df, err := db.OpenDataFile("cmd/db/data/test.txt")
	if err != nil {
		log.Panic(err)
	}

	i, err := df.Stat()
	checkErr(err)
	fmt.Printf("file info: name=%q, size=%d, modified=%v\n", i.Name(), i.Size(), i.ModTime())

	s, err := df.ReadString()
	checkErr(err)
	fmt.Printf("read string: %s", s)

	s, err = df.ReadString()
	checkErr(err)
	fmt.Printf("read string: %s", s)

	err = df.Close()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Printf("[error] %v\n", err)
	}
}
