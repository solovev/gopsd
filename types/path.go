package types

import "github.com/solovev/gopsd/util"

type Point struct {
	X, Y float32
}

type Knot struct {
	Controls []*Point
	Anchor   *Point
}

type Path struct {
	IsOpen              bool
	StartsWithAllPixels bool
	Knots               []*Knot
}

// TODO: If windows - reverse byte order?
func ReadPath(data []byte) *Path {
	r := util.NewReader(data)
	path := new(Path)
	index := 0
	for r.Position < len(data) {
		record := r.ReadInt16()
		if len(data)-r.Position >= 24 {
			switch record {
			case 0, 3:
				path.Knots = make([]*Knot, r.ReadInt16())
				r.Skip(22)
			case 1, 2, 4, 5:
				if len(path.Knots) == 0 {
					return nil
				}
				path.Knots[index] = readKnot(r)
				index++
			case 6: // Path fill
				r.Skip(24)
			case 7: // Clipboard
				r.Skip(24)
			case 8: // Initial fill
				path.StartsWithAllPixels = r.ReadInt16() == 1
				r.Skip(22)
			}
		}
	}
	return path
}

func readKnot(r *util.Reader) *Knot {
	knot := new(Knot)
	knot.Controls = make([]*Point, 2)

	knot.Controls[0] = readPoint(r)
	knot.Anchor = readPoint(r)
	knot.Controls[1] = readPoint(r)

	return knot
}

func readPoint(r *util.Reader) *Point {
	point := new(Point)
	point.Y = readComponent(r)
	point.X = readComponent(r)
	return point
}

func readComponent(r *util.Reader) float32 {
	i := float32(r.ReadByte())
	f := float32(r.ReadInt24()) / 16777216.0
	return i + f
}
