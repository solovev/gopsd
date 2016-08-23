package types

import (
	"fmt"

	"github.com/solovev/gopsd/util"
)

type Rectangle struct {
	top    int32
	left   int32
	bottom int32
	right  int32

	X, Y, Width, Height int32
}

type RectangleFloat struct {
	Top, Left, Bottom, Right float64
}

func CreateRectangle(x, y, width, height int32) *Rectangle {
	r := new(Rectangle)

	r.X = x
	r.Y = y
	r.Width = width
	r.Height = height

	return r
}

func NewRectangle(reader *util.Reader) *Rectangle {
	r := new(Rectangle)

	r.top = reader.ReadInt32()
	r.Y = r.top

	r.left = reader.ReadInt32()
	r.X = r.left

	r.bottom = reader.ReadInt32()
	r.right = reader.ReadInt32()

	r.Width = r.right - r.left
	r.Height = r.bottom - r.top

	return r
}

func (r Rectangle) ToString() string {
	return fmt.Sprintf("[X: %d, Y: %d, Width: %d, Height: %d]", r.X, r.Y, r.Width, r.Height)
}
