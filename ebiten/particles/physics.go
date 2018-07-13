package main

import (
	"math/rand"
	"time"

	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"
)

var ()

func physicsTicks() error {
	for i := 0; i < 25000; i++ {
		particles = append(particles,
			NewParticle(
				rand.Float64()*float64(cfg.ScreenWidth),
				rand.Float64()*float64(cfg.ScreenHeight),
				0, // radius, unused
				//0, 0, // velocity
				rand.Float64()-0.5, rand.Float64()-0.5,
				0, 0, // accel
			),
		)
	}

	// TODO: array of targets (mouse click = place target)
	// Target = green circle

	for {
		// TODO: if number of particles > cfg.Particles.Count
		// delete some (like 1% of difference, per tick?)
		// Same if < cfg.Particles.Count

		// Attract towards mouse cursor
		cX, cY := ebiten.CursorPosition()
		target := r3.Vector{X: float64(cX), Y: float64(cY)}

		for _, particle := range particles {
			_ = particle.Attract(target)

			// Attract to all other particles
			// TODO: if cfg.Particles.TargetOtherParticles == true
			/*
				for _, pTarget := range particles {
					_ = particle.Attract(pTarget.Pos)
				}
			*/

			_ = particle.Update()
		}
		time.Sleep(time.Second / 120.0)
	}

	return nil
}
