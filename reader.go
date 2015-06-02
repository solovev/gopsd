package gopsd

import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"
)

type Reader struct {
	buf      *bytes.Reader
	Position int // [CHECK] Must be int64?
}

func NewReader(b []byte) *Reader {
	return &Reader{bytes.NewReader(b), 0}
}

func (r *Reader) ReadByte() byte {
	var value byte
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		panic(err)
	}
	r.Position++
	return value
}

func (r *Reader) ReadString(n int) string {
	value := make([]byte, n)
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		panic(err)
	}
	r.Position += n
	return string(value)
}

func (r *Reader) ReadInt16() int16 {
	var value int16
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		panic(err)
	}
	r.Position += 2
	return value
}

func (r *Reader) ReadInt32() int32 {
	var value int32
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		panic(err)
	}
	r.Position += 4
	return value
}

func (r *Reader) ReadInt64() int64 {
	var value int64
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		panic(err)
	}
	r.Position += 8
	return value
}

func (r *Reader) ReadFloat32() float32 {
	var value float32
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		panic(err)
	}
	r.Position += 4
	return value
}

func (r *Reader) ReadFloat64() float64 {
	var value float64
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		panic(err)
	}
	r.Position += 8
	return value
}

func (r *Reader) ReadPascalString() string {
	var length byte
	if err := binary.Read(r.buf, binary.BigEndian, &length); err != nil {
		panic(err)
	}
	if length == 0 {
		length = 1
	}
	r.Position += 1
	return r.ReadString(int(length))
}

func (r *Reader) ReadUnicodeString() string {
	n := reader.ReadInt32()
	array := make([]uint16, n)
	for i := range array {
		if err := binary.Read(r.buf, binary.BigEndian, &array[i]); err != nil {
			panic(err)
		}
		r.Position += 2
	}
	return string(utf16.Decode(array))
}

func (r *Reader) ReadDynamicString() string {
	length := int(reader.ReadInt32())
	if length == 0 {
		length = 4
	}
	return reader.ReadString(length)
}

func (r *Reader) ReadBytes(number interface{}) []byte {
	n := getInteger(number)
	value := make([]byte, n)
	if err := binary.Read(r.buf, binary.BigEndian, value); err != nil {
		panic(err)
	}
	r.Position += n
	return value
}

func (r *Reader) ReadSignedBytes(number interface{}) []int8 {
	n := getInteger(number)
	value := make([]int8, n)
	if err := binary.Read(r.buf, binary.BigEndian, value); err != nil {
		panic(err)
	}
	r.Position += n
	return value
}

func (r *Reader) Skip(number interface{}) {
	n := getInteger(number)
	r.Position += n
	if _, err := r.buf.Seek(int64(n), 1); err != nil {
		panic(err)
	}
}
