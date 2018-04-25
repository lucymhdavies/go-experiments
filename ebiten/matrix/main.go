// Borrowing heavily from https://github.com/hajimehoshi/ebiten/tree/master/examples/life

package main

import (
	"fmt"
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480

	// How many of the cells will increment every tick
	// TODO: this should be a variable, inversely proportional to the sum of the matrix
	ChanceOfIncrement = 1
	// TODO: variable for how many cells to incrememnt every second

	// How much of a cell's value should spill out to neighbours every tick
	// This is in relation to half of the difference between the cell's value and the neighbour's value
	SpillRate = 1

	// At this value, pixel will be entirely white
	MaxDisplay = 1000

	// How much will it increment by
	IncrementAmount = 100000
)

var (
	world  = NewWorld(screenWidth, screenHeight, int((screenWidth*screenHeight)/10))
	pixels = make([]byte, screenWidth*screenHeight*4)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func update(screen *ebiten.Image) error {

	if err := world.Update(); err != nil {
		return err
	}

	if ebiten.IsRunningSlowly() {
		return nil
	}

	world.Draw(pixels)
	screen.ReplacePixels(pixels)

	// TODO: keyboard to toggle debugprint
	x, y := ebiten.CursorPosition()
	val, _ := world.GetVal(x, y)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f\nX: %d, Y: %d, Val: %.2f", ebiten.CurrentFPS(), x, y, val))

	return nil
}

func main() {
	if false {
		log.SetLevel(log.DebugLevel)
	}

	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Some Kinda Grid Flow Thing"); err != nil {
		panic(err)
	}
}
