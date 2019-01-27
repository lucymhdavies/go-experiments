package main

import (
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
		boids:      make([]*Boid, 10, MaxBoids),
		targetSize: 10,
	}
	// currently, flock is made of weird creepy nullboids
	// we create them as part of Update
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (f *Flock) Update() error {
	log.Tracef("f.Update()")

	for i, boid := range f.boids {

		if boid != nil {

			_ = boid.Update(f)

			if boid.IsDead() {
				f.boids[i] = nil
			}
		} else {
			// If nil, then either we have just started, and we are spawning
			// new boids from scratch...
			// or a boid has just died, and we need a new one

			f.boids[i] = NewBoid()
		}

	}

	for len(f.boids) < f.targetSize {
		f.boids = append(f.boids, NewBoid())
	}

	for len(f.boids) > f.targetSize {
		// kill the last boid (make it nil)
		f.boids[len(f.boids)-1] = nil
		f.boids = f.boids[:len(f.boids)-1]
	}

	log.Tracef("END f.Update()")
	return nil
}

func (f Flock) Show(screen *ebiten.Image) error {
	log.Tracef("f.Show()")

	for _, boid := range f.boids {
		if boid != nil {
			_ = boid.Show(screen)
		}
	}

	log.Tracef("END f.Show()")
	return nil
}

func (f Flock) Size() int {
	return len(f.boids)
}
