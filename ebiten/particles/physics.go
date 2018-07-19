package main

import (
	"math"
	"math/rand"
)

// TODO: stick these in config struct
var (
	targetNumParticles = 25000
	particlesPerTick   = 100
)

func physicsTicks() error {

	//for {
	// TODO: if number of particles > cfg.Particles.Count
	// delete some (like 1% of difference, per tick?)

	// If we don't have enough particles, spawn some in
	if len(particles) < targetNumParticles {

		// TODO: if the difference between len(particles) and targetNumParticles is < particlesPerTick
		// use diff

		// If we're really close to target number of particles, don't spawn more than necessary
		particlesToSpawnThisTick := int(
			math.Min(
				float64(targetNumParticles-len(particles)),
				float64(particlesPerTick)))

		for i := 0; i < particlesToSpawnThisTick; i++ {
			particles = append(particles,
				NewParticle(
					rand.Float64()*float64(cfg.ScreenWidth),
					rand.Float64()*float64(cfg.ScreenHeight),
					0, // radius, unused
					//0, 0, // velocity
					10*(rand.Float64()-0.5), 10*(rand.Float64()-0.5),
					0, 0, // accel
					int(rand.Float64()*float64(initialParticleTTL)), // ttl
				),
			)
		}
	}

	for i, particle := range particles {
		// If we've recently deleted the particle...
		if particle == nil {
			continue
		}

		if particle.TTL == 0 {
			// Delete without preserving order
			// https://github.com/golang/go/wiki/SliceTricks#delete
			particles[i] = particles[len(particles)-1]
			particles[len(particles)-1] = nil
			particles = particles[:len(particles)-1]

			continue
		}

		for _, target := range targets {
			if target.Enabled {
				target.Visible = true
				_ = particle.Attract(target.Circle.Pos)
			} else {
				target.Visible = false
			}
		}

		_ = particle.Update()
	}
	//	time.Sleep(time.Second / 120.0)
	//}

	return nil
}
