package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
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

	particles = []*Particle{}
)

func init() {
	rand.Seed(time.Now().UnixNano())

	dot, _ = ebiten.NewImage(1, 1, ebiten.FilterNearest)
	dot.Fill(color.White)

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
		_ = particle.Draw(screen)
	}

	x, y := ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		`FPS: %.2f
X: %d, Y: %d
Particles: %d`,
		ebiten.CurrentFPS(),
		x, y,
		len(particles),
	),
	)

	return nil
}

func input() error {
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

	// TODO: Increase/Decrease number of particles

	return nil
}
