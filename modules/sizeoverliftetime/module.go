package sizeoverliftetime

import (
	"github.com/peterhellberg/gfx"
)

type Sizer interface {
	Size(normalizedTime float64) float64
}

var _ Sizer = &SizeBetweenTwoConstants{}
var _ Sizer = &SizeConstant{}

type SizeBetweenTwoConstants struct {
	Size1, Size2 float64
	Easing       func(t float64) float64
}

type SizeConstant struct {
	Size0 float64
}

func (c SizeConstant) Size(_ float64) float64 {
	return c.Size0
}

func (c SizeBetweenTwoConstants) Size(normalizedTime float64) float64 {
	return gfx.Lerp(c.Size1, c.Size2, c.Easing(normalizedTime))
}
