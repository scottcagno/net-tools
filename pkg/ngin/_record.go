package ngin

import (
	"bufio"
	"bytes"
)

const (
	szKB     = 1 << 10
	szMB     = 1 << 20
	szGB     = 1 << 30
	szHeader = 16
	szPage   = 4 * szKB
)

// alignBytes takes n number of input bytes (should be length of the data) and returns
// a page aligned byte count, making sure to include the header in the calculation.
func alignBytes(n int) uint64 {
	if n > 0 {
		return uint64(((n + szHeader) + szPage - 1) &^ (szPage - 1))
	}
	return uint64(szPage)
}

// alignPages takes n number of input bytes (should be length of the data) and returns
// a page aligned page count, making sure to include the header in the calculation.
func alignPages(n int) uint16 {
	return uint16(alignBytes(n) / szPage)
}

// align takes n number of input bytes (should be length of the data) and returns
// both a page aligned byte count, and a page aligned page count, making sure to
// include the header in the calculation.
func align(n int) (uint64, uint16) {
	return alignBytes(n), alignPages(n)
}

// header represents a the header of a data record. It stores the status of the
// record, a magic byte, an extra uint16 marker, a page count marker, the length
// of the actual data in the record and the padding size to page align the record.
type header struct {
	status  byte   // 0 - free, 1 - active, 2 - deleted
	magic   byte   // magic is currently unused, but was put there for future use (in case)
	extra   uint16 // extra is currently unused, but was put there for future use (in case)
	pages   uint16 // number of aligned pages, 65535 pages is the max (255mb)
	length  uint64 // total length of record in bytes, 268431360 (not bound by uint64 type) bytes is the max (255mb)
	padding uint16 // number of bytes to pad after the header and raw data
}

// record represents a data record. It has a small header that stores the size
// of the record as well as some other markers for empty, deleted, etc. It may
// vary in size, but will always be perfectly page aligned. The header bytes
// (including the padding) occupy 16 bytes in addition to the raw data itself.
type record struct {
	*header        // embedded header
	data    []byte // actual data
}

// newHeader takes the length of the raw data, known as 'dl', and creates
// and returns a new filled header struct based on provided data length
func newHeader(dl int) *header {
	abc, apc := align(dl)
	return &header{
		status:  byte(1),
		magic:   byte(0),
		extra:   uint16(0),
		pages:   uint16(apc),
		length:  uint64(dl),
		padding: uint16(abc - uint64(dl+szHeader)),
	}
}

func newRecord(b []byte) *record {
	return &record{
		header: newHeader(len(b)),
		data:   b,
	}
}

func (r *record) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	w := NewWriter(bufio.NewWriter(&buf))
	// write header
	err := w.WriteByte(r.status)
	if err != nil {
		return nil, err
	}
	err = w.WriteByte(r.magic)
	if err != nil {
		return nil, err
	}
	err = w.WriteUint16(r.extra)
	if err != nil {
		return nil, err
	}
	err = w.WriteUint16(r.pages)
	if err != nil {
		return nil, err
	}
	err = w.WriteUint64(r.length)
	if err != nil {
		return nil, err
	}
	err = w.WriteUint16(r.padding)
	if err != nil {
		return nil, err
	}
	// write data
	err = w.WriteBytes(r.data)
	if err != nil {
		return nil, err
	}
	// write padding
	err = w.WriteBytes(make([]byte, r.padding, r.padding))
	if err != nil {
		return nil, err
	}
	err = w.Flush()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *record) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	rd := NewReader(bufio.NewReader(buf))
	// read header
	var err error
	r.header.status, err = rd.ReadByte()
	if err != nil {
		return err
	}
	// TODO: CONTINUE FROM HERE
	return nil
	err = w.WriteByte(r.magic)
	if err != nil {
		return err
	}
	err = w.WriteUint16(r.extra)
	if err != nil {
		return err
	}
	err = w.WriteUint16(r.pages)
	if err != nil {
		return err
	}
	err = w.WriteUint64(r.length)
	if err != nil {
		return err
	}
	err = w.WriteUint16(r.padding)
	if err != nil {
		return err
	}
	// write data
	err = w.WriteBytes(r.data)
	if err != nil {
		return err
	}
	// write padding
	err = w.WriteBytes(make([]byte, r.padding, r.padding))
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	return nil
}
