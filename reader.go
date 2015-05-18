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

func (r *Reader) ReadString32() (string, error) {
	value := make([]byte, 4)
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		return "", err
	}
	return string(value), nil
}

func (r *Reader) ReadInt16() (int16, error) {
	var value int16
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

func (r *Reader) ReadInt32() (int32, error) {
	var value int32
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

func (r *Reader) ReadFloat32() (float32, error) {
	var value float32
	if err := binary.Read(r.buf, binary.BigEndian, &value); err != nil {
		return 0, err
	}
	return value, nil
}

func (r *Reader) ReadPascalString() (string, error) {
	var l byte
	if err := binary.Read(r.buf, binary.BigEndian, &l); err != nil {
		return "", err
	}
	if l%2 == 0 {
		l = 1
	}
	value := make([]byte, l)
	if err := binary.Read(r.buf, binary.BigEndian, value); err != nil {
		return "", err
	}
	return string(value), nil
}

func (r *Reader) Skip(n int32) error {
	if _, err := r.buf.Seek(int64(n), 1); err != nil {
		return err
	}
	return nil
}
