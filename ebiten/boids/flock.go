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

func worker(id int, jobs <-chan *Boid, results chan<- bool, f *Flock) {
	log.Tracef("Worker %v started", id)

	for boid := range jobs {

		log.Tracef("Worker %v processing boid %v", id, boid)

		_ = boid.Update(f)

		results <- true
	}
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

	// Set up some channels to spawn goroutines
	jobs := make(chan *Boid, f.Size())
	results := make(chan bool, f.Size())
	// Launch the worker goroutines
	for w := 1; w <= workerPools; w++ {
		go worker(w, jobs, results, f)
	}

	for i, boid := range f.boids {
		log.Tracef("Updating boid %v", i)

		if boid != nil {
			if boid.IsDead() {
				// TODO: killing boids should instead be a case of removing them
				// from the slice, in addition to setting them to nil
				f.boids[i] = nil
			} else {
				log.Tracef("Putting boid %v onto jobs channel", i)
				jobs <- boid
			}
		} else {
			// If nil, then either we have just started, and we are spawning
			// new boids from scratch...
			// or a boid has just died, and we need a new one

			// TODO: this whole bit may not be necessary, if we're creating boids below
			f.boids[i] = NewBoid(f)
			results <- true
		}

	}
	close(jobs)
	for i, _ := range f.boids {
		log.Tracef("Result from boid %v", i)
		<-results
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

		// Get distance from me to potential neighbour
		mNPos := maybeNeighbour.GetPos()
		distanceVector := mNPos.Sub(position)

		// first, consider taxicab distance, a much quicker calculation
		if math.Abs(distanceVector.X) > NeighbourhoodDistance || math.Abs(distanceVector.Y) > NeighbourhoodDistance {
			continue
		}

		// Knowing that it's in my taxicab neighbourhood, get it's actual distance from me
		distance := distanceVector.Norm()

		// if position within NeighbourhoodDisance...
		if distance <= NeighbourhoodDistance {
			neighbours = append(neighbours, maybeNeighbour)
		}
	}

	// Show me my neighbours!
	return neighbours, nil
}
