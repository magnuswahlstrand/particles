package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/particles"
	"github.com/kyeett/particles/generators"
	"github.com/kyeett/particles/modules/coloroverliftetime"
	"github.com/kyeett/particles/shapes"
	"github.com/peterhellberg/gfx"
	"golang.org/x/image/colornames"
	"math"
)

type Game struct {
	particles        *particles.ParticleSystem
	rateOverDistance float64
}

func (g *Game) Update(_ *ebiten.Image) error {
	cx, cy := ebiten.CursorPosition()
	if gfx.MathAbs(float64(cx)) < 10000 && gfx.MathAbs(float64(cy)) < 10000 {
		g.particles.Move(float64(cx), float64(cy))
	}

	dt := 1.0 / 60.0
	g.particles.Update(dt)
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowWidth, windowHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.particles.Draw(screen)
}

func toRad(a int32) float64 {
	return math.Pi * float64(a) / 180
}

const (
	initialLifetime = 1

	// Shape
	radius       = 0
	initialAngle = 0

	windowWidth  = 800
	windowHeight = 600
)

var (
	initialRate             = float64(0)
	initialRateOverDistance = float64(0.01)
)

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	g := &Game{
		particles: particles.NewParticleSystem(particles.Options{
			PositionX:        300,
			PositionY:        500,
			StartLifetime:    generators.FloatConstant{initialLifetime},
			StartSize:        generators.FloatConstant{0.1},
			StartSpeed:       generators.FloatConstant{0.0},
			Rate:             &initialRate,
			RateOverDistance: &initialRateOverDistance,
			Shape:            shapes.NewCone(toRad(initialAngle), float64(radius)),
			Material:         particles.MaterialHeart,

			ColorOverLifetime: coloroverliftetime.ColorBetweenTwoConstants{colornames.Red, colornames.Yellow, coloroverliftetime.Linear},
		}),
	}

	ebiten.RunGame(g)
}
