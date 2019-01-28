package main

import (
	"bytes"
	"image"
	_ "image/png"
	"math"
	"math/rand"

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

	// TODO later: actual boid logic
	// Separation
	// Alignment
	// Cohesion

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
