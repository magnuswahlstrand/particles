package easing

func InCubic(t float64) float64  { return t * t * t }
func OutCubic(t float64) float64 { return 1 - (1-t)*(1-t)*(1-t) }
func Linear(t float64) float64   { return t }
