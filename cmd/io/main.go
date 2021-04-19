package main

import (
	"fmt"
	"github.com/scottcagno/net-tools/pkg/file"
	"io"
	"log"
)

func OpenFileAndReadDataIntoBuffer(filepath string, lines *[][]byte, verbose bool) {
	// open a file
	if verbose {
		fmt.Printf("Opening file (%s) and reading data into memory\n", filepath)
	}
	f1, err := file.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	// create place to store lines in memory
	//var lines [][]byte

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
		*lines = append(*lines, line)
		if verbose {
			fmt.Printf("Reading line %d into memory...\n", ln)
		}
		ln++
	}
	if verbose {
		fmt.Printf("Finished. Closing file (%s)\n", filepath)
	}
	// close file, and open different file
	f1.Close()
}

func OpenFileAndWriteDataFromBuffer(filepath string, lines *[][]byte, verbose bool) {
	// open a new file and write contents to it
	if verbose {
		fmt.Printf("Opening file (%s) and writing data from memory\n", filepath)
	}
	f2, err := file.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	// write contents from memory onto disk
	w := file.NewWriter(f2)
	for ln, line := range *lines {
		if verbose {
			fmt.Printf("Writing line %d into file...\n", ln)
		}
		if err := w.Write(line); err != nil {
			log.Fatalf("error writing line to file: %v\n", err)
		}
	}
	if verbose {
		fmt.Printf("Finished. Flush, error check and close file (%s)\n", filepath)
	}
	w.Flush() // flush writer and check for any errors
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	f2.Close()
}

func main() {

	filepath1, verbose := "cmd/io/test/in.txt", false
	var lines [][]byte
	OpenFileAndReadDataIntoBuffer(filepath1, &lines, false)

	if verbose {
		for n, line := range lines {
			fmt.Printf(">> %d: %s\n", n, line)
		}
	}

	var n int
	fmt.Printf(">> line %d: %s\n", n, lines[n])

	filepath2 := "cmd/io/test/out.txt"
	OpenFileAndWriteDataFromBuffer(filepath2, &lines, false)

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
