package types

import "github.com/solovev/gopsd/util"

type Color struct {
	red, green, blue, alpha int16
}

func (c Color) Red() int16 {
	return c.red
}

func (c Color) Green() int16 {
	return c.green
}

func (c Color) Blue() int16 {
	return c.blue
}

func (c Color) Alpha() int16 {
	return c.alpha
}

func NewRGBAColor(reader *util.Reader) *Color {
	return &Color{reader.ReadInt16(), reader.ReadInt16(), reader.ReadInt16(), reader.ReadInt16()}
}
