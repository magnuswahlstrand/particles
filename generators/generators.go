package generators

import (
	"math/rand"
)

type Float interface {
	New() float64
}

var _ Float = &FloatConstant{}
var _ Float = &FloatRandomBetweenTwoConstants{}

type FloatConstant struct {
	Value float64
}

func (g FloatConstant) New() float64 {
	return g.Value
}

type FloatRandomBetweenTwoConstants struct {
	Min float64
	Max float64
}

func (g FloatRandomBetweenTwoConstants) New() float64 {
	return (g.Max-g.Min)*rand.Float64() + g.Min
}

