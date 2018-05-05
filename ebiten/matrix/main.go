// Borrowing heavily from https://github.com/hajimehoshi/ebiten/tree/master/examples/life

package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480

	// How many of the cells will increment every tick
	// TODO: this should be a variable, inversely proportional to the sum of the matrix
	ChanceOfIncrement = 1000
	// TODO: variable for how many cells to incrememnt every second

	// How much of a cell's value should spill out to neighbours every tick
	// This is in relation to half of the difference between the cell's value and the neighbour's value
	SpillRate = 1

	// At this value, pixel will be static
	// Only applies to normal color mode
	MaxDisplay = 1000

	// At this value, we cycle
	// For weird, this is  RGB := Modulo(value, CycleValue)
	// For spiral, this is Hue := Modulo(value, CycleValue)
	CycleValue = 1000

	// HSV Value Cycle ratios
	// e.g. for 1 0.2, we will hue value 5 times per value cycle, and cycle value once per CycleValue
	SpiralHueCycleRatio   = 0.1
	SpiralValueCycleRatio = 1

	// How much will it increment by
	IncrementAmount       = 10000
	ManualIncrementAmount = IncrementAmount * 100
	ManualDecrementAmount = IncrementAmount * 100

	// How to render colors. See color.go
	ColorMode = colorSpiral

	// Whether to restrict manual increments by stored value
	RestrictIncrementToStoredValue = true
	// TODO: once I've implemented pressure in addition to value
	// this will make more sense
	RestrictDecrementToMinZero = true
	// TODO: ManualIncrementRange ?
	// TODO: ManualDecrementRange ?
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
	totalVal := world.TotalValue()
	storedVal := world.StoredValue
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("FPS: %.2f\nX: %d, Y: %d, Val: %.2f\nTotal: %g\nStored: %g",
			ebiten.CurrentFPS(),
			x, y, val,
			totalVal,
			storedVal,
		),
	)

	return nil
}

func main() {
	// TODO: cobra? viper? some way of getting config vars

	ebiten.SetRunnableInBackground(true)

	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Some Kinda Grid Flow Thing"); err != nil {
		panic(err)
	}
}
