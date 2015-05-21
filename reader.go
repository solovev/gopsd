package gopsd

import (
	"bytes"
	"encoding/binary"
)

type Reader struct {
	buf *bytes.Reader
}

func NewReader(b []byte) *Reader {
	return &Reader{bytes.NewReader(b)}
}

func (r Reader) ReadByte() byte {
	var value byte
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		panic(err)
	}
	return value
}

func (r Reader) ReadString32() string {
	value := make([]byte, 4)
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		panic(err)
	}
	return string(value)
}

func (r Reader) ReadInt16() int16 {
	var value int16
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		panic(err)
	}
	return value
}

func (r Reader) ReadInt32() int32 {
	var value int32
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		panic(err)
	}
	return value
}

func (r Reader) ReadFloat32() float32 {
	var value float32
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		panic(err)
	}
	return value
}

func (r Reader) ReadFloat64() float64 {
	var value float64
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		panic(err)
	}
	return value
}

func (r Reader) ReadPascalString() string {
	var l byte
	if err := binary.Read(r.buf, binary.BigEndian, &l); err != nil {
		panic(err)
	}
	if l%2 == 0 {
		l = 1
	}
	value := make([]byte, l)
	if err := binary.Read(r.buf, binary.BigEndian, value); err != nil {
		panic(err)
	}
	return string(value)
}

func (r Reader) ReadRectangle() *Rectangle {
	return &Rectangle{r.ReadInt32(), r.ReadInt32(), r.ReadInt32(), r.ReadInt32()}
}

func (r Reader) Skip(n int32) {
	if _, err := r.buf.Seek(int64(n), 1); err != nil {
		panic(err)
	}
}
