package main

import (
	"encoding/binary"
	"fmt"
	"github.com/scottcagno/net-tools/pkg/engine"
	"io"
	"os"
)

// https://play.golang.org/p/gaWCeLMiXZG

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

// blocksize is the minimum size a record should align to
const blocksize = 8

// align aligns by rounding the input number 'n'
// up to the nearest supplied block size 'b'
func align(n, b int) int {
	bsz := b - 1 // blocksize
	return (n + bsz) &^ bsz
}

// pad takes data 'd' and pads it by 'n' bytes
func pad(d []byte, n int) []byte {
	d = append(d, make([]byte, n)...)
	return d
}

func insert(s []byte, k int, vs ...byte) []byte {
	if n := len(s) + len(vs); n <= cap(s) {
		s2 := s[:n]
		copy(s2[k+len(vs):], s[k:])
		copy(s2[k:], vs)
		return s2
	}
	s2 := make([]byte, len(s)+len(vs))
	copy(s2, s[:k])
	copy(s2[k:], vs)
	copy(s2[k+len(vs):], s[k:])
	return s2
}

// create new record
func record(d []byte) []byte {
	// pad out the record to the specified blocksize
	d = insert(d, 0, []byte{0, 0}...)
	z := align(len(d), blocksize) - len(d)
	d = pad(d, z)
	blocks := len(d) / blocksize
	binary.PutUvarint(d[0:2], uint64(blocks))
	return d
}
