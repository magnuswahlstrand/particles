package particles

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/particles/assets"
	"github.com/kyeett/particles/generators"
	"github.com/kyeett/particles/modules/coloroverliftetime"
	"github.com/kyeett/particles/modules/rotationoverlifetime"
	"github.com/kyeett/particles/modules/sizeoverliftetime"
	"github.com/kyeett/particles/shapes"
	"github.com/peterhellberg/gfx"
	"golang.org/x/image/colornames"
	"image/color"
	"log"
	"math"
)

var (
	MaterialStar  *ebiten.Image
	MaterialHeart *ebiten.Image
	MaterialDot   *ebiten.Image
	MaterialLeaf  *ebiten.Image
)

func init() {
	img, err := gfx.DecodeImageBytes(assets.MustAsset("star.png"))
	if err != nil {
		log.Fatal(err)
	}
	MaterialStar, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	img, err = gfx.DecodeImageBytes(assets.MustAsset("heart.png"))
	if err != nil {
		log.Fatal(err)
	}
	MaterialHeart, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	img, err = gfx.DecodeImageBytes(assets.MustAsset("dot.png"))
	if err != nil {
		log.Fatal(err)
	}
	MaterialDot, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	img, err = gfx.DecodeImageBytes(assets.MustAsset("leaf.png"))
	if err != nil {
		log.Fatal(err)
	}
	MaterialLeaf, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
}

type ParticleSystem struct {
	particles []*Particle

	pos        gfx.Vec
	initialPos gfx.Vec

	Gravity gfx.Vec

	Material *ebiten.Image

	StartLifetime generators.Float
	StartSpeed    generators.Float
	StartSize     generators.Float
	Color         generators.Color

	ColorOverLifetime coloroverliftetime.Colorizer
	SizeOverLifetime  sizeoverliftetime.Sizer

	spawner         float64
	distanceSpawner float64

	Rate             float64
	RateOverDistance float64

	Shape shapes.Shape

	burst                Burst
	RotationOverLifetime rotationoverlifetime.Rotator

	// Todo
	// SimulationSpace

	// Emission
	// Bursts

	// Shapes
	//

	// Renderer
}

type Options struct {
	PositionX float64
	PositionY float64

	StartLifetime generators.Float
	StartSize     generators.Float
	StartSpeed    generators.Float
	Color         generators.Color

	ColorOverLifetime    coloroverliftetime.Colorizer
	SizeOverLifetime     sizeoverliftetime.Sizer
	RotationOverLifetime rotationoverlifetime.Rotator

	Rate             *float64
	RateOverDistance *float64

	Shape   shapes.Shape
	Gravity gfx.Vec

	Material *ebiten.Image
	Burst    *Burst
}

var (
	defaultShape = shapes.NewCone(math.Pi/8, 50)
	defaultRate  = 100.0
)

func NewParticleSystem(options Options) *ParticleSystem {
	ps := &ParticleSystem{
		Material:      MaterialStar,
		StartLifetime: &generators.FloatConstant{5.0},
		StartSpeed:    &generators.FloatConstant{1.0},
		StartSize:     &generators.FloatConstant{1.0},
		Color:         &generators.ColorConstant{colornames.White},

		RotationOverLifetime: rotationoverlifetime.RotatorConstant{0},

		Rate:  defaultRate,
		Shape: defaultShape,
	}

	ps.pos.X = options.PositionX
	ps.pos.Y = options.PositionY
	ps.initialPos = ps.pos

	ps.Gravity = options.Gravity

	if options.StartLifetime != nil {
		ps.StartLifetime = options.StartLifetime
	}

	if options.StartSpeed != nil {
		ps.StartSpeed = options.StartSpeed
	}

	if options.StartSize != nil {
		ps.StartSize = options.StartSize
	}

	if options.Color != nil {
		ps.Color = options.Color
	}

	// Emission
	if options.Rate != nil {
		ps.Rate = *options.Rate
	}
	if options.RateOverDistance != nil {
		ps.RateOverDistance = *options.RateOverDistance
	}
	if options.Burst != nil {
		ps.burst = *options.Burst
	}

	if options.ColorOverLifetime != nil {
		ps.ColorOverLifetime = options.ColorOverLifetime
	}

	if options.SizeOverLifetime != nil {
		ps.SizeOverLifetime = options.SizeOverLifetime
	}

	if options.RotationOverLifetime != nil {
		ps.RotationOverLifetime = options.RotationOverLifetime
	}

	if options.Shape != nil {
		ps.Shape = options.Shape
	}

	if options.Material != nil {
		ps.Material = options.Material
	}

	return ps
}

func applyColor(opt *ebiten.DrawImageOptions, clr color.Color) {

	r0, g0, b0, a0 := clr.RGBA()
	r := (float64(r0) / 65536)
	g := (float64(g0) / 65536)
	b := (float64(b0) / 65536)
	a := (float64(a0) / 65536)
	opt.ColorM.Scale(r, g, b, a)
}

func (s *ParticleSystem) newParticle() {
	s.newParticleAt(s.pos)
}

func (s *ParticleSystem) newParticleAt(v gfx.Vec) {
	x, y, angle := s.Shape.New()
	speed := s.StartSpeed.New()
	lifetime := s.StartLifetime.New()

	s.particles = append(s.particles, &Particle{
		pos:          gfx.V(x, y).Add(v).Sub(s.initialPos),
		velocity:     gfx.Unit(angle - gfx.Pi/2).Scaled(speed),
		acceleration: gfx.ZV,

		currentLifetime: lifetime,
		startLifetime:   lifetime,
		startSize:       s.StartSize.New(),

		color: s.Color.New(),
	})
}

func (s *ParticleSystem) Move(cx float64, cy float64) {
	newPos := gfx.V(cx, cy)

	// Get direction and directional vector
	dv := s.pos.To(newPos)
	d := dv.Len()
	dvUnit := dv.Unit()

	// update spawner
	s0 := s.distanceSpawner
	s.distanceSpawner += d * s.RateOverDistance

	// Spawn particles along the path
	for i := 0.0; s.distanceSpawner > 1.0; i++ {
		p := dvUnit.Scaled((i - s0 + 1) / s.RateOverDistance)
		s.newParticleAt(s.pos.Add(p))
		s.distanceSpawner--
	}

	s.pos = newPos
}
