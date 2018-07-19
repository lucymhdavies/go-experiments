package main

import (
	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

var inMenu = false
var debugText = ""

func input() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		inMenu = !inMenu
	}

	if inMenu {
		return menuInput()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		switch cfg.ScreenScale {
		case 1:
			cfg.ScreenScale = 2
			cfg.ScreenWidth = 640
			cfg.ScreenHeight = 480
		case 2:
			cfg.ScreenScale = 1
			cfg.ScreenWidth = 1280
			cfg.ScreenHeight = 960
		default:
			panic("not reached")
		}
	}
	ebiten.SetScreenSize(cfg.ScreenWidth, cfg.ScreenHeight)
	ebiten.SetScreenScale(cfg.ScreenScale)

	//
	// TODO: Increase/Decrease number of particles
	//

	//
	// reset on R
	//
	if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		targets = []*Target{cursorTarget}
	}

	//
	// Targets
	//

	// Cursor position
	cX, cY := ebiten.CursorPosition()

	// Move cursorTarget with mouse cursor
	cursorTarget.Circle.Pos = r3.Vector{X: float64(cX), Y: float64(cY)}

	// Toggle cursor target enabled with space
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		cursorTarget.Enabled = !cursorTarget.Enabled
	}

	// Left click: add new targets
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		target := NewTarget(float64(cX), float64(cY))
		targets = append(targets, target)
	}

	// TODO: right-click (within X pixels) = delete

	return nil
}
