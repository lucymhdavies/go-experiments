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

func NewHexTile(x, y int) *HexTile {
	sizeX, sizeY := hexImage.Size()

	// Tiles are vertical pointy, so no need for modifying sizeX
	xOffset := sizeX
	yOffset := int(float64(sizeY) * 0.775) // rough guestimate

	// Go to midpoint of screen
	xTranslate, yTranslate := float64(screenWidth)/2, float64(screenHeight)/2

	// Offset by tile position
	xTranslate += float64(xOffset) * float64(x)
	// +y is up
	yTranslate += float64(yOffset) * float64(-y)

	// if odd row, offset by 1/2 xOffset
	if y%2 == 1 || y%2 == -1 {
		xTranslate += float64(xOffset) / 2
	}

	screenMidX := xTranslate
	screenMidY := yTranslate

	// Offset by size of tile
	xTranslate -= float64(sizeX) / 2
	yTranslate -= float64(sizeY) / 2

	return &HexTile{
		image: hexImage,

		x: x,
		y: y,

		screenX: xTranslate,
		screenY: yTranslate,

		screenMidX: screenMidX,
		screenMidY: screenMidY,
	}
}

type HexTile struct {
	image *ebiten.Image

	// Grid coordinates
	x int
	y int

	// Screen coordinates
	// Refers to the top left pixel of the tile
	// Used for drawing on screen
	screenX float64
	screenY float64

	// Midpoint Screen coordinates
	// used to detect closest tile to mouse
	screenMidX float64
	screenMidY float64

	// Highlight this tile?
	highlighted bool
	clicked     bool
}

func (t *HexTile) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}

	if t.highlighted {
		op.ColorM.Scale(1, 1, 1, 1)
	} else {
		op.ColorM.Scale(1, 1, 1, 0.5)
	}
	if t.clicked {
		op.ColorM.Scale(1, 0, 0, 1)
	}

	op.GeoM.Translate(t.screenX, t.screenY)

	screen.DrawImage(t.image, op)
}
