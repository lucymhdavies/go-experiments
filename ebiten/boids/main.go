package main

import (
	"errors"
	"fmt"
	"runtime"
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
	MaxBoids     = 5000
	InitialBoids = 1000
	MaxSpeed     = 5
	MaxForce     = 1

	// Weighting for each boid behaviour
	AlignmentMultiplier  = 1.5
	SeparationMultiplier = 1
	CohesionMultiplier   = 1

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

var (
	// How many boids we can update concurrently
	workerPools = runtime.NumCPU()
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
