package ngin

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

// record represents a data record. It has a small header that stores the size
// of the record as well as some other markers for empty, deleted, etc. It may
// vary in size, but will always be perfectly page aligned. The header bytes
// (including the padding) occupy 16 bytes in addition to the raw data itself.
type record struct {
	status  byte   // 0 - free, 1 - active, 2 - deleted
	extra   byte   // extra is currently unused, but was put there for future use (in case)
	pages   uint16 // number of aligned pages, 65535 pages is the max (255mb)
	length  uint64 // total length of record in bytes, 268431360 (not bound by uint64 type) bytes is the max (255mb)
	padding uint16 // number of bytes to pad after the header and raw data
	data    []byte // actual data
}

func newRecord(b []byte) *record {
	abc, apc := align(len(b))
	return &record{
		status:  1,
		pages:   apc,
		length:  uint64(len(b)),
		padding: uint16(abc - uint64(len(b)+szHeader)),
		data:    b,
	}
}

func (r *record) MarshalBinary() ([]byte, error) {
	return nil, nil
}

func (r *record) UnmarshalBinary(data []byte) error {
	return nil
}
