package main

import (
	"errors"
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

//
// Ebiten Game Interface
//

var (
	game               *Game
	regularTermination = errors.New("regular termination")
)

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
	// For High DPI size
	//s := ebiten.DeviceScaleFactor()
	//return int(float64(outsideWidth) * s), int(float64(outsideHeight) * s)

	// For normal size
	return screenWidth, screenHeight
}
