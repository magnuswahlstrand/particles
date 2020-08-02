package shapes

import (
	"github.com/peterhellberg/gfx"
	"math/rand"
)

type Shape interface {
	New() (float64, float64, float64)
}

var _ Shape = &Cone{}

type Cone struct {
	angle  float64
	radius float64
}

func NewCone(angle, radius float64) Cone {
	return Cone{angle, radius}
}

func (c Cone) New() (float64, float64, float64) {
	rnd := rand.Float64()
	x := 2 * c.radius * (rnd - 0.5)
	angle := gfx.Lerp(-c.angle, c.angle, rnd)
	return x, 0, angle
}
