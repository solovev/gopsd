package gopsd

import (
	"encoding/binary"
	"errors"
	"os"
	"strings"
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

func inRange(i interface{}, min, max int) bool {
	val := getInteger(i)
	if val >= min && val <= max {
		return true
	}

	return false
}

func valueIs(i interface{}, numbers ...int) bool {
	val := getInteger(i)
	for n := range numbers {
		if val == numbers[n] {
			return true
		}
	}
	return false
}

func stringValueIs(value string, values ...string) bool {
	for i := range values {
		if value == values[i] {
			return true
		}
	}
	return false
}

func getInteger(unk interface{}) int {
	switch i := unk.(type) {
	case int32:
		return int(i)
	case int16:
		return int(i)
	case byte:
		return int(i)
	case int:
		return i
	default:
		return 0
	}
}
