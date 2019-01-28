package main

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// stick these all into some nice struct later...
const (
	WorldWidth  = 800
	WorldHeight = 600

	MinBoids     = 10
	MaxBoids     = 1000
	InitialBoids = 100
	MaxSpeed     = 5
	MaxForce     = 1

	// How close do other boids need to be to be considered a neighbour
	NeighbourhoodDistance = 50.0
	SeparationDistance    = 10.0

	//
	// Debug Options
	//
	logLevel = log.DebugLevel

	// Whether or not to run at 1 TPS, for debugging
	OneTPS = false

	// Highlight or not
	HighlightPrimary = true
)

var regularTermination = errors.New("regular termination")

func update(screen *ebiten.Image) error {
	log.Tracef("update")

	//
	// Handle user input
	//

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return regularTermination
	}

	// Decrease the nubmer of the sprites.
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		flock.targetSize -= 1
		if flock.targetSize < MinBoids {
			flock.targetSize = MinBoids
		}
	}

	// Increase the nubmer of the sprites.
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		flock.targetSize += 1
		if MaxBoids < flock.targetSize {
			flock.targetSize = MaxBoids
		}
	}

	//
	// Update the flock
	//

	flock.Update()

	//
	// Draw (unless FPS is low)
	//

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	flock.Show(screen)

	msg := fmt.Sprintf(`TPS: %0.2f
FPS: %0.2f
Num of boids: %d
Press <- or -> to change the number of sprites
Press Q to quit`,
		ebiten.CurrentTPS(),
		ebiten.CurrentFPS(),
		flock.Size(),
	)
	ebitenutil.DebugPrint(screen, msg)

	if OneTPS {
		// force slow TPS, for debugging
		time.Sleep(1 * time.Second)
	}

	log.Tracef("END update")

	return nil
}

func main() {
	log.SetLevel(logLevel)
	ebiten.SetRunnableInBackground(true)
	if err := ebiten.Run(update, WorldWidth, WorldHeight, 1, "Boids!"); err != nil && err != regularTermination {
		panic(err)
	}
}
