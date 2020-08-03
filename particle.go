package particles

import (
	"github.com/peterhellberg/gfx"
	"image/color"
)

type Particle struct {
	pos          gfx.Vec
	velocity     gfx.Vec
	acceleration gfx.Vec

	rotation        float64
	angularVelocity float64

	currentLifetime float64
	startLifetime   float64
	startSize       float64

	color color.Color
}

func (p Particle) normalizedLifetime() float64 {
	return 1 - p.currentLifetime/p.startLifetime
}
