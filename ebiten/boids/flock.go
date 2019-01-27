package main

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
)

type Flock struct {
	boids      []*Boid
	targetSize int
}

var (
	flock = &Flock{
		// TODO; init this to 0, and add the boids in Update
		boids:      make([]*Boid, 10, 10),
		targetSize: 10,
	}
	// currently, flock is made of weird creepy nullboids
	// we create them as part of Update
)

// TODO: move this to an AddBoid function later
// the idea being, we start with an empty flock, and add one boid per tick
func init() {
	// make them exist pls

	rand.Seed(time.Now().UnixNano())
}

func (f Flock) Update() error {

	for i, boid := range f.boids {

		if boid != nil {

			_ = boid.Update(&f)

			if boid.IsDead() {
				flock.boids[i] = nil
			}
		} else {
			// If nil, then either we have just started, and we are spawning
			// new boids from scratch...
			// or a boid has just died, and we need a new one

			// TODO: call a NewBoid function instead of doing all the logic here

			// size of image is important; we need it to get the center of the image
			w, h := ebitenImage.Size()

			// random position somewhere in the world
			// between 0 and world height/width
			x, y := rand.Float64()*float64(WorldWidth-w), rand.Float64()*float64(WorldHeight-h)

			// random velocity, between -MaxSpeed and MaxSpeed
			vx, vy := MaxSpeed*(rand.Float64()*2-1), MaxSpeed*(rand.Float64()*2-1)

			flock.boids[i] = &Boid{
				imageWidth:  w,
				imageHeight: h,
				x:           x,
				y:           y,
				vx:          vx,
				vy:          vy,
				ttl:         100 + rand.Intn(100),
			}

		}

	}

	return nil
}

func (f Flock) Show(screen *ebiten.Image) error {

	for _, boid := range f.boids {
		if boid != nil {
			_ = boid.Show(screen)
		}
	}

	return nil
}

func (f Flock) Size() int {
	return len(f.boids)
}
