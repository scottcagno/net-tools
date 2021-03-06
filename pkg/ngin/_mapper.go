package ngin

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

const (
	KB int = (1 << 10) // kilobyte
	MB int = (1 << 20) // megabyte

	PAGE  int  = 4 * KB // page size
	EMPTY byte = 0xC1   // empty marker
)

type Mapper struct {
	file   *os.File // underlying file
	data   mmap     // memory mapping
	count  int      // record count
	cursor int      // coursor
}

func align(size int) int {
	if size > 0 {
		return ((size + 2) + PAGE - 1) &^ (PAGE - 1)
	}
	return PAGE // 4KB
}

func OpenMapper(path string) (*Mapper, error) {
	path += `.db`
	// check to see if we need to create a new file
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		// sanitize the filepath
		dirs, _ := filepath.Split(path)
		// create any directories
		if err := os.MkdirAll(dirs, os.ModeDir); err != nil {
			return nil, err
		}
		// create the new file
		fd, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		// initally size it to 4MB
		if err = fd.Truncate(4 * MB); err != nil {
			return nil, err
		}
		// mark beginning of each page with empty marker
		for off := 0; off < (4 * MB); off += PAGE {
			if _, err := fd.Write([]byte{EMPTY}, off); err != nil {
				log.Fatalf("OpenMapper: mark beg of each page with empty marker: ERROR BELOW\n\t%s\n", err)
				return nil, err
			}
		}
		// close the file
		if err = fd.Close(); err != nil {
			return nil, err
		}
	}
	// already existing
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, os.ModeSticky)
	if err != nil {
		return nil, err
	}
	fi, err := fd.Stat()
	if err != nil {
		return nil, err
	}
	// map file into virtual address space
	mm, err := mmap_at(0, fd.Fd(), 0, fi.Size(), PROT, FLAGS)
	if err != nil {
		return nil, err
	}
	// create new mapper instance
	m := &Mapper{fd, mm, 0, 0}
	// populate record count
	for m.cursor < fi.Size() {
		meta := m.GetMeta(m.cursor)
		if meta.IsEmpty {
			m.cursor++
			continue
		}
		m.count++
		m.cursor += meta.PgCount
	}
	m.cursor = 0
	return m, nil
}

func (m *Mapper) GetMeta(pos int) *struct {
	IsEmpty bool
	PgCount int
} {
	off := pos * PAGE
	if off%PAGE != 0 {
		return nil
	}
	return &struct {
		IsEmpty bool
		PgCount int
	}{m.data[off] == EMPTY, int(m.data[off+1])}
}

// return offset of next available n*pages
func (m *Mapper) FindEmpty(n int) (int, bool) {
	var npages int
	for m.cursor < len(m.data) {
		if npages == n {
			return m.cursor / PAGE, true
		}
		meta := m.GetMeta(m.cursor)
		if meta.IsEmpty {
			m.cursor++
			npages++
			continue
		}
		m.cursor += meta.PgCount
		npages = 0
	}
	// NOTE:  check back for empty pages
	return m.cursor / PAGE, false
}

func (m *Mapper) Read(b []byte) (int, error) {
	return -1, nil
}

// add a new record to the mapper at the first available slot
// return a non-nil error if there is an issue growing the file
func (m *Mapper) Write(b []byte) (int, error) {
	// pgs is bigger number
	pgs := align(len(b))

	pos, ok := m.FindEmpty(pgs)
	if !ok {
		if err := m.grow(); err != nil {
			return -1, err
		}
	}
	m.write(pos*PAGE, pgs, b)
	return pos, nil

}

func (m *Mapper) write(off, pgs int, b []byte) {
	m.data[off+1] = byte(pgs)
	copy(m.data[off+2:off+2+pgs], b)
}

// grow the underlying mapped file
func (m *Mapper) grow() error {
	// resize the size to double the current, ie. len * 2
	size := int64(((len(m.data) * 2) + PAGE - 1) &^ (PAGE - 1))
	// unmap current mapping before growing underlying file...
	m.data.Munmap()
	// truncate underlying file to updated size, check for errors
	if err := syscall.Ftruncate(int(m.file.Fd()), size); err != nil {
		return err
	}
	// remap underlying file now that it has grown
	mm, err := mmap_at(0, m.file.Fd(), 0, size, PROT, FLAGS)
	if err != nil {
		return err
	}
	m.data = mm

	// there were no errors, so return nil
	return nil
}

// add a new record to the engine at the first available slot
// return a non-nil error if there is an issue growing the file
func (e *engine) addRecord(r *record) (int, error) {
	// initialize block position k at beginning of mapped file, as well as future byte offset
	var k, o int
	// start iterating through mapped file reigon one page at a time
	for o < len(e.data) {
		// checking for empty page
		if bytes.Equal(e.data[o:o+e.page], e.zero) {
			// found an empty page, re-use it; copy data into it
			copy(e.data[o:o+e.page], r.data)
			// return location of block in page offset
			return o / e.page, nil
		}
		// go to next page offset
		k++
		o = k * e.page
	}
	// haven't found any empty pages, so let's grow the file
	if err := e.grow(); err != nil {
		return -1, err
	}
	// write.data to page
	copy(e.data[o:o+e.page], r.data)
	// return location of block in page offset
	return o / e.page, nil
}

func (m *Mapper) ReadAt(b []byte, pos int64) (int, error) {
	return -1, nil
}

func (m *Mapper) WriteAt(b []byte, pos int64) (int, error) {
	return -1, nil
}

// close the mapper, return any errors encountered
func CloseMapper(m *Mapper) error {
	m.data.Munmap()                        // unmap memory mappings (Munmap automatically flushes)
	if err := m.file.Close(); err != nil { // close underlying file
		return err
	}
	// set everything to nil and gc before closing
	m.file, m.data, m.count = nil, nil, -1
	runtime.GC()
	return nil
}

// add a new record to the engine at the first available slot
// return a non-nil error if there is an issue growing the file
func (e *engine) addRecord(r *record) (int, error) {
	// initialize block position k at beginning of mapped file, as well as future byte offset
	var k, o int
	// start iterating through mapped file reigon one page at a time
	for o < len(e.data) {
		// checking for empty page
		if bytes.Equal(e.data[o:o+e.page], e.zero) {
			// found an empty page, re-use it; copy data into it
			copy(e.data[o:o+e.page], r.data)
			// return location of block in page offset
			return o / e.page, nil
		}
		// go to next page offset
		k++
		o = k * e.page
	}
	// haven't found any empty pages, so let's grow the file
	if err := e.grow(); err != nil {
		return -1, err
	}
	// write.data to page
	copy(e.data[o:o+e.page], r.data)
	// return location of block in page offset
	return o / e.page, nil
}

// update a record at provided offset, assuming one exists
// return a non-nil error if offset is outside of mapped reigon
func (e *engine) setRecord(k int, r *record) error {
	// get byte offset from block position k
	o := k * e.page
	// do a bounds check; if outside of mapped reigon...
	if o+e.page > len(e.data) {
		// do not grow, return an error
		return fmt.Errorf("engine[set]: cannot update record at block %d (offset %d)\n", k, o)
	}
	// wipe page in case updated data is smaller than original dataset
	copy(e.data[o:o+e.page], e.zero)
	// write updated data to page
	copy(e.data[o:o+e.page], r.data)
	// there were no errors, so return nil
	return nil
}

var ErrEmptyRecord error = errors.New("engine: empty record found")
var ErrEngineEOF error = io.EOF

// return a record at provided offset, assuming one exists
// return a non-nil error if offset is outside of mapped reigon
func (e *engine) getRecord(k int) (*record, error) {
	// get byte offset from block position k
	o := k * e.page
	// do a bounds check; if outside of mapped reigon...
	if o+e.page > len(e.data) {
		// ...return an error
		return nil, ErrEngineEOF //fmt.Errorf("engine[get]: cannot return record at block %d (offset %d)\n", k, o)
	}
	// create record to return
	r := new(record)
	// fill out record data if not empty, returning no error
	if n := bytes.LastIndexByte(e.data[o+maxKey-1:o+e.page], eofVal); n > 0 {
		r.data = e.data[o : o+n]
		return r, nil
	}
	// otherwise, return empty record, with an error
	return r, ErrEmptyRecord //fmt.Errorf("engine[get]: empty record found at block %d (offset %d)", k, o)
}

func (e *engine) getRecordKey(k int) ([]byte, error) {
	// get byte offset from block position k
	o := k * e.page
	// do a bounds check; if outside of mapped reigon...
	if o+e.page > len(e.data) {
		// ...return an error
		return nil, fmt.Errorf("engine[getKey]: cannot return key at block %d (offset %d)\n", k, o)
	}
	if !bytes.Equal(e.data[o:o+page], e.zero) {
		return e.data[o : o+maxKey], nil
	}
	// otherwise, return empty record, with an error
	return nil, fmt.Errorf("engine[getKey]: empty key found at block %d (offset %d)", k, o)
}

func (e *engine) getRecordVal(k int) ([]byte, error) {
	// get byte offset from block position k
	o := k * e.page
	// do a bounds check; if outside of mapped reigon...
	if o+e.page > len(e.data) {
		// ...return an error
		return nil, fmt.Errorf("engine[getVal]: cannot return val at block %d (offset %d)\n", k, o)
	}
	// fill out record data if not empty, returning no error
	if n := bytes.LastIndexByte(e.data[o+maxKey:o+e.page], eofVal); n > 0 {
		v := e.data[o+maxKey : o+maxKey+n]
		return v, nil
	}
	// otherwise, return empty record, with an error
	return nil, fmt.Errorf("engine[getVal]: empty val found at block %d (offset %d)", k, o)
}

// delete a record at provided offset, assuming one exists
// return a non-nil error if offset is outside of mapped reigon
func (e *engine) delRecord(k int) error {
	// get byte offset from block position k
	o := k * e.page
	// do a bounds check; if outside of mapped reigon...
	if o+e.page > len(e.data) {
		// ...return an error
		return fmt.Errorf("engine[del]: cannot delete record at block %d (offset %d)\n", k, o)
	}
	// otherwise, wipe page block at offset
	copy(e.data[o:o+page], e.zero)
	// there were no errors, so return nil
	return nil
}

// close the engine, return any errors encountered
func (e *engine) close() error {
	e.munmap()                             // unmap memory mappings (Munmap automatically flushes)
	if err := e.file.Close(); err != nil { // close underlying file
		return err
	}
	e.file = nil // set file descriptor to nil
	// there were no errors, so return nil
	return nil
}

// temp structure
type payload struct {
	key []byte
	pos int
}

// get all of the record data payloads from the engine
func (e *engine) loadAllRecords() <-chan payload {
	// initialize the channels to return the keys and blocks on
	loader := make(chan payload)
	go func() {
		var o, k int
		// start iterating through mapped file reigon one page at a time
		for o < len(e.data) {
			// checking for non-empty page
			if !bytes.Equal(e.data[o:o+e.page], e.zero) {
				// found one; return key and block offset
				loader <- payload{e.data[o : o+maxKey], (o / e.page)}
			}
			k++
			o = k * e.page
		}
		close(loader)
	}()
	return loader
}

/*
func Mmap(fd *os.File, off, len int) mmap {
	mm, err := mmap_at(0, fd.Fd(), int64(off), int64(len), PROT, FLAGS)
	if err != nil {
		panic(err)
	}
	return mm
}
*/
