package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/particles"
	"github.com/kyeett/particles/easing"
	"github.com/kyeett/particles/generators"
	"github.com/kyeett/particles/modules/coloroverliftetime"
	"github.com/kyeett/particles/shapes"
	"golang.org/x/image/colornames"
	"math"
)

type Game struct {
	particles        *particles.ParticleSystem
	rateOverDistance float64
}

func (g *Game) Update(_ *ebiten.Image) error {
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
	initialLifetime = 0.5

	// Shape
	radius       = 50
	initialAngle = 60

	windowWidth  = 800
	windowHeight = 600
)

var (
	initialRate             = float64(0)
	initialRateOverDistance = float64(0)
)

func main() {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	g := &Game{
		particles: particles.NewParticleSystem(particles.Options{
			PositionX:        windowWidth/2,
			PositionY:        windowHeight/2,
			StartLifetime:    generators.FloatConstant{initialLifetime},
			StartSize:        generators.FloatConstant{0.3},
			StartSpeed:       generators.FloatConstant{5.0},
			Rate:             &initialRate,
			RateOverDistance: &initialRateOverDistance,
			Burst: &particles.Burst{
				Count:    100,
				Cycles:   5,
				Interval: 0.3,
			},

			Shape:            shapes.NewCone(toRad(initialAngle), float64(radius)),
			Material:         particles.MaterialDot,

			ColorOverLifetime: coloroverliftetime.ColorBetweenTwoConstants{colornames.Red, colornames.Yellow, easing.OutQuint},
		}),
	}

	ebiten.RunGame(g)
}
