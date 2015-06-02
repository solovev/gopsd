package gopsd

import (
	"bytes"
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

func unpackRLEBits(data []int8, length int) []int8 {
	result := make([]int8, length)
	wPos, rPos := 0, 0
	for rPos < len(data) {
		n := data[rPos]
		rPos++
		if n > 0 {
			count := int(n) + 1
			for j := 0; j < count; j++ {
				result[wPos] = data[rPos]
				wPos++
				rPos++
			}
		} else {
			b := data[rPos]
			rPos++
			count := int(-n) + 1
			for j := 0; j < count; j++ {
				result[wPos] = b
				wPos++
			}
		}
	}
	return result
}

type StringMixer struct {
	Indent int
	buffer bytes.Buffer
}

func newStringMixer(indent int) *StringMixer {
	sm := new(StringMixer)
	sm.Indent = indent
	return sm
}

func (s *StringMixer) Add(values ...string) *StringMixer {
	for _, value := range values {
		for i := 0; i < s.Indent; i++ {
			s.buffer.WriteString("    ")
		}
		s.buffer.WriteString(value)
	}
	return s
}

func (s *StringMixer) NewLine() *StringMixer {
	s.buffer.WriteString("\n")
	return s
}

func (s *StringMixer) String() string {
	return s.buffer.String()
}
