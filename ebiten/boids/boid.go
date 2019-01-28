package main

import (
	"bytes"
	"image"
	_ "image/png"
	"math"
	"math/rand"

	"github.com/golang/geo/r2"
	log "github.com/sirupsen/logrus"

	"github.com/hajimehoshi/ebiten"
	"github.com/lucymhdavies/go-experiments/ebiten/resources/images"
)

var (
	ebitenImage *ebiten.Image
	op          = &ebiten.DrawImageOptions{}
)

type Boid struct {
	w, h   int
	x, y   float64
	vx, vy float64
	ax, ay float64
	angle  float64
	ttl    int

	// reference to flock that boid is a part of
	flock *Flock

	// For debugging, highlight boid 1, and its neighbours
	IsHighlighted bool
	IsNeighbour   bool
}

func init() {
	img, _, err := image.Decode(bytes.NewReader(images.Bullet_png))
	if err != nil {
		log.Fatal(err)
	}
	origEbitenImage, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	w, h := origEbitenImage.Size()
	ebitenImage, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)

	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 0.5)
	ebitenImage.DrawImage(origEbitenImage, op)
}

func (b *Boid) GetPos() (float64, float64) {
	return b.x, b.y
}

func (b *Boid) Update(f *Flock) error {
	// Decrement its TTL
	// (eventually, the boid will die)
	//b.ttl -= 1 // keep boids imortal for now
	if b.ttl <= 0 {
		return nil
	}

	// Move!
	b.x += b.vx
	b.y += b.vy

	// TODO: in future, maybe consider killing the boids if they get off screen instead of warping?

	// If we have left the world, warp to the other side
	if b.x > WorldWidth {
		b.x = 0
	}
	if b.x < 0 {
		b.x = WorldWidth
	}
	if b.y > WorldHeight {
		b.y = 0
	}
	if b.y < 0 {
		b.y = WorldHeight
	}

	// Calculate angle
	// add pi/2, to account for sprite direction
	b.angle = math.Atan2(b.vy, b.vx) + math.Pi/2

	// Get my neighbours
	neighbours, _ := flock.GetNeighbours(b)

	// For debugging, highlight neighbours of primary boid
	if b.IsHighlighted {
		for _, boid := range neighbours {
			boid.IsNeighbour = true
		}
	}

	// Before we do any boiding, reset accelleration
	b.ax, b.ay = 0.0, 0.0

	axAlignment, ayAlignment := b.Alignment(neighbours)
	axSeparation, aySeparation := b.Separation(neighbours)
	axCohesion, ayCohesion := b.Cohesion(neighbours)
	// TODO: multiply each of these by some individually
	// configurable scale factor

	b.ax += (axSeparation + axAlignment + axCohesion) / 3
	b.ay += (aySeparation + ayAlignment + ayCohesion) / 3

	// TODO: no, seriously, just use r2.Point everywhere pls.
	// you'll thank me later

	// Limit accelleration (force) to our MaxForce
	a := r2.Point{b.ax, b.ay}
	a = ConstrainPoint(a, MaxForce)
	b.ax = a.X
	b.ay = a.Y

	// Apply our accelleration to our velocity
	b.vx += b.ax
	b.vy += b.ay

	// TODO: no, seriously, just use r2.Point everywhere pls.
	// you'll thank me later

	// Constrain our velocity to MaxSpeed
	v := r2.Point{b.vx, b.vy}
	v = ConstrainPoint(v, MaxSpeed)
	b.vx = v.X
	b.vy = v.Y

	return nil
}

func NewBoid(f *Flock) *Boid {

	// size of image is important; we need it to get the center of the image
	w, h := ebitenImage.Size()

	// random position somewhere in the world
	// between 0 and world height/width
	x, y := rand.Float64()*float64(WorldWidth-w), rand.Float64()*float64(WorldHeight-h)

	// random velocity, between -MaxSpeed and MaxSpeed
	vx, vy := MaxSpeed*(rand.Float64()*2-1), MaxSpeed*(rand.Float64()*2-1)

	// Set angle when creating it
	angle := math.Atan2(vy, vx) + math.Pi/2

	return &Boid{
		w:     w,
		h:     h,
		x:     x,
		y:     y,
		vx:    vx,
		vy:    vy,
		angle: angle,
		flock: f,

		// Not actually in use yet
		ttl: 100 + rand.Intn(100),
	}
}

func (b *Boid) Show(screen *ebiten.Image) error {

	op.GeoM.Reset()
	op.ColorM.Reset()

	// Rotate around midpoint:
	// Translate to midpoint, rotate, translate back
	op.GeoM.Translate(-float64(b.w)/2, -float64(b.h)/2)
	op.GeoM.Rotate(b.angle)
	op.GeoM.Translate(float64(b.w)/2, float64(b.h)/2)

	// Move it to its position
	op.GeoM.Translate(float64(b.x), float64(b.y))

	if HighlightPrimary {
		// If we're highlighting it...
		if b.IsHighlighted {
			op.ColorM.Scale(1, 0, 0, 1)
		}
		if b.IsNeighbour {
			op.ColorM.Scale(0, 1, 0, 1)
		}
	}

	// Draw the damn boid!
	screen.DrawImage(ebitenImage, op)

	return nil

}

func (b *Boid) IsDead() bool {
	return b.ttl <= 0
}

// Alignment:
// steer towards the average heading of local flockmates
func (b *Boid) Alignment(neighbours []*Boid) (float64, float64) {

	//forceX, forceY := 0.0, 0.0
	// 	averageAngle := 0.0
	//
	// 	for _, neighbour := range neighbours {
	// 		averageAngle += neighbour.angle
	// 	}
	// 	averageAngle = averageAngle / float64(len(neighbours))
	//
	// 	magnitude = 1
	//
	// 	forceX = magnitude * math.Cos(averageAngle)
	// 	forceY = magnitude * math.Sin(averageAngle)

	force := r2.Point{0.0, 0.0}

	// Add up all their velocities, normalize, then multiply by MaxForce
	for _, neighbour := range neighbours {
		force.X += neighbour.vx
		force.Y += neighbour.vy
	}
	force = ConstrainPoint(force, MaxForce)

	return force.X, force.Y
}

// Separation:
// steer to avoid crowding local flockmates
func (b *Boid) Separation(neighbours []*Boid) (float64, float64) {

	return 0.0, 0.0
}

// Cohesion:
// steer to move toward the average position of local flockmates
func (b *Boid) Cohesion(neighbours []*Boid) (float64, float64) {

	return 0.0, 0.0
}

func ConstrainPoint(p r2.Point, max float64) r2.Point {

	// Get our current magnitude
	magnitude := math.Sqrt(math.Pow(p.X, 2) + math.Pow(p.Y, 2))

	// If we're too fast, slow down
	if magnitude > max {
		// Unit point in same direction * max
		p = p.Normalize()
		p = p.Mul(max)
	}

	return p
}
