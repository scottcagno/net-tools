package main

import (
	"fmt"
	"github.com/scottcagno/net-tools/pkg/db"
	"io"
	"log"
	"os"
)

func main() {

	// open data file
	df, err := db.OpenDataFile("cmd/db/data/test.txt")
	if err != nil {
		log.Panic(err)
	}

	// get data file info
	i, err := df.Stat()
	checkErr(err)
	fmt.Printf("file info: name=%q, size=%d, modified=%v\n", i.Name(), i.Size(), i.ModTime())

	// check if need to resize
	kb := int64(1 << 10)
	if i.Size() < 4*kb {
		err = df.Truncate(4 * kb)
		checkErr(err)
		fmt.Printf("resizing file...\n")
	}

	seekBlock(df, 0)
	readBlockData(df, 5)

	checkErr(err)

	// close data file
	err = df.Close()
	checkErr(err)
}

func writeBlockData(df *db.DataFile, n int) {
	// writing block data
	seekBlock(df, 0)
	for i := 0; i < n; i++ {
		d := fmt.Sprintf("%d:this is data for block #%d\n", i, i)
		err := df.WriteData([]byte(d))
		checkErr(err)
	}
}

func readBlockData(df *db.DataFile, n int) {
	// reading block data
	seekBlock(df, 0)
	for i := 0; i < n; i++ {
		d, err := df.ReadData()
		checkErr(err)
		fmt.Printf("reading data at block %d: %s\n", i, d)
	}
}

func test() {

	fd, err := db.OpenFile("cmd/db/data/debug.txt")
	checkErr(err)

	_, err = fd.Seek(128, io.SeekStart)
	checkErr(err)

	//r := bufio.NewReaderSize(fd, 16)
	//w := bufio.NewWriterSize(fd, 16)
	//rw := bufio.NewReadWriter(r, w)

	_, err = fd.Write([]byte("FOOOOOOOO\n"))
	checkErr(err)

	//err = rw.Flush()
	//checkErr(err)

	err = fd.Close()
	checkErr(err)

	os.Exit(0)
}

func seekBlock(df *db.DataFile, blk int64) {
	n, err := df.SeekBlock(blk, io.SeekStart)
	checkErr(err)
	fmt.Printf("seek block: %d (offset %d)\n", blk, n)
}

func writeString(df *db.DataFile, n int) {
	mystring := fmt.Sprintf("%d:foobar\n", n)
	err := df.WriteString(mystring)
	checkErr(err)
	fmt.Printf("wrote string: %q\n", mystring)
}
func readString(df *db.DataFile) {
	s, err := df.ReadString()
	checkErr(err)
	fmt.Printf("read string: %q\n", s)
}

func checkErr(err error) {
	if err != nil {
		log.Printf("[error] %v\n", err)
	}
}
