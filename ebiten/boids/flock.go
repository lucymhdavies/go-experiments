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
)

// TODO: move this to an AddBoid function later
// the idea being, we start with an empty flock, and add one boid per tick
func init() {
	// currently, flock is made of weird creepy nullboids
	// make them exist pls

	rand.Seed(time.Now().UnixNano())

	for i := range flock.boids {
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
			ttl:         100,
		}
	}
}

func (f Flock) Update() error {

	for _, boid := range f.boids {

		_ = boid.Update(&f)
		// TODO: error handling

		// TODO: if boid.ttl <= 0
		// kill it
	}

	// TODO: if flock is too small, add another boid

	return nil
}

func (f Flock) Show(screen *ebiten.Image) error {

	for _, boid := range f.boids {
		_ = boid.Show(screen)
		// TODO: error handling
	}

	return nil
}

func (f Flock) Size() int {
	return len(f.boids)
}
