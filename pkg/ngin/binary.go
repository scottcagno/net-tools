package ngin

import (
	"bufio"
	"encoding/binary"
	"io"
	"math"
	"unsafe"
)

// BinaryWriter for writing
type BinaryWriter struct {
	bw *bufio.Writer
}

// NewBinaryWriter returns a new Writer
func NewBinaryWriter(w io.Writer) *BinaryWriter {
	return &BinaryWriter{bufio.NewWriter(w)}
}

// Flush the buffered bytes to the underlying writer
func (w *BinaryWriter) Flush() error {
	return w.bw.Flush()
}

// WriteUvarint writes a uvarint
func (w *BinaryWriter) WriteUvarint(x uint64) error {
	var buf [10]byte
	n := binary.PutUvarint(buf[:], x)
	_, err := w.bw.Write(buf[:n])
	return err
}

// WriteVarint writes a varint
func (w *BinaryWriter) WriteVarint(x int64) error {
	var buf [10]byte
	n := binary.PutVarint(buf[:], x)
	_, err := w.bw.Write(buf[:n])
	return err
}

// WriteUint32 writes a uint32
func (w *BinaryWriter) WriteUint32(x uint32) error {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], x)
	_, err := w.bw.Write(buf[:])
	return err
}

// WriteInt32 writes an int32
func (w *BinaryWriter) WriteInt32(x int32) error {
	return w.WriteUint32(uint32(x))
}

// WriteUint16 writes a uint16
func (w *BinaryWriter) WriteUint16(x uint16) error {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], x)
	_, err := w.bw.Write(buf[:])
	return err
}

// WriteInt16 writes an int16
func (w *BinaryWriter) WriteInt16(x int16) error {
	return w.WriteUint16(uint16(x))
}

// WriteUint8 writes a uint8
func (w *BinaryWriter) WriteUint8(x uint8) error {
	return w.bw.WriteByte(x)
}

// WriteInt8 writes an int8
func (w *BinaryWriter) WriteInt8(x int8) error {
	return w.WriteUint8(uint8(x))
}

// WriteByte writes a byte
func (w *BinaryWriter) WriteByte(x byte) error {
	return w.WriteUint8(uint8(x))
}

// WriteUint64 writes a uint64
func (w *BinaryWriter) WriteUint64(x uint64) error {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], x)
	_, err := w.bw.Write(buf[:])
	return err
}

// WriteInt64 writes an int64
func (w *BinaryWriter) WriteInt64(x int64) error {
	return w.WriteUint64(uint64(x))
}

// WriteFloat64 writes a float64
func (w *BinaryWriter) WriteFloat64(f float64) error {
	return w.WriteUint64(math.Float64bits(f))
}

// WriteFloat32 writes a float32
func (w *BinaryWriter) WriteFloat32(f float32) error {
	return w.WriteUint32(math.Float32bits(f))
}

// WriteBool writes a bool
func (w *BinaryWriter) WriteBool(t bool) error {
	if t {
		return w.bw.WriteByte(1)
	}
	return w.bw.WriteByte(0)
}

// WriteString writes a string
func (w *BinaryWriter) WriteString(s string) error {
	if err := w.WriteUvarint(uint64(len(s))); err != nil {
		return err
	}
	_, err := w.bw.WriteString(s)
	return err
}

// WriteBytes writes bytes
func (w *BinaryWriter) WriteBytes(b []byte) error {
	return w.WriteString(*(*string)(unsafe.Pointer(&b)))
}

// BinaryReader for reading
type BinaryReader struct {
	br *bufio.Reader
}

// NewBinaryReader returns a new Reader
func NewBinaryReader(r io.Reader) *BinaryReader {
	return &BinaryReader{bufio.NewReader(r)}
}

// ReadUvarint reads a uvarint
func (r *BinaryReader) ReadUvarint() (uint64, error) {
	return binary.ReadUvarint(r.br)
}

// ReadVarint reads a varint
func (r *BinaryReader) ReadVarint() (int64, error) {
	return binary.ReadVarint(r.br)
}

// ReadUint64 reads a uint64
func (r *BinaryReader) ReadUint64() (uint64, error) {
	var buf [8]byte
	if _, err := io.ReadFull(r.br, buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf[:]), nil
}

// ReadInt64 reads an int64
func (r *BinaryReader) ReadInt64() (int64, error) {
	x, err := r.ReadUint64()
	return int64(x), err
}

// ReadUint32 reads a uint32
func (r *BinaryReader) ReadUint32() (uint32, error) {
	var buf [4]byte
	if _, err := io.ReadFull(r.br, buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf[:]), nil
}

// ReadInt32 reads an int32
func (r *BinaryReader) ReadInt32() (int32, error) {
	x, err := r.ReadUint32()
	return int32(x), err
}

// ReadUint16 reads a uint16
func (r *BinaryReader) ReadUint16() (uint16, error) {
	var buf [2]byte
	if _, err := io.ReadFull(r.br, buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(buf[:]), nil
}

// ReadInt16 reads an int16
func (r *BinaryReader) ReadInt16() (int16, error) {
	x, err := r.ReadUint16()
	return int16(x), err
}

// ReadUint8 reads a uint8
func (r *BinaryReader) ReadUint8() (uint8, error) {
	return r.br.ReadByte()
}

// ReadInt8 reads an int8
func (r *BinaryReader) ReadInt8() (int8, error) {
	x, err := r.ReadUint8()
	return int8(x), err
}

// ReadByte reads a byte
func (r *BinaryReader) ReadByte() (byte, error) {
	return r.br.ReadByte()
}

// ReadFloat64 reads a float64
func (r *BinaryReader) ReadFloat64() (float64, error) {
	x, err := r.ReadUint64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(x), nil
}

// ReadFloat32 reads a float32
func (r *BinaryReader) ReadFloat32() (float32, error) {
	x, err := r.ReadUint32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(x), nil
}

// ReadBool reads a bool
func (r *BinaryReader) ReadBool() (bool, error) {
	b, err := r.br.ReadByte()
	return b != 0, err
}

// ReadBytes reads bytes
func (r *BinaryReader) ReadBytes() ([]byte, error) {
	n, err := r.ReadUvarint()
	if err != nil {
		return nil, err
	}
	b := make([]byte, n)
	if _, err := io.ReadFull(r.br, b); err != nil {
		return nil, err
	}
	return b, nil
}

// ReadString reads a string
func (r *BinaryReader) ReadString() (string, error) {
	b, err := r.ReadBytes()
	if err != nil {
		return "", err
	}
	return *(*string)(unsafe.Pointer(&b)), nil
}
