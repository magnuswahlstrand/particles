package generators

import (
	"github.com/peterhellberg/gfx"
	"image/color"
	"math/rand"
)

type Color interface {
	New() color.Color
}

var _ Color = &ColorRandomBetweenTwoConstants{}
var _ Color = &ColorConstant{}

type ColorConstant struct {
	Color0 color.Color
}

func (c ColorConstant) New() color.Color {
	return c.Color0
}

type ColorRandomBetweenTwoConstants struct {
	Color1, Color2 color.Color
}

func (g ColorRandomBetweenTwoConstants) New() color.Color {
	return gfx.LerpColors(g.Color1, g.Color2, rand.Float64())
}
