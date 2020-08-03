package rotationoverlifetime

import (
	"github.com/peterhellberg/gfx"
)

type RotationOverLifetime struct {
	angularVelocity float64
}

type Rotator interface {
	Rotation(normalizedTime float64) float64
}

var _ Rotator = &RotatorBetweenTwoConstants{}
var _ Rotator = &RotatorConstant{}

type RotatorBetweenTwoConstants struct {
	W1, W2 float64
	Easing func(t float64) float64
}

func (r RotatorBetweenTwoConstants) Rotation(normalizedTime float64) float64 {
	return gfx.Lerp(r.W1, r.W2, r.Easing(normalizedTime))
}

type RotatorConstant struct {
	W1 float64
}

func (r RotatorConstant) Rotation(_ float64) float64 {
	return r.W1
}
