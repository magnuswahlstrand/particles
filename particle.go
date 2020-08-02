package particles

import (
	"github.com/peterhellberg/gfx"
	"image/color"
)

type Particle struct {
	pos          gfx.Vec
	velocity     gfx.Vec
	acceleration gfx.Vec

	currentLifetime float64
	startLifetime   float64

	color color.Color
}
