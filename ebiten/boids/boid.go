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
	w, h         int
	position     r2.Point
	velocity     r2.Point
	acceleration r2.Point
	angle        float64

	// Unused... for now
	ttl int

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

func (b *Boid) GetPos() r2.Point {
	return b.position
}

func (b *Boid) Update(f *Flock) error {
	// Decrement its TTL
	// (eventually, the boid will die)
	//b.ttl -= 1 // keep boids imortal for now
	if b.ttl <= 0 {
		return nil
	}

	// Move!
	b.position = b.position.Add(b.velocity)

	// TODO: in future, maybe consider killing the boids if they get off screen instead of warping?

	// If we have left the world, warp to the other side
	if b.position.X > WorldWidth {
		b.position.X = 0
	}
	if b.position.X < 0 {
		b.position.X = WorldWidth
	}
	if b.position.Y > WorldHeight {
		b.position.Y = 0
	}
	if b.position.Y < 0 {
		b.position.Y = WorldHeight
	}

	// Calculate angle
	// add pi/2, to account for sprite direction
	b.angle = math.Atan2(b.velocity.Y, b.velocity.X) + math.Pi/2

	// Get my neighbours
	neighbours, _ := flock.GetNeighbours(b)

	// For debugging, highlight neighbours of primary boid
	if b.IsHighlighted {
		for _, boid := range neighbours {
			boid.IsNeighbour = true
		}
	}

	// Before we do any boiding, reset acceleration
	b.acceleration = r2.Point{0.0, 0.0}

	alignment := b.Alignment(neighbours)
	alignment = alignment.Mul(AlignmentMultiplier)

	separation := b.Separation(neighbours)
	separation = separation.Mul(SeparationMultiplier)

	cohesion := b.Cohesion(neighbours)
	cohesion = cohesion.Mul(CohesionMultiplier)

	b.acceleration = alignment.Add(separation.Add(cohesion))
	b.acceleration = b.acceleration.Mul(1.0 / 3.0)

	// Limit acceleration (force) to our MaxForce
	b.acceleration = ConstrainPoint(b.acceleration, MaxForce)

	// Apply our acceleration to our velocity
	b.velocity = b.velocity.Add(b.acceleration)

	// Constrain our velocity to MaxSpeed
	b.velocity = ConstrainPoint(b.velocity, MaxSpeed)

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
		w:        w,
		h:        h,
		position: r2.Point{x, y},
		velocity: r2.Point{vx, vy},
		angle:    angle,
		flock:    f,

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
	op.GeoM.Translate(b.position.X, b.position.Y)

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
func (b *Boid) Alignment(neighbours []*Boid) r2.Point {

	force := r2.Point{0.0, 0.0}

	// Add up all their velocities, normalize, then multiply by MaxForce
	for _, neighbour := range neighbours {
		force = force.Add(neighbour.velocity)
	}
	force = force.Sub(b.velocity)
	force = ConstrainPoint(force, MaxForce)

	return force
}

// Separation:
// steer to avoid crowding local flockmates
func (b *Boid) Separation(neighbours []*Boid) r2.Point {

	force := r2.Point{0.0, 0.0}

	for _, neighbour := range neighbours {
		// vector from neighbour to us
		distanceVector := b.position.Sub(neighbour.position)
		distance := distanceVector.Norm()

		// how many neighbours are too close
		count := 0

		// if we are too close...
		if distance < SeparationDistance {
			// unit vector pointing away from neighbour
			distanceVector = distanceVector.Normalize()

			// Weight inversely to distance
			distanceVector = distanceVector.Mul(1.0 / distance)

			force.Add(distanceVector)

			count++
		}

		if count > 0 {
			force.Mul(1.0 / float64(count))
		}
	}

	force = force.Sub(b.velocity)
	force = ConstrainPoint(force, MaxForce)

	return force
}

// Cohesion:
// steer to move toward the average position of local flockmates
func (b *Boid) Cohesion(neighbours []*Boid) r2.Point {

	force := r2.Point{0.0, 0.0}

	return force
}

func ConstrainPoint(p r2.Point, max float64) r2.Point {

	// Get our current magnitude
	magnitude := p.Norm()

	// If we're too fast, slow down
	if magnitude > max {
		// Unit point in same direction * max
		p = p.Normalize()
		p = p.Mul(max)
	}

	return p
}
