package main

import (
	"math"

	"github.com/golang/geo/r3"
)

type Particle struct {
	Circle
	Velocity     r3.Vector
	Acceleration r3.Vector
	TTL          int
}

// TODO: stick these in Config

var (
	// Infinite = -1
	initialParticleTTL = 2000
	deleteOnBounce     = false //true
)

func NewParticle(x, y, r, dx, dy, ddx, ddy float64, ttl int) *Particle {
	return &Particle{
		Circle:       NewCircle(x, y, r),
		Velocity:     r3.Vector{X: dx, Y: dy},
		Acceleration: r3.Vector{X: ddx, Y: ddy},
		TTL:          ttl,
	}
}

func (p *Particle) Update() error {
	// TODO: take into account FPS?
	// Or maybe move particles in a separate goroutine

	// TODO: friction? slow down when we get close?

	//
	// Velocity
	//

	p.Velocity = p.Velocity.Add(p.Acceleration)
	// Constrain speed
	direction := p.Velocity.Normalize()
	speed := constrain(p.Velocity.Distance(r3.Vector{0, 0, 0}), 0, 10)
	p.Velocity = direction.Mul(speed)

	// Reset accelleration
	p.Acceleration = p.Acceleration.Mul(0)

	//
	// Position
	//

	p.Pos = p.Pos.Add(p.Velocity)

	//
	// TTL
	//

	if p.TTL > 0 {
		p.TTL--
	}

	return nil
}

func (p *Particle) Attract(cPos r3.Vector) error {
	//
	// Accelleration
	//

	// force = target (mouse cursor) - my current location
	force := cPos.Sub(p.Pos)
	// Constrain force
	// https://github.com/CodingTrain/website/blob/master/CodingChallenges/CC_056_attraction_repulsion/particle.js#L35-L41
	distance := force.Distance(r3.Vector{0, 0, 0})
	forceDirection := force.Normalize()
	scalarForce := constrain(distance, 1, 25)

	// Force should be stronger the closer you are
	G := 50.0 // Gravitational Constant equivalent
	strength := G / (scalarForce * scalarForce)

	force = forceDirection.Mul(strength)

	// If we're too close, go away!
	// https://github.com/CodingTrain/website/blob/master/CodingChallenges/CC_056_attraction_repulsion/particle.js#L39-L41
	// TODO: disable bounce with config?

	if distance < 20 {

		if deleteOnBounce {
			// It got to close! kill it
			p.TTL = 0
		} else {
			force = force.Mul(-10)
			// Too close: bounce it away
		}
	}

	// Force = Mass * Acceleration
	// Assuming mass of 1 for now, F = A
	//force = force.Mul(0.1)
	p.Acceleration = p.Acceleration.Add(force)

	return nil
}

func constrain(x, min, max float64) float64 {
	return math.Max(min, math.Min(max, x))
}
