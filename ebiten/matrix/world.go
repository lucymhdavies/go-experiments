package main

import (
	"fmt"
	"math"
	"math/rand"
)

// TODO:
// we currently have x,y coords, and a single value
// how about x,y coords, and r,g,b values?

type Matrix = [][]float32

type World struct {
	matrix        Matrix // [x][y]
	width, height int    // convenience vars
}

// newMatrix initialises a new empty matrix
func newMatrix(width, height int) Matrix {
	a := make(Matrix, width)
	for i := 0; i < width; i++ {
		a[i] = make([]float32, height)
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

// Update game state by one tick.
func (w *World) Update() error {

	// Randomly Increment
	w.spillToNeighbours()
	w.randomIncrementMatrix()

	return nil
}

// Draw paints current game state.
func (w *World) Draw(pix []byte) {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			idx := 4*y*w.width + 4*x

			color := transformValToColor(w.matrix[x][y])

			// Red
			pix[idx] = 0
			// Green
			pix[idx+1] = color
			// Blue
			pix[idx+2] = color
			// Alpha
			pix[idx+3] = color
		}
	}
}

func transformValToColor(cellValue float32) byte {

	// 	// https://math.stackexchange.com/a/377174
	// 	colorCode := byte(cellValue * (255 / float32(MaxDisplay)))
	//
	// 	// make sure it's within range
	// 	colorCode = byte(math.Max(float64(colorCode), 0))
	// 	colorCode = byte(math.Min(float64(colorCode), 255))

	// https://math.stackexchange.com/a/377174
	colorCode := float32(cellValue * (255 / float32(MaxDisplay)))

	// make sure it's within range
	colorCode = float32(math.Max(float64(colorCode), 0))
	colorCode = float32(math.Min(float64(colorCode), 255))

	return byte(colorCode)

}

// Get val at specified coord, converted to int
func (w *World) GetVal(x, y int) (float32, error) {
	m := w.matrix

	if x >= 0 && x < w.width && y >= 0 && y < w.height {
		return m[x][y], nil
	}

	return -1, fmt.Errorf("Out of Bounds")
}

// randomIncrementMatrix will, for every cell, with a certain percentage, either increment or not
func (w *World) randomIncrementMatrix() {

	// TODO: do this in a for loop, as per original version

	cellsToIncrement := 0
	if rand.Float32() < ChanceOfIncrement {
		cellsToIncrement = 1
	}

	for i := 0; i < cellsToIncrement; i++ {
		randY := rand.Intn(w.height)
		randX := rand.Intn(w.width)

		w.matrix[randX][randY] = w.matrix[randX][randY] + IncrementAmount
	}

}

// TODO: this is horrible; refactor it
// spillToNeighbours will take some of a cell's value and spill it out to neighbours
func (w *World) spillToNeighbours() {

	matrix := w.matrix

	// Each row has the same number of cols, so use first col
	numRows := w.height
	numCols := w.width

	// Dupe the matrix, to store spill values
	spillMatrix := make([][]float32, numCols)
	for i := range spillMatrix {
		spillMatrix[i] = make([]float32, numRows)
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
			var spillEastAmount float32 = 0
			var spillWestAmount float32 = 0
			var spillNorthAmount float32 = 0
			var spillSouthAmount float32 = 0

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
				spillNorthAmount = spillNorthAmount / float32(spillDirections)
				spillMatrix[col][row] = spillMatrix[col][row] - spillNorthAmount
				spillMatrix[col][row-1] = spillMatrix[col][row-1] + spillNorthAmount
			}

			if spillSouthAmount > 0 {
				// Spill South
				spillSouthAmount = spillSouthAmount / float32(spillDirections)
				spillMatrix[col][row] = spillMatrix[col][row] - spillSouthAmount
				spillMatrix[col][row+1] = spillMatrix[col][row+1] + spillSouthAmount
			}

			if spillEastAmount > 0 {
				// Spill East
				spillEastAmount = spillEastAmount / float32(spillDirections)
				spillMatrix[col][row] = spillMatrix[col][row] - spillEastAmount
				spillMatrix[col+1][row] = spillMatrix[col+1][row] + spillEastAmount
			}

			if spillWestAmount > 0 {
				// Spill West
				spillWestAmount = spillWestAmount / float32(spillDirections)
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
