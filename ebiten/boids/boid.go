package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/lucymhdavies/go-experiments/ebiten/resources/images"
)

var (
	ebitenImage *ebiten.Image
	op          = &ebiten.DrawImageOptions{}
)

type Boid struct {
	imageWidth, imageHeight int
	x, y                    float64
	vx, vy                  float64
	angle                   float64
	ttl                     int
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

func (b Boid) Update() error {
	b.ttl -= 1

	// TODO: move, update angle, etc. all that fun stuff.

	return nil
}

func (b Boid) Show(screen *ebiten.Image) error {

	op.GeoM.Reset()

	// Rotate around midpoint:
	// Translate to midpoint, rotate, translate back
	op.GeoM.Translate(-float64(b.imageWidth)/2, -float64(b.imageHeight)/2)
	op.GeoM.Rotate(b.angle)
	op.GeoM.Translate(float64(b.imageWidth)/2, float64(b.imageHeight)/2)

	// Move it to its position
	op.GeoM.Translate(float64(b.x), float64(b.y))

	// Draw the damn boid!
	screen.DrawImage(ebitenImage, op)

	return nil

}
