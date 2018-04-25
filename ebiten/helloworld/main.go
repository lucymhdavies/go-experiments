package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func update(screen *ebiten.Image) error {

	// Get the x, y position of the cursor from the CursorPosition() function
	x, y := ebiten.CursorPosition()

	// Display the information with "X: xx, Y: xx" format
	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %d, Y: %d", x, y))

	return nil
}

func main() {
	if err := ebiten.Run(update, 320, 240, 2, "Hello world!"); err != nil {
		panic(err)
	}
}
