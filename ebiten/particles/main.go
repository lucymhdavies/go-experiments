package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Config struct {
	ScreenWidth  int
	ScreenHeight int
}

var (
	cfg *Config = &Config{
		ScreenHeight: 480,
		ScreenWidth:  640,
	}

	// TODO: Array of these
	particles []*Particle
)

func init() {
	rand.Seed(time.Now().UnixNano())

	dot, _ = ebiten.NewImage(1, 1, ebiten.FilterNearest)
	dot.Fill(color.White)

	// TODO: a bunch of random particles, with random velocities
	particles = []*Particle{}

	for i := 0; i < 25000; i++ {
		particles = append(particles,
			NewParticle(
				rand.Float64()*float64(cfg.ScreenWidth),
				rand.Float64()*float64(cfg.ScreenHeight),
				0, 0, 0, 0, 0),
		)
	}

}

func main() {
	// TODO: cobra? viper? some way of getting config vars

	// TODO: modify conf in menu, with spacebar = reset

	ebiten.SetRunnableInBackground(true)

	if err := ebiten.Run(update, cfg.ScreenWidth, cfg.ScreenHeight, 2, "Particles!"); err != nil {
		panic(err)
	}
}

func update(screen *ebiten.Image) error {

	if ebiten.IsRunningSlowly() {
		return nil
	}

	screen.Fill(color.Black)

	// Attract towards mouse cursor
	cX, cY := ebiten.CursorPosition()
	target := r3.Vector{X: float64(cX), Y: float64(cY)}

	for _, particle := range particles {
		_ = particle.Attract(target)

		/*
			// Attract to all other particles
			for _, pTarget := range particles {
				_ = particle.Attract(pTarget.Pos)
			}
		*/

		_ = particle.Update()
		_ = particle.Draw(screen)
	}

	x, y := ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f\nX: %d, Y: %d", ebiten.CurrentFPS(), x, y))

	return nil
}
