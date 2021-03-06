package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

// TODO:
// we currently have x,y coords, and a single value
// how about x,y coords, and r,g,b values?

type Matrix = [][]float64

type World struct {
	matrix        Matrix // [x][y]
	width, height int    // convenience vars
	StoredValue   float64
}

// newMatrix initialises a new empty matrix
func newMatrix(width, height int) Matrix {
	a := make(Matrix, width)
	for i := 0; i < width; i++ {
		a[i] = make([]float64, height)
	}
	return a
}

// NewWorld creates a new world.
func NewWorld(width, height int, maxInitLiveCells int) *World {
	w := &World{
		matrix: newMatrix(width, height),
		width:  width,
		height: height,
	}
	return w
}

func (w *World) Reset() error {
	w.matrix = newMatrix(w.width, w.height)
	w.StoredValue = 0
	randomIncrements = false

	return nil
}

var (
	randomIncrements = false // default to off
)

// Update game state by one tick.
func (w *World) Update() error {

	// Spill to neighbours
	w.spillToNeighbours()

	// Randomly Increment
	if randomIncrements {
		w.randomIncrementMatrix()
	}

	// Respond to clicks
	w.respondToInput()

	return nil
}

// Draw paints current game state.
func (w *World) Draw(pix []byte) {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			idx := 4*y*w.width + 4*x

			color := Float64ToColor(w.matrix[x][y])

			pix[idx] = color.R
			pix[idx+1] = color.G
			pix[idx+2] = color.B
			pix[idx+3] = color.A
		}
	}
}

// Get val at specified coord, converted to int
func (w *World) GetVal(x, y int) (float64, error) {
	m := w.matrix

	if x >= 0 && x < w.width && y >= 0 && y < w.height {
		return m[x][y], nil
	}

	return -1, fmt.Errorf("Out of Bounds")
}

// randomIncrementMatrix will, for every cell, with a certain percentage, either increment or not
func (w *World) randomIncrementMatrix() {

	// TODO: do this in a for loop, as per original version

	// TODO: option for cell's chance of incrementing to be inversely proportional to its value

	cellsToIncrement := 0
	if rand.Float64() < ChanceOfIncrement {
		cellsToIncrement = 10
	}

	for i := 0; i < cellsToIncrement; i++ {
		randY := rand.Intn(w.height)
		randX := rand.Intn(w.width)

		w.matrix[randX][randY] = w.matrix[randX][randY] + IncrementAmount
	}

}

// Keep track of which keys have been pressed, but not released
var (
	keySpaceDown = false
	keyRDown     = false
)

// respondToInput responds to user input: clicks, key presses, etc
func (w *World) respondToInput() {

	// TODO: do this in a for loop, as per original version
	x, y := ebiten.CursorPosition()
	if x >= 0 && x < w.width && y >= 0 && y < w.height {

		// Left click, or Q: incremement
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || ebiten.IsKeyPressed(ebiten.KeyQ) {
			if RestrictIncrementToStoredValue {
				// only incremement if we have stored some value already
				incrementAmount := math.Min(float64(ManualIncrementAmount), w.StoredValue)

				w.matrix[x][y] = w.matrix[x][y] + incrementAmount
				w.StoredValue = w.StoredValue - incrementAmount
			} else {
				w.matrix[x][y] = w.matrix[x][y] + ManualIncrementAmount
			}
		}

		// Right click, or W: decremement
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) || ebiten.IsKeyPressed(ebiten.KeyW) {
			if RestrictIncrementToStoredValue {

				decrementAmount := float64(ManualDecrementAmount)

				if RestrictDecrementToMinZero {
					decrementAmount = math.Min(float64(ManualDecrementAmount), w.matrix[x][y])
				}

				w.matrix[x][y] = w.matrix[x][y] - decrementAmount
				w.StoredValue = w.StoredValue + decrementAmount
			} else {
				w.matrix[x][y] = w.matrix[x][y] - ManualDecrementAmount
			}
		}

		// TODO: place a permanent incrementor
		// TODO: place a permanent decrementer

		// toggle random incremements
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			keySpaceDown = true
		} else {
			// If space key was previously down, toggle random incremements
			if keySpaceDown {
				randomIncrements = !randomIncrements
			}
			keySpaceDown = false
		}

		// Reset
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			keyRDown = true
		} else {
			// If R key was previously down, reset
			if keyRDown {
				w.Reset()
			}
			keyRDown = false
		}

	}

}

// TotalValue returns the total value of all cells
func (w *World) TotalValue() float64 {
	val := float64(0)

	for col := range w.matrix {
		for row := range w.matrix[col] {
			val = val + w.matrix[col][row]
		}
	}

	return val
}

// TODO: waves?
// i.e. when calculating spill, consider previous spill in that direction
// Or, perhaps track not only value, but pressure
// higher pressure = higher spill

// TODO: this is horrible; refactor it
// spillToNeighbours will take some of a cell's value and spill it out to neighbours
func (w *World) spillToNeighbours() {

	matrix := w.matrix

	// Each row has the same number of cols, so use first col
	numRows := w.height
	numCols := w.width

	// Dupe the matrix, to store spill values
	spillMatrix := make([][]float64, numCols)
	for i := range spillMatrix {
		spillMatrix[i] = make([]float64, numRows)
	}

	// TODO: keep track of non-empty cells

	for col := range matrix {
		for row := range matrix[col] {
			cellValue := matrix[col][row]

			//
			// First, check if we can spill anywhere
			//

			// Assume we cannot spill anywhere
			var spillDirections uint8 = 0
			var spillEastAmount float64 = 0
			var spillWestAmount float64 = 0
			var spillNorthAmount float64 = 0
			var spillSouthAmount float64 = 0

			// TODO: loop?

			// Check north
			if row > 0 {
				cellValueNorth := matrix[col][row-1]

				// If the value is lower than ours...
				if cellValueNorth < cellValue {
					spillDirections++

					diffNorth := cellValue - cellValueNorth

					// Spill at most half the diff
					spillNorthAmount = diffNorth * SpillRate / 2
				}
			}

			// Check south
			if row < numRows-1 {
				cellValueSouth := matrix[col][row+1]

				// If the value is lower than ours...
				if cellValueSouth < cellValue {
					spillDirections++

					diffSouth := cellValue - cellValueSouth

					// Spill at most half the diff
					spillSouthAmount = diffSouth * SpillRate / 2
				}
			}

			// Check east
			if col < numCols-1 {
				cellValueEast := matrix[col+1][row]

				// If the value is lower than ours...
				if cellValueEast < cellValue {
					spillDirections++

					diffEast := cellValue - cellValueEast

					// Spill at most half the diff
					spillEastAmount = diffEast * SpillRate / 2
				}
			}

			// Check west
			if col > 0 {
				cellValueWest := matrix[col-1][row]

				// If the value is lower than ours...
				if cellValueWest < cellValue {
					spillDirections++

					diffWest := cellValue - cellValueWest

					// Spill at most half the diff
					spillWestAmount = diffWest * SpillRate / 2
				}
			}

			//
			// keep track of spill amounts in a separate matrix
			// otherwise there will be a kind of prevailing wind, in the south easterly direction
			//

			if spillNorthAmount > 0 {
				// Spill North
				spillNorthAmount = spillNorthAmount / float64(spillDirections)
				spillMatrix[col][row] = spillMatrix[col][row] - spillNorthAmount
				spillMatrix[col][row-1] = spillMatrix[col][row-1] + spillNorthAmount
			}

			if spillSouthAmount > 0 {
				// Spill South
				spillSouthAmount = spillSouthAmount / float64(spillDirections)
				spillMatrix[col][row] = spillMatrix[col][row] - spillSouthAmount
				spillMatrix[col][row+1] = spillMatrix[col][row+1] + spillSouthAmount
			}

			if spillEastAmount > 0 {
				// Spill East
				spillEastAmount = spillEastAmount / float64(spillDirections)
				spillMatrix[col][row] = spillMatrix[col][row] - spillEastAmount
				spillMatrix[col+1][row] = spillMatrix[col+1][row] + spillEastAmount
			}

			if spillWestAmount > 0 {
				// Spill West
				spillWestAmount = spillWestAmount / float64(spillDirections)
				spillMatrix[col][row] = spillMatrix[col][row] - spillWestAmount
				spillMatrix[col-1][row] = spillMatrix[col-1][row] + spillWestAmount
			}

		}
	}

	//
	// Spill!
	//

	for col := range matrix {
		for row := range matrix[col] {
			cellValue := matrix[col][row]
			spillValue := spillMatrix[col][row]

			matrix[col][row] = cellValue + spillValue
		}
	}

}
