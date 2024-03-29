package main

import (
	"github.com/hajimehoshi/ebiten"
	imgui "github.com/kyeett/guigi"
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

type UI struct {
	lifetime float64
	size     float64
	speed    float64

	colorOverLifetime ColorOverLifetime
	sizeOverLifeTime  SizeOverLifetime

	shape coneShape
}

type Game struct {
	particles *particles.ParticleSystem

	ui UI
}

type coneShape struct {
	angle  float64
	radius float64
}

func (g *Game) Update(_ *ebiten.Image) error {
	dt := 1.0 / 60.0
	g.particles.Update(dt)
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowWidth, windowHeight
}

func floatPtr(v float64) *float64 {
	return &v
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.particles.Draw(screen)

	imgui.BeginFrame(10, 10)
	if imgui.CollapsingHeader("Examples") {
		imgui.BeginListBox("ExamplesBox")
		if imgui.Selectable("pulsating dot") {
			g.particles, g.ui = dotExample()
		}
		if imgui.Selectable("fire") {
			g.particles, g.ui = fireExample()
		}
		imgui.EndListBox()
	}

	if imgui.DragFloatV("Lifetime", &g.ui.lifetime, 0.1, floatPtr(0), floatPtr(10)) {
		g.particles.StartLifetime = generators.FloatConstant{g.ui.lifetime}
	}
	if imgui.DragFloatV("StartSize", &g.ui.size, 0.1, floatPtr(0), floatPtr(3.0)) {
		g.particles.StartSize = generators.FloatConstant{g.ui.size}
	}
	if imgui.DragFloatV("StartSpeed", &g.ui.speed, 0.1, floatPtr(0), floatPtr(5)) {
		g.particles.StartSpeed = generators.FloatConstant{g.ui.speed}
	}

	if imgui.CollapsingHeader("Emission") {
		imgui.DragFloatV("Rate", &g.particles.Rate, 10, floatPtr(0), floatPtr(1000))
	}

	// Shape
	if imgui.CollapsingHeader("Shape") {
		angleChanged := imgui.DragFloatV("Angle", &g.ui.shape.angle, 1, floatPtr(0), floatPtr(90))
		radiusChanged := imgui.DragFloatV("Radius", &g.ui.shape.radius, 2, floatPtr(0), floatPtr(200))
		if angleChanged || radiusChanged {
			g.particles.Shape = shapes.NewCone(toRad(g.ui.shape.angle), g.ui.shape.radius)
		}
	}

	if imgui.CollapsingHeader("Rendering") {
		imgui.BeginListBox("Material")
		if imgui.Selectable("star") {
			g.particles.Material = particles.MaterialStar
		}
		if imgui.Selectable("heart") {
			g.particles.Material = particles.MaterialHeart
		}
		if imgui.Selectable("dot") {
			g.particles.Material = particles.MaterialDot
		}
		if imgui.Selectable("leaf") {
			g.particles.Material = particles.MaterialLeaf
		}
		imgui.EndListBox()
	}

	imgui.EndFrame(screen)
}

func newSizeOverLifetime(v1, v2 float32) sizeoverliftetime.SizeBetweenTwoConstants {
	return sizeoverliftetime.SizeBetweenTwoConstants{float64(v1), float64(v2), easing.Linear}
}

func newColorOverLifetime(c1, c2 [4]float32) coloroverliftetime.ColorBetweenTwoConstants {
	return coloroverliftetime.ColorBetweenTwoConstants{
		Color1: colorFromArray(c1),
		Color2: colorFromArray(c2),
		Easing: easing.OutCubic,
	}
}

func toRad(a float64) float64 {
	return math.Pi * a / 180
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
	ebiten.SetWindowSize(windowWidth, windowHeight)

	g := NewGame(initialLifetime, initialSize, initialSpeed, initialRate, initialAngle, radius, 0)
	ebiten.RunGame(g)
}

func NewGame(lifetime, size, speed, rate float64, angle float64, radius, gravity float64) *Game {
	ps := particles.NewParticleSystem(particles.Options{
		PositionX:     550,
		PositionY:     400,
		StartLifetime: generators.FloatConstant{lifetime},
		StartSize:     generators.FloatConstant{size},
		StartSpeed:    generators.FloatConstant{speed},
		Rate:          &rate,
		Shape:         shapes.NewCone(toRad(angle), radius),
		Gravity:       gfx.V(0, gravity),
	})

	g := &Game{
		particles: ps,

		ui: UI{
			lifetime: lifetime,
			size:     size,
			speed:    speed,

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
				angle:  angle,
				radius: radius,
			},
		},
	}
	return g
}
