package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 640
	screenHeight = 480
	screenScale  = 2

	// Size of grid in tiles
	gridWidth  = 15
	gridHeight = 13
)

func main() {
	ebiten.SetWindowSize(screenWidth*screenScale, screenHeight*screenScale)
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
