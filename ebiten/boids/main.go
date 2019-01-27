package main

import (
	"errors"
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// stick these all into some nice struct later...
const (
	WorldWidth  = 800
	WorldHeight = 600

	MaxSpeed = 5
)

var regularTermination = errors.New("regular termination")

func update(screen *ebiten.Image) error {

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return regularTermination
	}

	flock.Update()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	flock.Show(screen)

	msg := fmt.Sprintf(`TPS: %0.2f
FPS: %0.2f
Num of boids: %d
Press Q to quit`,
		ebiten.CurrentTPS(),
		ebiten.CurrentFPS(),
		flock.Size(),
	)
	ebitenutil.DebugPrint(screen, msg)

	// force slow TPS, for debugging
	//time.Sleep(1 * time.Second)

	return nil
}

func main() {
	if err := ebiten.Run(update, WorldWidth, WorldHeight, 1, "Boids!"); err != nil && err != regularTermination {
		panic(err)
	}
}
