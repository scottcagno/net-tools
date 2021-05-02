package engine

import (
	"bytes"
	"fmt"
)

// !!! working, but check to make sure maxd still works properly if page size changes
// https://play.golang.org/p/_jyyoAZn0A_z
// https://play.golang.org/p/Gq5cGzY7D-g <<--- UPDATED VERSION!

var (
	page   = 16 //os.Getpagesize()
	maxlen = (0xff * page) - 1
	recsep = byte(0x1e)
	header = func(asz int) []byte {
		// return aligned size
		return []byte{byte(asz / page)}
	}
)

func align(n int) int {
	return ((2 + n) + page - 1) &^ (page - 1)
}

type record struct {
	d []byte
}

func makerecord(data []byte) (*record, error) {
	datalen := len(data)
	if datalen > maxlen {
		return nil, fmt.Errorf("error: data too large: %d > MAX(%d)\n", datalen, maxlen)
	}
	// get proper alignment
	sz := align(datalen)
	// make document record
	rec := make([]byte, sz, sz)
	// write header to record
	copy(rec[0:1], header(sz))
	// write data to record
	copy(rec[1:], data)
	// write record separator to end record
	copy(rec[1+len(data):], []byte{recsep})
	return &record{d: rec}, nil
}

func (r *record) empty() bool {
	return r.d[0] == 0x00
}

func (r *record) pages() int {
	return int(r.d[0])
}

func (r *record) data() []byte {
	end := bytes.LastIndexByte(r.d, recsep)
	if end == -1 || r.d[0] == 0x00 {
		return nil
	}
	return r.d[1:end]
}

func (r *record) record() []byte {
	return r.d
}

func (r *record) remove() bool {
	end := bytes.LastIndexByte(r.d, recsep)
	if end == -1 || r.d[0] == 0x00 {
		return false
	}
	r.d[0], r.d[end] = 0x00, byte(0x00)
	return true
}
