package main

import (
	"math"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/hajimehoshi/ebiten"
)

type Flock struct {
	boids      []*Boid
	targetSize int
}

var (
	flock = &Flock{
		boids:      make([]*Boid, InitialBoids, MaxBoids),
		targetSize: InitialBoids,
	}
	// currently, flock is made of weird creepy nullboids
	// we create them as part of Update
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (f *Flock) Update() error {
	log.Tracef("f.Update()")

	// Very first thing: mark every boid as not a neightbour of our
	// primary boid.
	for i, boid := range f.boids {
		if boid != nil {
			f.boids[i].IsNeighbour = false

			if i == 0 {
				// Highlight primary boid
				f.boids[i].IsHighlighted = true
			}
		}
	}

	for i, boid := range f.boids {

		if boid != nil {

			_ = boid.Update(f)

			if boid.IsDead() {
				// TODO: killing boids should instead be a case of removing them
				// from the slice, in addition to setting them to nil
				f.boids[i] = nil
			}

		} else {
			// If nil, then either we have just started, and we are spawning
			// new boids from scratch...
			// or a boid has just died, and we need a new one

			// TODO: this whole bit may not be necessary, if we're creating boids below
			f.boids[i] = NewBoid(f)
		}

	}

	// Create new boids if we don't have enough
	for len(f.boids) < f.targetSize {
		f.boids = append(f.boids, NewBoid(f))
	}

	// Delete boids if we have too much
	for len(f.boids) > f.targetSize {
		// kill the last boid (make it nil)
		f.boids[len(f.boids)-1] = nil
		f.boids = f.boids[:len(f.boids)-1]
	}

	// TODO: in future, index all boids by position
	// e.g. with a QuadTree

	log.Tracef("END f.Update()")
	return nil
}

func (f *Flock) Show(screen *ebiten.Image) error {
	log.Tracef("f.Show()")

	for _, boid := range f.boids {
		if boid != nil {
			_ = boid.Show(screen)
		}
	}

	log.Tracef("END f.Show()")
	return nil
}

func (f *Flock) Size() int {
	return len(f.boids)
}

// GetNeighbours returns boids within a specific distance of a specific boid
func (f *Flock) GetNeighbours(b *Boid) ([]*Boid, error) {

	// Assume I'm a lonely hermit with no friends :(
	neighbours := []*Boid{}

	// Get my position
	position := b.GetPos()

	// TODO: query the QuadTree which I've not implemented yet
	for _, maybeNeighbour := range f.boids {
		if b == maybeNeighbour {
			continue
		}

		// Get it's distance from me
		mNPos := maybeNeighbour.GetPos()
		dX := mNPos.X - position.X
		dY := mNPos.Y - position.Y

		// a^2 + b^2 = c^2
		// simple pythagoras
		distance := math.Sqrt(math.Pow(dY, 2) + math.Pow(dX, 2))

		// if position within NeighbourhoodDisance...
		if distance <= NeighbourhoodDistance {
			neighbours = append(neighbours, maybeNeighbour)
		}
	}

	// Show me my neighbours!
	return neighbours, nil
}
