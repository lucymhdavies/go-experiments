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

	debugMessage string
}

func (g *Game) Update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return regularTermination
	}

	mouseUpdate(g)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.grid.Draw(screen)

	msg := fmt.Sprintf(`TPS: %0.2f
FPS: %0.2f
Press Q to quit
%s`,
		ebiten.CurrentTPS(), ebiten.CurrentFPS(),
		g.debugMessage)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// For High DPI size
	//s := ebiten.DeviceScaleFactor()
	//return int(float64(outsideWidth) * s), int(float64(outsideHeight) * s)

	// For normal size
	return screenWidth, screenHeight
}
