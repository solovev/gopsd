package util

import "bytes"

type StringMixer struct {
	buffer bytes.Buffer
}

func (s *StringMixer) Add(values ...string) *StringMixer {
	for _, value := range values {
		s.buffer.WriteString(value)
	}
	return s
}

func (s *StringMixer) AddIndent(value int) *StringMixer {
	for i := 0; i < value; i++ {
		s.buffer.WriteString("    ")
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
