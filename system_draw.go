package particles

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

func (s *ParticleSystem) Draw(screen *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	//size := s.StartSize.New()
	w := float64(s.Material.Bounds().Dx())
	h := float64(s.Material.Bounds().Dy())

	opt.GeoM.Translate(s.initialPos.X, s.initialPos.Y)

	for _, p := range s.particles {
		// Colorize
		var clr color.Color
		if s.ColorOverLifetime != nil {
			clr = s.ColorOverLifetime.Color(p.normalizedLifetime())
		} else {
			clr = p.color
		}

		// Scale
		scale := p.startSize
		if s.SizeOverLifetime != nil {
			scale *= s.SizeOverLifetime.Size(p.normalizedLifetime())
		}

		o := &ebiten.DrawImageOptions{}
		o.GeoM.Translate(-w/2, -h/2)
		o.GeoM.Rotate(p.rotation)
		o.GeoM.Scale(scale, scale)
		o.GeoM.Translate(p.pos.X, p.pos.Y)

		o.GeoM.Add(opt.GeoM)

		// Colorize
		applyColor(o, clr)

		screen.DrawImage(s.Material, o)
	}
}