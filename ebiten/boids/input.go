package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

func input() error {
	//
	// Handle user input
	//

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return regularTermination
	}

	// Increase/Decrease faster when shift is held
	incrementAmount := 1
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		incrementAmount = 10
	}

	// Decrease the nubmer of the sprites.
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		flock.targetSize -= incrementAmount
		if flock.targetSize < MinBoids {
			flock.targetSize = MinBoids
		}
	}

	// Increase the nubmer of the sprites.
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		flock.targetSize += incrementAmount
		if MaxBoids < flock.targetSize {
			flock.targetSize = MaxBoids
		}
	}

	// Cursor position
	cX, cY := ebiten.CursorPosition()

	// Left click: add new obstacle
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		obstacle := NewObstacle(float64(cX), float64(cY))
		obstacles = append(obstacles, obstacle)
	}

	return nil
}
