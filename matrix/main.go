package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	tm "github.com/buger/goterm"
	"github.com/mgutz/ansi"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	// Display stuff; controls how big the matrix is
	CellWidth  = 2
	CellBuffer = 1

	// How many of the cells will increment every tick
	// TODO: this should be a variable, inversely proportional to the sum of the matrix
	ChanceOfIncrement = 0.0005

	// How much of a cell's value should spill out to neighbours every tick
	// This is in relation to half of the difference between the cell's value and the neighbour's value
	SpillRate = 1

	// TODO: variable, how frequently per second should something update
	// or, how many things should update per second
	// Use this to calcualte UpdatesPerSecond

	// How much will it increment by
	IncrementAmount = 255

	// How many ticks to attempt every second
	// This can be super high, but would then essentially fill up immediatelly
	TicksPerSecond = 24 // TODO: allow this to be < 1

	// How often to refresh the view
	FramesPerSecond = 24

	// Whether to display cell value or not
	DisplayValue = false

	// Whether to enable colors for cells
	EnableColor = true
)

var (
	// what's the biggest hex number we can display?
	MaxDisplay = int(math.Pow(16, CellWidth))
	FullString = strings.Repeat("X", CellWidth)

	matrix [][]float32

	resetColorCode = ansi.ColorCode("reset")
	// Color Codes
	minColor int = 232
	maxColor int = 255
	midColor int = 243 // white below, black above
)

func transformValToColor(cellValue float32) (int, string) {

	// https://math.stackexchange.com/a/377174
	colorCode := int(cellValue*(float32(maxColor-minColor)/float32(MaxDisplay)) + float32(minColor))

	// make sure it's within range
	colorCode = int(math.Max(float64(colorCode), float64(minColor)))
	colorCode = int(math.Min(float64(colorCode), float64(maxColor)))

	fg := "white"

	if colorCode > midColor {
		fg = "black"
	}

	return colorCode, fg

}

// cellColorString
func cellString(cellValue float32) string {

	cellString := ""

	if DisplayValue {
		// if bigger than we can display
		if int(cellValue) >= MaxDisplay {
			cellString = fmt.Sprintf(" %s", FullString)
		} else {
			cellString = fmt.Sprintf(" %*x", CellWidth, int(cellValue))
		}
	} else {
		cellString = fmt.Sprintf(" %s", strings.Repeat(" ", CellWidth))
	}

	if EnableColor {
		cellColor, fg := transformValToColor(float32(cellValue))
		color := ansi.ColorCode(fmt.Sprintf("%s:%d", fg, cellColor))
		cellString = color + cellString + resetColorCode
	}

	return cellString

}

// printMatrix just prints it to screen
func printMatrix(matrix [][]float32) {

	var buffer bytes.Buffer

	// print the matrix

	// TODO: concurrency?
	// e.g. generate the strings, 1 row per thread, then dump them all out in one go
	for row := range matrix {
		for col := range matrix[row] {
			cellValue := matrix[row][col]

			buffer.WriteString(cellString(cellValue))
		}
		buffer.WriteString("\n")
	}

	fmt.Println(buffer.String())

}

// randomIncrementMatrix will, for every cell, with a certain percentage, either increment or not
func randomIncrementMatrix(matrix [][]float32) {

	// Each row has the same number of cols, so use first col
	numRows := len(matrix)
	numCols := len(matrix[1])
	numCells := numRows * numCols

	cellsToIncrement := int(float32(numCells) * ChanceOfIncrement)
	// TODO: what happens if this is zero?
	if cellsToIncrement == 0 {
		cellsToIncrement = 1
		// For now, this will do
		// In this case, we need to use this to determine if ANY cells increment
	}

	for i := 0; i < cellsToIncrement; i++ {
		randCellRow := rand.Intn(numRows)
		randCellCol := rand.Intn(numCols)

		matrix[randCellRow][randCellCol] = matrix[randCellRow][randCellCol] + IncrementAmount
	}

}

// spillToNeighbours will take some of a cell's value and spill it out to neighbours
func spillToNeighbours(matrix [][]float32) {

	// Each row has the same number of cols, so use first col
	numRows := len(matrix)
	numCols := len(matrix[1])

	// Dupe the matrix, to store spill values
	spillMatrix := make([][]float32, numRows)
	for i := range spillMatrix {
		spillMatrix[i] = make([]float32, numCols)
	}

	// TODO: keep track of non-empty cells

	for row := range matrix {
		for col := range matrix[row] {
			cellValue := matrix[row][col]

			//
			// First, check if we can spill anywhere
			//

			// Assume we cannot spill anywhere
			var spillDirections uint8 = 0
			var spillEastAmount float32 = 0
			var spillWestAmount float32 = 0
			var spillNorthAmount float32 = 0
			var spillSouthAmount float32 = 0

			// Check north
			if row > 0 {
				cellValueNorth := matrix[row-1][col]

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
				cellValueSouth := matrix[row+1][col]

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
				cellValueEast := matrix[row][col+1]

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
				cellValueWest := matrix[row][col-1]

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
				spillMatrix[row][col] = spillMatrix[row][col] - spillNorthAmount
				spillMatrix[row-1][col] = spillMatrix[row-1][col] + spillNorthAmount
			}

			if spillSouthAmount > 0 {
				// Spill South
				spillSouthAmount = spillSouthAmount / float32(spillDirections)
				spillMatrix[row][col] = spillMatrix[row][col] - spillSouthAmount
				spillMatrix[row+1][col] = spillMatrix[row+1][col] + spillSouthAmount
			}

			if spillEastAmount > 0 {
				// Spill East
				spillEastAmount = spillEastAmount / float32(spillDirections)
				spillMatrix[row][col] = spillMatrix[row][col] - spillEastAmount
				spillMatrix[row][col+1] = spillMatrix[row][col+1] + spillEastAmount
			}

			if spillWestAmount > 0 {
				// Spill West
				spillWestAmount = spillWestAmount / float32(spillDirections)
				spillMatrix[row][col] = spillMatrix[row][col] - spillWestAmount
				spillMatrix[row][col-1] = spillMatrix[row][col-1] + spillWestAmount
			}

		}
	}

	//
	// Spill!
	//

	for row := range matrix {
		for col := range matrix[row] {
			cellValue := matrix[row][col]
			spillValue := spillMatrix[row][col]

			matrix[row][col] = cellValue + spillValue
		}
	}

}

// tick is the loop where stuff actually happens
func tick() {

	for {
		// Distribute some of your value to neighbours
		spillToNeighbours(matrix)

		// Random Increment
		randomIncrementMatrix(matrix)

		if TicksPerSecond < 1 {
			SecondsPerTick := 1 / TicksPerSecond

			time.Sleep(time.Second * time.Duration(SecondsPerTick))
		} else {
			time.Sleep(time.Second / TicksPerSecond)
		}
	}
}

func init() {
	// Not going to bother handling screen resizing
	fd := int(os.Stdin.Fd())
	tw, th, _ := terminal.GetSize(fd)

	rows := th - 3
	cols := tw / (CellWidth + CellBuffer)

	// Create a matrix of integers, same dimensions as the terminal
	matrix = make([][]float32, rows)
	for i := range matrix {
		matrix[i] = make([]float32, cols)
	}

	// Seed the RNG
	rand.Seed(time.Now().UnixNano())

}

func main() {
	tm.Clear()                   // Clear current screen
	fmt.Print("\033[?25l")       // hide cursor
	defer fmt.Print("\033[?25h") //unhide cursor

	go tick()

	for {
		// By moving cursor to top-left position we ensure that console output
		// will be overwritten each time, instead of adding new.
		tm.MoveCursor(1, 1)

		// Display
		printMatrix(matrix)

		tm.Flush() // Call it every time at the end of rendering
		time.Sleep(time.Second / FramesPerSecond)

	}

}
