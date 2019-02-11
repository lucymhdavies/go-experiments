package main

import (
	"image/color"
	"math"

	"github.com/golang/geo/r2"
	"github.com/hajimehoshi/ebiten"
)

type Obstacles []*Obstacle

var (
	obstacles = make(Obstacles, 0)
	dot       *ebiten.Image
)

func init() {
	dot, _ = ebiten.NewImage(1, 1, ebiten.FilterNearest)
	dot.Fill(color.RGBA{255, 0, 0, 255})
}

// Generic obstacle
// a thing in the world with size, position, movement

type Obstacle struct {
	position r2.Point
	radius   float64
}

func (os Obstacles) GetNearbyObstacles(b *Boid, distance float64) (Obstacles, error) {

	// Assume no obstacles
	nearbyObstacles := Obstacles{}

	// get position of boid looking for obstacles
	position := b.GetPos()

	// TODO: query the QuadTree which I've not implemented yet
	for _, maybeObstacle := range os {

		// Distance vector from me to potential obstacle
		mOPos := maybeObstacle.position
		distanceVector := mOPos.Sub(position)

		// first, consider taxicab distance, a much quicker calculation
		if math.Abs(distanceVector.X) > distance || math.Abs(distanceVector.Y) > distance {
			continue
		}

		// Knowing that it's in my taxicab neighbourhood, get it's actual distance from me
		distance := distanceVector.Norm()

		// if position within NeighbourhoodDisance...
		if distance <= NeighbourhoodDistance {
			nearbyObstacles = append(nearbyObstacles, maybeObstacle)
		}
	}

	return nearbyObstacles, nil
}

func NewObstacle(x, y float64) *Obstacle {
	o := &Obstacle{
		position: r2.Point{x, y},
		radius:   10.0,
	}

	return o
}

func (o *Obstacle) Show(screen *ebiten.Image) error {

	op.GeoM.Reset()
	op.ColorM.Reset()

	// Move it to its position
	op.GeoM.Scale(o.radius*2+1, o.radius*2+1)
	op.GeoM.Translate(o.position.X-o.radius, o.position.Y-o.radius)

	// Draw the obstacle
	screen.DrawImage(dot, op)

	return nil

}

func (obstacles Obstacles) Show(screen *ebiten.Image) error {

	for _, obstacle := range obstacles {
		obstacle.Show(screen)
	}

	return nil
}
