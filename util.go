package gopsd

import (
	"encoding/binary"
	"errors"
	"os"
	"strings"
)

type Rectangle struct {
	Top, Left, Bottom, Right int32
}

func (r *Rectangle) X() int32 {
	return r.Left - r.Right
}

func (r *Rectangle) Y() int32 {
	return r.Bottom - r.Top
}

func IsValid(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	b := make([]byte, 4)
	_, err = f.Read(b)
	if err != nil {
		return false, err
	}
	if string(b) != "8BPS" {
		return false, errors.New("Wrong document signature.")
	}

	b = make([]byte, 2)
	_, err = f.Read(b)
	if err != nil {
		return false, err
	}
	ver := binary.BigEndian.Uint16(b)
	if (strings.HasSuffix(path, "psd") && ver != 1) ||
		(strings.HasSuffix(path, "psb") && ver != 2) {
		return false, errors.New("Wrong document version.")
	}
	return true, nil
}
