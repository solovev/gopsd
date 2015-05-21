package gopsd

type Rectangle struct {
	Top, Left, Bottom, Right int32
}

func (r *Rectangle) X() int32 {
	return r.Left - r.Right
}

func (r *Rectangle) Y() int32 {
	return r.Bottom - r.Top
}
