package main

import (
	"fmt"
	"github.com/scottcagno/net-tools/pkg/file"
	"io"
	"log"
	"strings"
)

func OpenFileAndReadDataIntoBuffer(filepath string, lines *[][]byte, verbose bool) {
	// open a file
	if verbose {
		fmt.Printf("Opening file (%s) and reading data into memory\n", filepath)
	}
	f, err := file.NewFileHandler(filepath)
	if err != nil {
		log.Fatal(err)
	}

	// create place to store lines in memory
	//var lines [][]byte

	// read contents into memory
	ln := 1
	for {
		line, err := f.ReadLine()
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
	f.Close()
}

func OpenFileAndWriteDataFromBuffer(filepath string, lines *[][]byte, verbose bool) {
	// open a new file and write contents to it
	if verbose {
		fmt.Printf("Opening file (%s) and writing data from memory\n", filepath)
	}

	// write contents from memory onto disk
	f, err := file.NewFileHandler(filepath)
	if err != nil {
		log.Fatal(err)
	}
	for ln, line := range *lines {
		if verbose {
			fmt.Printf("Writing line %d into file...\n", ln)
		}
		if err := f.WriteLine(line); err != nil {
			log.Fatalf("error writing line to file: %v\n", err)
		}
	}
	if verbose {
		fmt.Printf("Finished. Flush, error check and close file (%s)\n", filepath)
	}
	f.Flush() // flush writer and check for any errors
	if err := f.Error(); err != nil {
		log.Fatal(err)
	}
	f.Close()
}

func handleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main() {

	f, err := file.NewFileHandler("cmd/io/test/in.txt")
	handleErr(err)

	lc := f.LineCache()
	for i, l := range lc {
		if i == 0 {
			continue
		}
		fmt.Printf("line %d begins at offset %d\n", i, l)
	}
	return

	var ln int64 = 32
	err = f.SeekLine(ln)
	handleErr(err)

	line, err := f.ReadLine()
	handleErr(err)

	fmt.Printf("go to line %d, read, and print: %s\n", ln, line)

	return

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

func TestLargeFiles(filepath string, verbose bool) {
	var lines [][]byte
	OpenFileAndReadDataIntoBuffer(filepath, &lines, verbose)

	filepath = strings.Replace(filepath, ".txt", ".out.txt", 1)
	OpenFileAndWriteDataFromBuffer(filepath, &lines, verbose)
}
