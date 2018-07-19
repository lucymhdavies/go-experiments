package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Config struct {
	ScreenWidth  int
	ScreenHeight int
	ScreenScale  float64
}

// TODO: new config struct, something like this:

// type Config struct {
// 	Screen struct {
// 		Height int     `json:"Height"`
// 		Width  int     `json:"Width"`
// 		Scale  float64 `json:"Scale"`
// 	} `json:"Screen"`
// 	Particles struct {
// 		Count int `json:"Count"`
// 	} `json:"Particles"`
// }

var (
	cfg *Config = &Config{
		ScreenWidth:  1280,
		ScreenHeight: 960,
		ScreenScale:  1,
	}

	particles    = []*Particle{}
	targets      = []*Target{}
	cursorTarget *Target
)

func init() {
	rand.Seed(time.Now().UnixNano())

	dot, _ = ebiten.NewImage(1, 1, ebiten.FilterNearest)
	dot.Fill(color.White)

	// Init cursorTarget
	cX, cY := ebiten.CursorPosition()
	cursorTarget = NewTarget(float64(cX), float64(cY))
	targets = append(targets, cursorTarget)

}

func main() {
	// TODO: cobra? viper? some way of getting config vars

	// TODO: modify conf in menu, with spacebar = reset

	ebiten.SetRunnableInBackground(true)

	// Simulate the physics independently of game tick, to allow varying simulation speed
	go physicsTicks()

	if err := ebiten.Run(update, cfg.ScreenWidth, cfg.ScreenHeight, cfg.ScreenScale, "Particles!"); err != nil {
		panic(err)
	}
}

func update(screen *ebiten.Image) error {

	// Handle input
	_ = input()

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.Fill(color.Black)

	for _, particle := range particles {
		// If it's just died, skip it
		if particle == nil {
			continue
		}

		// Still alive, so draw it
		_ = particle.Draw(screen)
	}

	for _, target := range targets {
		// If it's just died, skip it
		if target == nil {
			continue
		}

		// Still alive, so draw it
		_ = target.Draw(screen)
	}

	if !inMenu {
		debugText = fmt.Sprintf(
			`FPS: %.2f
Particles: %d / %d`,
			ebiten.CurrentFPS(),
			len(particles), targetNumParticles,
		)
	}

	ebitenutil.DebugPrint(screen, debugText)

	return nil
}
