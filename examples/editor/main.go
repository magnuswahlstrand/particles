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

type UI struct {
	color [4]float32
	rate  float32

	size     float32
	lifetime float32
	speed    float32
	gravity  float32

	colorOverLifetime ColorOverLifetime
	sizeOverLifeTime  SizeOverLifetime

	exampleIndex  int32
	materialIndex int32
	shape         coneShape
}

type Game struct {
	manager   *renderer.Manager
	particles *particles.ParticleSystem

	ui UI
}

type coneShape struct {
	angle  int32
	radius int32
}

func (g *Game) Update(_ *ebiten.Image) error {
	dt := 1.0 / 60.0
	g.manager.Update(float32(dt), windowWidth, windowHeight)
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

		if imgui.CollapsingHeader("Examples") {
			if imgui.ListBox("", &g.ui.exampleIndex, []string{"pulsating dot", "fire"}) {
				switch g.ui.exampleIndex {
				case 0:
					g.particles, g.ui = dotExample()
				case 1:
					g.particles, g.ui = fireExample()
				}
			}
		}

		imgui.Spacing()
		imgui.Spacing()

		// General

		if imgui.SliderFloat("Lifetime", &g.ui.lifetime, 0.0, 5.0) {
			g.particles.StartLifetime = generators.FloatConstant{float64(g.ui.lifetime)}
		}
		if imgui.SliderFloat("Size", &g.ui.size, 0.0, 1.0) {
			g.particles.StartSize = generators.FloatConstant{float64(g.ui.size)}
		}
		if imgui.SliderFloat("Speed", &g.ui.speed, 0.0, 10.0) {
			g.particles.StartSpeed = generators.FloatConstant{float64(g.ui.speed)}
		}

		if imgui.ColorEdit4("StartColor", &g.ui.color) {
			g.particles.Color = generators.ColorConstant{colorFromArray(g.ui.color)}
		}

		if imgui.SliderFloat("Gravity", &g.ui.gravity, 0.0, 10.0) {
			g.particles.Gravity = gfx.V(0, float64(g.ui.gravity))
		}

		imgui.Spacing()
		imgui.Spacing()

		// Emission
		if imgui.CollapsingHeader("Emission") {
			if imgui.SliderFloat("Rate", &g.ui.rate, 0, 1000) {
				g.particles.Rate = float64(g.ui.rate)
			}
		}

		// Shape
		if imgui.CollapsingHeader("Shape") {
			angleChanged := imgui.SliderInt("Angle", &g.ui.shape.angle, 0, 180)
			radiusChanged := imgui.SliderInt("Radius", &g.ui.shape.radius, 0, 200)
			if angleChanged || radiusChanged {
				g.particles.Shape = shapes.NewCone(toRad(g.ui.shape.angle), float64(g.ui.shape.radius))
			}
		}

		// Rendering
		if imgui.CollapsingHeader("Rendering") {
			if imgui.ListBox("Material", &g.ui.materialIndex, []string{"star", "heart", "dot", "leaf"}) {
				switch g.ui.materialIndex {
				case 0:
					g.particles.Material = particles.MaterialStar
				case 1:
					g.particles.Material = particles.MaterialHeart
				case 2:
					g.particles.Material = particles.MaterialDot
				case 3:
					g.particles.Material = particles.MaterialLeaf
				}
			}
		}

		// Modules
		imgui.Text("")
		imgui.Text("Modules")

		// Module: Color over lifetime
		if imgui.Checkbox("ColorOverLifetime", &g.ui.colorOverLifetime.enabled) {
			if g.ui.colorOverLifetime.enabled {
				g.particles.ColorOverLifetime = newColorOverLifetime(g.ui.colorOverLifetime.startColor, g.ui.colorOverLifetime.endColor)
			} else {
				g.particles.ColorOverLifetime = nil
			}
		}
		if g.ui.colorOverLifetime.enabled {
			c1Changed := imgui.ColorEdit4("Start", &g.ui.colorOverLifetime.startColor)
			c2Changed := imgui.ColorEdit4("End", &g.ui.colorOverLifetime.endColor)
			if c1Changed || c2Changed {
				g.particles.ColorOverLifetime = newColorOverLifetime(g.ui.colorOverLifetime.startColor, g.ui.colorOverLifetime.endColor)
			}
		}

		// Module: Size over lifetime
		if imgui.Checkbox("SizeOverLifetime", &g.ui.sizeOverLifeTime.enabled) {
			if g.ui.sizeOverLifeTime.enabled {
				g.particles.SizeOverLifetime = newSizeOverLifetime(g.ui.sizeOverLifeTime.start, g.ui.sizeOverLifeTime.end)
			} else {
				g.particles.SizeOverLifetime = nil
			}
		}
		if g.ui.sizeOverLifeTime.enabled {
			v1Changed := imgui.SliderFloat("Start", &g.ui.sizeOverLifeTime.start, 0, 2)
			v2Changed := imgui.SliderFloat("End", &g.ui.sizeOverLifeTime.end, 0, 2)
			if v1Changed || v2Changed {
				g.particles.SizeOverLifetime = newSizeOverLifetime(g.ui.sizeOverLifeTime.start, g.ui.sizeOverLifeTime.end)
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
		Easing: easing.OutCubic,
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

//const (
//	initialSize     = 0.1
//	initialSpeed    = 2.0
//	initialLifetime = 2
//
//	// Shape
//	radius       = 50
//	initialAngle = 45
//
//	windowWidth  = 800
//	windowHeight = 600
//)

//var (
//	initialRate = float64(100)
//)

const (
	initialSize     = 0.1
	initialSpeed    = 3.5
	initialLifetime = 2

	initialRate = 5

	// Shape
	radius       = 30
	initialAngle = 0

	windowWidth  = 800
	windowHeight = 600
)

func main() {
	mgr := renderer.New(nil)
	ebiten.SetWindowSize(windowWidth, windowHeight)

	g := NewGame(initialLifetime, initialSize, initialSpeed, initialRate, initialAngle, radius, 0)
	g.manager = mgr

	ebiten.RunGame(g)
}

func NewGame(lifetime, size, speed, rate float64, angle int32, radius, gravity float64) *Game {
	g := &Game{
		particles: particles.NewParticleSystem(particles.Options{
			PositionX:     550,
			PositionY:     400,
			StartLifetime: generators.FloatConstant{lifetime},
			StartSize:     generators.FloatConstant{size},
			StartSpeed:    generators.FloatConstant{speed},
			Rate:          &rate,
			Shape:         shapes.NewCone(toRad(angle), float64(radius)),
			Gravity:       gfx.V(0, gravity),
		}),

		ui: UI{
			color:    [4]float32{1, 1, 1, 1},
			lifetime: float32(lifetime),
			size:     float32(size),
			speed:    float32(speed),
			rate:     float32(rate),
			gravity:  float32(gravity),

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
				radius: int32(radius),
			},

			materialIndex: 0,
			exampleIndex: -1,
		},
	}
	return g
}
