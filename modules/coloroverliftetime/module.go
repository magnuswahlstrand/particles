package coloroverliftetime

import (
	"github.com/peterhellberg/gfx"
	"image/color"
)

type Colorizer interface {
	Color(normalizedTime float64) color.Color
}

var _ Colorizer = &ColorBetweenTwoConstants{}
var _ Colorizer = &ColorConstant{}

func EaseInCubic(t float64) float64  { return t * t * t }
func EaseOutCubic(t float64) float64 { return 1 - (1-t)*(1-t)*(1-t) }
func Linear(t float64) float64       { return t }

type ColorBetweenTwoConstants struct {
	Color1, Color2 color.Color
	Easing         func(t float64) float64
}

type ColorConstant struct {
	color1 color.Color
}

func (c ColorConstant) Color(normalizedTime float64) color.Color {
	return c.color1
}

func (c ColorBetweenTwoConstants) Color(normalizedTime float64) color.Color {
	return gfx.LerpColors(c.Color1, c.Color2, c.Easing(normalizedTime))
}
