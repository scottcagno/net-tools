package main

import (
	"encoding/binary"
	"fmt"
	"github.com/scottcagno/net-tools/pkg/engine"
	"io"
	"os"
)

func main() {
	e := engine.Open("cmd/engine/test/data.txt")
	for i := 0; i < 10; i++ {
		b := make([]byte, 8)
		n, _ := e.Read(b)
		d := binary.BigEndian.Uint64(b)
		fmt.Printf("read %d bytes, (%d)", n, d)
		e.Seek(int64(os.Getpagesize()-len(b)), io.SeekCurrent)
	}
	e.Close()
}
