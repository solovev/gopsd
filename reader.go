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
