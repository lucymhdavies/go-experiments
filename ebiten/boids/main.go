package main

import (
	"errors"
	"fmt"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// stick these all into some nice struct later...
const (
	// TODO: all of these, set sensible defaults, but pull actual values from config file

	// Windowed dimensions
	//WorldWidth  = 800
	//WorldHeight = 600

	// Fullscreen dimensions. TODO: get this from os
	WorldWidth  = 1280
	WorldHeight = 720
	Fullscreen  = true

	ShowDebug = true

	MinBoids     = 10
	MaxBoids     = 5000
	InitialBoids = 1000
	MaxSpeed     = 2
	MaxForce     = 0.03

	// Weighting for each boid behaviour
	AlignmentMultiplier  = 1
	SeparationMultiplier = 1.5
	CohesionMultiplier   = 1

	// How close do other boids need to be to be considered a neighbour
	NeighbourhoodDistance = 50.0
	SeparationDistance    = 25.0

	// Avoid obstacles
	ObstacleDistance         = 50.0
	AvoidObstaclesMultiplier = 5

	// TTL
	BoidsHaveTTL  = false
	MaxInitialTTL = 1000
	MinInitialTTL = 100

	//
	// Debug Options
	//
	logLevel = log.DebugLevel

	// Whether or not to run at 1 TPS, for debugging
	OneTPS = false

	// Highlight or not
	HighlightPrimary = false
)

var (
	// How many boids we can update concurrently
	workerPools = runtime.NumCPU()
)

var regularTermination = errors.New("regular termination")

func update(screen *ebiten.Image) error {
	log.Tracef("update")

	// Handle input
	err := input()
	if err != nil {
		return err
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
	obstacles.Show(screen)

	if ShowDebug {
		msg := fmt.Sprintf(`TPS: %0.2f
FPS: %0.2f
Num of boids: %d
Num of obstacles: %d
Press <- or -> to change the number of boids
Click to add obstacles
Press Q to quit`,
			ebiten.CurrentTPS(),
			ebiten.CurrentFPS(),
			flock.Size(),
			len(obstacles),
		)
		ebitenutil.DebugPrint(screen, msg)
	}

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

	if Fullscreen {
		ebiten.SetFullscreen(true)
		//ebiten.SetCursorVisible(false)
	}

	if err := ebiten.Run(update, WorldWidth, WorldHeight, 1, "Boids!"); err != nil && err != regularTermination {
		panic(err)
	}
}
