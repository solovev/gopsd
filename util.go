package gopsd

import (
	"encoding/binary"
	"errors"
	"os"
)

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
		return false, errors.New("Wrong document signature!")
	}

	b = make([]byte, 2)
	_, err = f.Read(b)
	if err != nil {
		return false, err
	}
	if binary.BigEndian.Uint16(b) != 1 {
		return false, errors.New("Wrong document version!")
	}
	return true, nil
}
