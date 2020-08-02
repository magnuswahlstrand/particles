package particles

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/particles/assets"
	"github.com/kyeett/particles/generators"
	"github.com/kyeett/particles/modules/coloroverliftetime"
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
	spawner           float64
	distanceSpawner   float64

	Rate             float64
	RateOverDistance float64

	Shape shapes.Shape

	// Todo
	// SimulationSpace

	// Emission
	// RateOverDistance float64
	// Bursts

	// Renderer
}

type Options struct {
	PositionX float64
	PositionY float64

	StartLifetime generators.Float
	StartSize     generators.Float
	StartSpeed    generators.Float
	Color         generators.Color

	ColorOverLifetime coloroverliftetime.Colorizer

	Rate             *float64
	RateOverDistance *float64
	Shape            shapes.Shape
	Gravity          gfx.Vec

	Material *ebiten.Image
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

		ColorOverLifetime: nil,

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

	if options.Rate != nil {
		ps.Rate = *options.Rate
	}

	if options.RateOverDistance != nil {
		ps.RateOverDistance = *options.RateOverDistance
	}

	if options.ColorOverLifetime != nil {
		ps.ColorOverLifetime = options.ColorOverLifetime
	}

	if options.Shape != nil {
		ps.Shape = options.Shape
	}

	if options.Material != nil {
		ps.Material = options.Material
	}

	return ps
}

func (s *ParticleSystem) Draw(screen *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	size := s.StartSize.New()
	w := float64(s.Material.Bounds().Dx())
	h := float64(s.Material.Bounds().Dy())

	opt.GeoM.Translate(s.initialPos.X, s.initialPos.Y)

	for _, p := range s.particles {

		o := &ebiten.DrawImageOptions{}
		o.GeoM.Translate(-w/2, -h/2)
		o.GeoM.Scale(size, size)
		o.GeoM.Translate(p.pos.X, p.pos.Y)

		o.GeoM.Add(opt.GeoM)

		// Colorize
		var clr color.Color
		if s.ColorOverLifetime != nil {
			clr = s.ColorOverLifetime.Color(1 - p.currentLifetime/p.startLifetime)
		} else {
			clr = p.color
		}
		applyColor(o, clr)

		screen.DrawImage(s.Material, o)
	}
}

func applyColor(opt *ebiten.DrawImageOptions, clr color.Color) {

	r0, g0, b0, a0 := clr.RGBA()
	r := (float64(r0) / 65536)
	g := (float64(g0) / 65536)
	b := (float64(b0) / 65536)
	a := (float64(a0) / 65536)
	opt.ColorM.Scale(r, g, b, a)
}

func (s *ParticleSystem) Update(dt float64) {
	for i := 0; i < len(s.particles); i++ {
		p := s.particles[i]

		// Acceleration
		acc := p.acceleration.Add(s.Gravity.Scaled(dt))

		// Velocity
		p.velocity = p.velocity.Add(acc)

		// Position
		p.pos = p.pos.Add(p.velocity)

		// Life cycle
		p.currentLifetime -= dt
		if p.currentLifetime < 0 {
			// From https://github.com/golang/go/wiki/SliceTricks
			// Can be sped up, if needed
			copy(s.particles[i:], s.particles[i+1:])       // Shift a[i+1:] left one index.
			s.particles[len(s.particles)-1] = nil          // Erase last element (write zero value).
			s.particles = s.particles[:len(s.particles)-1] // Truncate slice.
			i--
		}
	}

	s.spawner += dt * s.Rate
	for s.spawner > 0.0 {
		s.newParticle()
		s.spawner--
	}
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

		color: s.Color.New(),
	})
}

func (s *ParticleSystem) Move(cx float64, cy float64) {
	newPos := gfx.V(cx, cy)

	// Get direction and directional vector
	dv := s.pos.To(newPos)
	d := dv.Len()
	dvUnit := dv.Unit()

	// Update spawner
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
