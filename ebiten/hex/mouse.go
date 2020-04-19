package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

func mouseUpdate(g *Game) {

	cursorX, cursorY := ebiten.CursorPosition()

	g.debugMessage = fmt.Sprintf("%d, %d", cursorX, cursorY)

	nearestTile, mouseOverTile := g.grid.FindNearestTile(float64(cursorX), float64(cursorY))

	g.debugMessage += fmt.Sprintf(" - %d, %d", nearestTile.x, nearestTile.y)

	if mouseOverTile {
		nearestTile.highlighted = true

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			nearestTile.clicked = !nearestTile.clicked
		}
	}
}
