package main

import (
	"fmt"
	"github.com/scottcagno/net-tools/pkg/file"
	"io"
	"log"
)

func main() {

	// open a file
	filepath1 := "cmd/io/test.txt"
	fmt.Printf("Opening file (%s) and reading data into memory\n", filepath1)
	f1, err := file.Open(filepath1)
	if err != nil {
		log.Fatal(err)
	}

	// create place to store lines in memory
	var lines [][]byte

	// read contents into memory
	r, ln := file.NewReader(f1), 1
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		lines = append(lines, line)
		fmt.Printf("Reading line %d into memory...\n", ln)
		ln++
	}

	fmt.Printf("Finished. Closing file (%s)\n", filepath1)
	// close file, and open different file
	f1.Close()

	// open a new file and write contents to it
	filepath2 := "cmd/io/out/test-copy.txt"
	fmt.Printf("Opening file (%s) and writing data from memory\n", filepath2)
	f2, err := file.Open(filepath2)
	if err != nil {
		log.Fatal(err)
	}

	// write contents from memory onto disk
	w := file.NewWriter(f2)
	for ln, line := range lines {
		fmt.Printf("Writing line %d into file...\n", ln)
		if err := w.Write(line); err != nil {
			log.Fatalf("error writing line to file: %v\n", err)
		}
	}

	fmt.Printf("Finished. Flush, error check and close file (%s)\n", filepath2)
	w.Flush() // flush writer and check for any errors
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	f2.Close()
}

func foo() {
	var err error
	// read all the lines in the file
	fmt.Println("Reading all the lines in the file")
	err = file.Read("cmd/io/test.txt")
	if err != nil {
		log.Println(err)
	}

	// read specified lines in the file
	fmt.Println("Reading line twenty three in the file")
	err = file.ReadLine("cmd/io/test.txt", 23)
	if err != nil {
		log.Println(err)
	}
}
