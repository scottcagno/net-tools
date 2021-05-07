package data

import (
	"bufio"
	"encoding/binary"
	"io"
	"math"
	"reflect"
	"unsafe"
)

// DataWriter for writing data
type DataWriter struct {
	bw *bufio.Writer
}

// NewDataWriter returns a new DataWriter
func NewDataWriter(w io.Writer) *DataWriter {
	return &DataWriter{
		bw: bufio.NewWriter(w),
	}
}

// NewDataWriterSize returns a new DataWriter
func NewDataWriterSize(w io.Writer, size int) *DataWriter {
	return &DataWriter{
		bw: bufio.NewWriterSize(w, size),
	}
}

// Flush the buffered bytes to the underlying writer
func (dw *DataWriter) Flush() error {
	return dw.bw.Flush()
}

// WriteUvarint writes a uvarint
func (dw *DataWriter) WriteUvarint(x uint64) error {
	var buf [10]byte
	n := binary.PutUvarint(buf[:], x)
	_, err := dw.bw.Write(buf[:n])
	return err
}

// WriteVarint writes a varint
func (dw *DataWriter) WriteVarint(x int64) error {
	var buf [10]byte
	n := binary.PutVarint(buf[:], x)
	_, err := dw.bw.Write(buf[:n])
	return err
}

// WriteUint32 writes a uint32
func (dw *DataWriter) WriteUint32(x uint32) error {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], x)
	_, err := dw.bw.Write(buf[:])
	return err
}

// WriteInt32 writes an int32
func (dw *DataWriter) WriteInt32(x int32) error {
	return dw.WriteUint32(uint32(x))
}

// WriteUint16 writes a uint16
func (dw *DataWriter) WriteUint16(x uint16) error {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], x)
	_, err := dw.bw.Write(buf[:])
	return err
}

// WriteInt16 writes an int16
func (dw *DataWriter) WriteInt16(x int16) error {
	return dw.WriteUint16(uint16(x))
}

// WriteUint8 writes a uint8
func (dw *DataWriter) WriteUint8(x uint8) error {
	return dw.bw.WriteByte(x)
}

// WriteInt8 writes an int8
func (dw *DataWriter) WriteInt8(x int8) error {
	return dw.WriteUint8(uint8(x))
}

// WriteByte writes a byte
func (dw *DataWriter) WriteByte(x byte) error {
	return dw.WriteUint8(uint8(x))
}

// WriteUint64 writes a uint64
func (dw *DataWriter) WriteUint64(x uint64) error {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], x)
	_, err := dw.bw.Write(buf[:])
	return err
}

// WriteInt64 writes an int64
func (dw *DataWriter) WriteInt64(x int64) error {
	return dw.WriteUint64(uint64(x))
}

// WriteFloat64 writes a float64
func (dw *DataWriter) WriteFloat64(f float64) error {
	return dw.WriteUint64(math.Float64bits(f))
}

// WriteFloat32 writes a float32
func (dw *DataWriter) WriteFloat32(f float32) error {
	return dw.WriteUint32(math.Float32bits(f))
}

// WriteBool writes a bool
func (dw *DataWriter) WriteBool(t bool) error {
	if t {
		return dw.bw.WriteByte(1)
	}
	return dw.bw.WriteByte(0)
}

// WriteString writes a string
func (dw *DataWriter) WriteString(s string) error {
	if err := dw.WriteUvarint(uint64(len(s))); err != nil {
		return err
	}
	_, err := dw.bw.WriteString(s)
	return err
}

func (dw *DataWriter) WriteBinary(b []byte) error {
	err := dw.WriteUint64(uint64(len(b)))
	if err != nil {
		return err
	}
	_, err = dw.bw.Write(b)
	if err != nil {
		return err
	}
	return nil
}

// WriteBytes writes bytes
func (dw *DataWriter) WriteBytes(b []byte) error {
	return dw.WriteString(*(*string)(unsafe.Pointer(&b)))
}

// DataReader for reading
type DataReader struct {
	br *bufio.Reader
}

// NewDataReader returns a new DataReader
func NewDataReader(r io.Reader) *DataReader {
	return &DataReader{
		br: bufio.NewReader(r),
	}
}

// NewDataReaderSize returns a new DataReader
func NewDataReaderSize(r io.Reader, size int) *DataReader {
	return &DataReader{
		br: bufio.NewReaderSize(r, size),
	}
}

// Peek returns the next n bytes without advancing the reader
func (dr *DataReader) Peek(n int) ([]byte, error) {
	return dr.br.Peek(n)
}

// ReadUvarint reads a uvarint
func (dr *DataReader) ReadUvarint() (uint64, error) {
	return binary.ReadUvarint(dr.br)
}

// ReadVarint reads a varint
func (dr *DataReader) ReadVarint() (int64, error) {
	return binary.ReadVarint(dr.br)
}

func (dr *DataReader) Discard(n int) (int, error) {
	return dr.br.Discard(n)
}

// ReadUint64 reads a uint64
func (dr *DataReader) ReadUint64() (uint64, error) {
	var buf [8]byte
	if _, err := io.ReadFull(dr.br, buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf[:]), nil
}

// ReadInt64 reads an int64
func (dr *DataReader) ReadInt64() (int64, error) {
	x, err := dr.ReadUint64()
	return int64(x), err
}

// ReadUint32 reads a uint32
func (dr *DataReader) ReadUint32() (uint32, error) {
	var buf [4]byte
	if _, err := io.ReadFull(dr.br, buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf[:]), nil
}

// ReadInt32 reads an int32
func (dr *DataReader) ReadInt32() (int32, error) {
	x, err := dr.ReadUint32()
	return int32(x), err
}

// ReadUint16 reads a uint16
func (dr *DataReader) ReadUint16() (uint16, error) {
	var buf [2]byte
	if _, err := io.ReadFull(dr.br, buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(buf[:]), nil
}

// ReadInt16 reads an int16
func (dr *DataReader) ReadInt16() (int16, error) {
	x, err := dr.ReadUint16()
	return int16(x), err
}

// ReadUint8 reads a uint8
func (dr *DataReader) ReadUint8() (uint8, error) {
	return dr.br.ReadByte()
}

// ReadInt8 reads an int8
func (dr *DataReader) ReadInt8() (int8, error) {
	x, err := dr.ReadUint8()
	return int8(x), err
}

// ReadByte reads a byte
func (dr *DataReader) ReadByte() (byte, error) {
	return dr.br.ReadByte()
}

// ReadFloat64 reads a float64
func (dr *DataReader) ReadFloat64() (float64, error) {
	x, err := dr.ReadUint64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(x), nil
}

// ReadFloat32 reads a float32
func (dr *DataReader) ReadFloat32() (float32, error) {
	x, err := dr.ReadUint32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(x), nil
}

// ReadBool reads a bool
func (dr *DataReader) ReadBool() (bool, error) {
	b, err := dr.br.ReadByte()
	return b != 0, err
}

// ReadBytes reads bytes
func (dr *DataReader) ReadBytes() ([]byte, error) {
	n, err := dr.ReadUvarint()
	if err != nil {
		return nil, err
	}
	b := make([]byte, n)
	if _, err := io.ReadFull(dr.br, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (dr *DataReader) ReadBinary() ([]byte, error) {
	n, err := dr.ReadUint64()
	if err != nil {
		return nil, err
	}
	b := make([]byte, n)
	if _, err := io.ReadFull(dr.br, b); err != nil {
		return nil, err
	}
	return b, nil
}

// ReadString reads a string
func (dr *DataReader) ReadString() (string, error) {
	b, err := dr.ReadBytes()
	if err != nil {
		return "", err
	}
	return *(*string)(unsafe.Pointer(&b)), nil
}

// UnsafeBytesToString converts bytes to string saving allocations
func UnsafeBytesToString(bytes []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))

	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}))
}

// UnsafeStringToBytes converts bytes to string saving allocations by re-using
func UnsafeStringToBytes(s string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}))
}
