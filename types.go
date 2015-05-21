package gopsd

type Rectangle struct {
	top, left, bottom, right int32
}

func (r Rectangle) X() int32 {
	return r.left
}

func (r Rectangle) Y() int32 {
	return r.top
}

func (r Rectangle) Width() int32 {
	return r.right - r.left
}

func (r Rectangle) Height() int32 {
	return r.bottom - r.top
}
