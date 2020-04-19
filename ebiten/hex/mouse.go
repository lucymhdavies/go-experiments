package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
)

func mouseUpdate(g *Game) {

	cursorX, cursorY := ebiten.CursorPosition()

	g.debugMessage = fmt.Sprintf("%d, %d", cursorX, cursorY)

	nearestTile := g.grid.FindNearestTile(float64(cursorX), float64(cursorY))

	g.debugMessage += fmt.Sprintf(" - %d, %d", nearestTile.x, nearestTile.y)
}
