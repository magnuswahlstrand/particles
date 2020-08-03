package particles

func (s *ParticleSystem) Update(dt float64) {
	for i := 0; i < len(s.particles); i++ {
		p := s.particles[i]

		// Acceleration
		acc := p.acceleration.Add(s.Gravity.Scaled(dt))

		// Velocity
		p.velocity = p.velocity.Add(acc)

		// Position
		p.pos = p.pos.Add(p.velocity)

		// Rotation
		p.rotation += s.RotationOverLifetime.Rotation(p.normalizedLifetime())

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

	s.burst.update(dt)
	n := s.burst.due()
	for i := n; i > 0; i-- {
		s.newParticle()
	}
}
