package main

// TODO: use a uint8
// That way, we know we have a max of 255 (which is more than enough for this, I reckon)
// And we don't need to calculate cell width on the fly

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	tm "github.com/buger/goterm"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	// Display stuff; controls how big the matrix is
	CellWidth  = 2
	CellBuffer = 1

	// How many of the cells will increment every tick
	ChanceOfIncrement = 0.01 // TODO: this should be a variable, inversely proportional to the sum of the matrix

	// How much will it increment by
	IncrementAmount = 1

	// How many ticks to attempt every second
	// This can be super high, but would then essentially fill up immediatelly
	UpdatesPerSecond = 600

	// How often to refresh the view
	FramesPerSecond = 30
)

var (
	// what's the biggest hex number we can display?
	MaxDisplay = int(math.Pow(16, CellWidth))
	FullString = strings.Repeat("X", CellWidth)

	matrix [][]int
)

// printMatrix just prints it to screen
func printMatrix(matrix [][]int) {

	// print the matrix

	// TODO: concurrency?
	// e.g. generate the strings, 1 row per thread, then dump them all out in one go
	for row := range matrix {
		for col := range matrix[row] {
			cellValue := matrix[row][col]
			cellString := ""

			// if bigger than we can display
			if cellValue >= MaxDisplay {
				cellString = fmt.Sprintf(" %s", FullString)
			} else {
				cellString = fmt.Sprintf(" %*x", CellWidth, cellValue)
			}

			fmt.Printf(cellString)
		}
		fmt.Print("\n")
	}
}

// randomIncrementMatrix will, for every cell, with a certain percentage, either increment or not
func randomIncrementMatrix(matrix [][]int) {

	// Each row has the same number of cols, so use first col
	numRows := len(matrix)
	numCols := len(matrix[1])
	numCells := numRows * numCols

	cellsToIncrement := int(float32(numCells) * ChanceOfIncrement)
	// TODO: what happens if this is zero?

	for i := 0; i < cellsToIncrement; i++ {
		randCellRow := rand.Intn(numRows)
		randCellCol := rand.Intn(numCols)

		matrix[randCellRow][randCellCol] = matrix[randCellRow][randCellCol] + IncrementAmount
	}

}

// stuffHappens is the loop where stuff actually happens
func stuffHappens() {

	for {
		// TODO
		// Distribute some of your value to neighbours
		// Requires keeping track of non-empty cells

		// Random Increment
		randomIncrementMatrix(matrix)

		time.Sleep(time.Second / UpdatesPerSecond)
	}
}

func init() {
	// Not going to bother handling screen resizing
	fd := int(os.Stdin.Fd())
	tw, th, _ := terminal.GetSize(fd)

	rows := th * 3 / 4
	cols := tw / (CellWidth + CellBuffer)

	// Create a matrix of integers, same dimensions as the terminal
	matrix = make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
	}

	// Seed the RNG
	rand.Seed(time.Now().UnixNano())

}

func main() {
	tm.Clear()                   // Clear current screen
	fmt.Print("\033[?25l")       // hide cursor
	defer fmt.Print("\033[?25h") //unhide cursor

	go stuffHappens()

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
