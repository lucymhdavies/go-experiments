package main

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/lucymhdavies/go-experiments/ebiten/resources/images"
)

//
// Hex Tiles
//

var (
	hexImage *ebiten.Image
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(images.Hex_png))
	if err != nil {
		log.Fatal(err)
	}
	hexImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	// TODO: highlighted tile (mouseover)
	// TODO: highlighted tile (neighbour)
}

type HexTile struct {
	image *ebiten.Image
	x     int
	y     int

	xOffset int
	yOffset int
}

func (t *HexTile) Draw(screen *ebiten.Image) {

	// TODO: calculate position relative to grid on first frame
	// then cache for later re-use

	op := &ebiten.DrawImageOptions{}

	// Go to midpoint of screen
	xTranslate, yTranslate := float64(screenWidth)/2, float64(screenHeight)/2

	// Offset by size of tile
	sizeX, sizeY := hexImage.Size()
	xTranslate -= float64(sizeX) / 2
	yTranslate -= float64(sizeY) / 2

	// Offset by tile position
	xTranslate += float64(t.xOffset) * float64(t.x)
	// +y is up
	yTranslate += float64(t.yOffset) * float64(-t.y)

	// if odd row, offset by 1/2 xOffset
	if t.y%2 == 1 || t.y%2 == -1 {
		xTranslate += float64(t.xOffset) / 2
	}

	op.GeoM.Translate(xTranslate, yTranslate)

	screen.DrawImage(t.image, op)
}
