package main

import "github.com/kyeett/particles"

func dotExample() (*particles.ParticleSystem, UI) {
	g := NewGame(2.5, 0.9, 0, 1.3, 0, 0, 0)
	//g.ui.exampleIndex = 0

	g.particles.Material = particles.MaterialDot
	//g.ui.materialIndex = 2

	c1 := [4]float32{1, 1, 1, 1}
	c2 := [4]float32{1, 1, 1, 0}
	g.particles.ColorOverLifetime = newColorOverLifetime(c1, c2)
	g.ui.colorOverLifetime = ColorOverLifetime{
		enabled:    true,
		startColor: c1,
		endColor:   c2,
	}

	g.particles.SizeOverLifetime = newSizeOverLifetime(0, 1.5)
	g.ui.sizeOverLifeTime = SizeOverLifetime{
		enabled: true,
		start:   0,
		end:     1.8,
	}
	return g.particles, g.ui
}

func fireExample() (*particles.ParticleSystem, UI) {
	g := NewGame(2, 0.2, 5, 700, 8, 75, 2.5)
	//g.ui.exampleIndex = 1

	g.particles.Material = particles.MaterialDot
	//g.ui.materialIndex = 2

	c1 := [4]float32{1, 0.25, 0, 1}
	c2 := [4]float32{1, 1, 0, 0}
	g.particles.ColorOverLifetime = newColorOverLifetime(c1, c2)
	g.ui.colorOverLifetime = ColorOverLifetime{
		enabled:    true,
		startColor: c1,
		endColor:   c2,
	}

	return g.particles, g.ui
}
