package types

import "github.com/solovev/gopsd/util"

// Matrix stores information about 2d transformation
//  XX - Horizontal scale
//  XY - Horizontal incline
//  YX - Vertical incline
//  YY - Vertical scale
//  TX - Horizontal offset (pixels)
//  TY - Vertical offset (pixels)
type Matrix struct {
	XX, XY, YX, YY, TX, TY float64
}

func ReadMatrix(r *util.Reader) *Matrix {
	return &Matrix{r.ReadFloat64(), r.ReadFloat64(), r.ReadFloat64(), r.ReadFloat64(), r.ReadFloat64(), r.ReadFloat64()}
}
