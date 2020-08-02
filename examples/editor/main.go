package main

import (
	"github.com/gabstv/ebiten-imgui/renderer"
	"github.com/hajimehoshi/ebiten"
	"github.com/inkyblackness/imgui-go/v2"
	"github.com/kyeett/particles"
	"github.com/kyeett/particles/easing"
	"github.com/kyeett/particles/generators"
	"github.com/kyeett/particles/modules/coloroverliftetime"
	"github.com/kyeett/particles/modules/sizeoverliftetime"
	"github.com/kyeett/particles/shapes"
	"github.com/peterhellberg/gfx"
	"image/color"
	"math"
)

type Game struct {
	manager   *renderer.Manager
	particles *particles.ParticleSystem

	color [4]float32
	rate  int32

	size     float32
	lifetime float32
	speed    float32
	gravity  float32

	colorOverLifetime ColorOverLifetime
	sizeOverLifeTime  SizeOverLifetime

	materialIndex int32
	shape         coneShape
}

type coneShape struct {
	angle  int32
	radius int32
}

func (g *Game) Update(_ *ebiten.Image) error {
	dt := 1.0 / 60.0
	g.manager.Update(float32(dt), 800, 600)
	g.particles.Update(dt)
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowWidth, windowHeight
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.particles.Draw(screen)

	g.manager.BeginFrame()
	{
		// General
		imgui.Text("General")
		if imgui.SliderFloat("Lifetime", &g.lifetime, 0.0, 5.0) {
			g.particles.StartLifetime = generators.FloatConstant{float64(g.lifetime)}
		}
		if imgui.SliderFloat("Size", &g.size, 0.0, 1.0) {
			g.particles.StartSize = generators.FloatConstant{float64(g.size)}
		}
		if imgui.SliderFloat("Speed", &g.speed, 0.0, 10.0) {
			g.particles.StartSpeed = generators.FloatConstant{float64(g.speed)}
		}

		if imgui.ColorEdit4("StartColor", &g.color) {
			g.particles.Color = generators.ColorConstant{colorFromArray(g.color)}
		}

		if imgui.SliderFloat("Gravity", &g.gravity, 0.0, 10.0) {
			g.particles.Gravity = gfx.V(0, float64(g.gravity))
		}

		// Emission
		imgui.Text("")
		imgui.Text("Emission")
		if imgui.SliderInt("Rate", &g.rate, 0, 1000) {
			g.particles.Rate = float64(g.rate)
		}

		// Shape
		imgui.Text("")
		imgui.Text("Shape (Cone)")
		angleChanged := imgui.SliderInt("Angle", &g.shape.angle, 0, 180)
		radiusChanged := imgui.SliderInt("Radius", &g.shape.radius, 0, 200)
		if angleChanged || radiusChanged {
			g.particles.Shape = shapes.NewCone(toRad(g.shape.angle), float64(g.shape.radius))
		}

		// Modules
		imgui.Text("")
		imgui.Text("Modules")

		// Module: Color over lifetime
		if imgui.Checkbox("ColorOverLifetime", &g.colorOverLifetime.enabled) {
			if g.colorOverLifetime.enabled {
				g.particles.ColorOverLifetime = newColorOverLifetime(g.colorOverLifetime.startColor, g.colorOverLifetime.endColor)
			} else {
				g.particles.ColorOverLifetime = nil
			}
		}
		if g.colorOverLifetime.enabled {
			c1Changed := imgui.ColorEdit4("Start", &g.colorOverLifetime.startColor)
			c2Changed := imgui.ColorEdit4("End", &g.colorOverLifetime.endColor)
			if c1Changed || c2Changed {
				g.particles.ColorOverLifetime = newColorOverLifetime(g.colorOverLifetime.startColor, g.colorOverLifetime.endColor)
			}
		}

		// Module: Size over lifetime
		if imgui.Checkbox("SizeOverLifetime", &g.sizeOverLifeTime.enabled) {
			if g.sizeOverLifeTime.enabled {
				g.particles.SizeOverLifetime = newSizeOverLifetime(g.sizeOverLifeTime.start, g.sizeOverLifeTime.end)
			} else {
				g.particles.SizeOverLifetime = nil
			}
		}
		if g.sizeOverLifeTime.enabled {
			v1Changed := imgui.SliderFloat("Start", &g.sizeOverLifeTime.start, 0, 2)
			v2Changed := imgui.SliderFloat("End", &g.sizeOverLifeTime.end, 0, 2)
			if v1Changed || v2Changed {
				g.particles.SizeOverLifetime = newSizeOverLifetime(g.sizeOverLifeTime.start, g.sizeOverLifeTime.end)
			}
		}

		// Rendering
		imgui.Text("")
		imgui.Text("Rendering")

		if imgui.ListBox("Material", &g.materialIndex, []string{"star", "heart", "dot"}) {
			switch g.materialIndex {
			case 0:
				g.particles.Material = particles.MaterialStar
			case 1:
				g.particles.Material = particles.MaterialHeart
			case 2:
				g.particles.Material = particles.MaterialDot
			}
		}
	}

	g.manager.EndFrame(screen)
}

func newSizeOverLifetime(v1, v2 float32) sizeoverliftetime.SizeBetweenTwoConstants {
	return sizeoverliftetime.SizeBetweenTwoConstants{float64(v1), float64(v2), easing.Linear}
}

func newColorOverLifetime(c1, c2 [4]float32) coloroverliftetime.ColorBetweenTwoConstants {
	return coloroverliftetime.ColorBetweenTwoConstants{
		Color1: colorFromArray(c1),
		Color2: colorFromArray(c2),
		Easing: easing.Linear,
	}
}

func toRad(a int32) float64 {
	return math.Pi * float64(a) / 180
}

func colorFromArray(clr [4]float32) color.Color {
	return color.RGBA{uint8(clr[0] * 255), uint8(clr[1] * 255), uint8(clr[2] * 255), uint8(clr[3] * 255)}
}

type ColorOverLifetime struct {
	enabled    bool
	startColor [4]float32
	endColor   [4]float32
}

type SizeOverLifetime struct {
	enabled bool
	start   float32
	end     float32
}

const (
	initialSize     = 0.1
	initialSpeed    = 2.0
	initialLifetime = 2

	// Shape
	radius       = 50
	initialAngle = 45

	windowWidth  = 800
	windowHeight = 600

)

var (
	initialRate = float64(100)
)

func main() {
	mgr := renderer.New(nil)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	g := &Game{
		manager: mgr,
		particles: particles.NewParticleSystem(particles.Options{
			PositionX:     550,
			PositionY:     400,
			StartLifetime: generators.FloatConstant{initialLifetime},
			StartSize:     generators.FloatConstant{initialSize},
			StartSpeed:    generators.FloatConstant{initialSpeed},
			Rate:          &initialRate,
			Shape:         shapes.NewCone(toRad(initialAngle), float64(radius)),
			Gravity:       gfx.Vec{},
		}),

		color:    [4]float32{1, 1, 1, 1},
		lifetime: initialLifetime,
		size:     initialSize,
		speed:    initialSpeed,
		rate:     int32(initialRate),

		colorOverLifetime: ColorOverLifetime{
			enabled:    false,
			startColor: [4]float32{1, 0, 0, 1},
			endColor:   [4]float32{1, 1, 0, 0.2},
		},

		sizeOverLifeTime: SizeOverLifetime{
			enabled: false,
			start:   1.0,
			end:     0.0,
		},

		shape: coneShape{
			angle:  initialAngle,
			radius: radius,
		},

		materialIndex: 0,
	}

	ebiten.RunGame(g)
}
