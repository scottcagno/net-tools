package engine

import "fmt"

// !!! working, but check to make sure maxd still works properly if page size changes
// https://play.golang.org/p/_jyyoAZn0A_z

var page = 16 //os.Getpagesize()
var maxd = 0xfffff - page

func align(n int) int {
	return ((1 + n) + page - 1) &^ (page - 1)
}

func record(d []byte) ([]byte, error) {
	dn := len(d)
	if dn > maxd-1 {
		return nil, fmt.Errorf("error: data too large: %d > MAX(%d)\n", dn, maxd-1)
	}
	// get proper alignment
	sz := align(dn)
	// make document record
	dc := make([]byte, sz, sz)
	// write header to document record
	copy(dc[0:1], []byte{byte(sz / page)})
	// write data
	copy(dc[1:], d)
	return dc, nil
}
