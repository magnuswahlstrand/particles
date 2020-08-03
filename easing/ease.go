package easing

import "math"

func InCubic(t float64) float64  { return t * t * t }
func OutCubic(t float64) float64 { return 1 - math.Pow(1-t, 3) }
func OutQuint(t float64) float64 { return math.Pow(t, 5) }
func Linear(t float64) float64   { return t }
