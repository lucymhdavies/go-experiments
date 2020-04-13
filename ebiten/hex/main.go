package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/lucymhdavies/go-experiments/ebiten/resources/images"
)

const (
	screenWidth  = 640
	screenHeight = 480

	// Size of grid in tiles
	gridWidth  = 7
	gridHeight = 7
)

var (
	hexImage           *ebiten.Image
	regularTermination = errors.New("regular termination")

	game *Game
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

//
// Hex Grid
//

type HexGrid struct {
	width  int
	height int

	// TODO: Store this in a structure that allows searching based on coords
	tiles []*HexTile
}

func NewHexGrid(width, height int) *HexGrid {
	sizeX, sizeY := hexImage.Size()
	// Tiles are vertical pointy, so no need for modifying sizeX
	sizeY = int(float64(sizeY) * 0.775) // rough guestimate

	tiles := []*HexTile{}

	// TODO: generate all tiles based on width/height
	// 0,0 is center

	for x := -3; x <= 3; x++ {
		for y := -3; y <= 3; y++ {
			tile := &HexTile{
				image: hexImage,
				x:     x,
				y:     y,

				xOffset: sizeX,
				yOffset: sizeY,
			}

			tiles = append(tiles, tile)
		}
	}

	return &HexGrid{
		width:  width,
		height: height,

		tiles: tiles,
	}
}

func (g *HexGrid) Draw(screen *ebiten.Image) {

	for _, tile := range g.tiles {
		tile.Draw(screen)
	}
}

//
// Hex Tiles
//

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

//
// Ebiten Game Interface
//

type Game struct {
	grid *HexGrid
}

func (g *Game) Update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return regularTermination
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.grid.Draw(screen)

	msg := fmt.Sprintf(`TPS: %0.2f
FPS: %0.2f
Press Q to quit`, ebiten.CurrentTPS(), ebiten.CurrentFPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	//s := ebiten.DeviceScaleFactor()
	//return int(float64(outsideWidth) * s), int(float64(outsideHeight) * s)
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Hex Tiles")

	game = &Game{
		grid: NewHexGrid(gridWidth, gridHeight),
	}

	if err := ebiten.RunGame(game); err != nil {
		if err != regularTermination {
			log.Fatal(err)
		}
	}
}
