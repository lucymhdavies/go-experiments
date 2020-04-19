package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 1280
	screenHeight = 960

	// Size of grid in tiles
	gridWidth  = 25
	gridHeight = 21
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
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
