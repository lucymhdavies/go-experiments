package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	tm "github.com/buger/goterm"
	// 	"github.com/mgutz/ansi"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	CellWidth         = 2
	CellBuffer        = 1
	ChanceOfIncrement = 0.01 // TODO: this should be a variable, inversely proportional to the sum of the matrix
	IncrementAmount   = 1
	UpdatesPerSecond  = 100
)

var (
	// what's the biggest hex number we can display?
	MaxDisplay = int(math.Pow(16, CellWidth))
	FullString = strings.Repeat("X", CellWidth)
)

// // getColor maps an integer to an ansi color code
// func getColor(i int) string {
// 	// TODO: use grayscale
// 	// 232 - 255
// 	// <= 243, white fg
// 	// else black fg
//
// 	// Some ansi colors from https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
//
// 	if i < MaxDisplay*1/6 {
// 		return "white:16"
// 	}
//
// 	if i < MaxDisplay*2/6 {
// 		return "white:17"
// 	}
//
// 	if i < MaxDisplay*3/6 {
// 		return "white:18"
// 	}
//
// 	if i < MaxDisplay*4/6 {
// 		return "white:19"
// 	}
//
// 	if i < MaxDisplay*5/6 {
// 		return "white:20"
// 	}
//
// 	return "white:21"
// }

// printMatrix just prints it to screen
func printMatrix(matrix [][]int) {

	// 	reset := ansi.ColorCode("reset")

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

			// 			// TODO: find a better color library. this one is slow
			// 			color := ansi.ColorCode(getColor(cellValue))
			// 			fmt.Printf("%s%s%s", color, fmt.Sprintf("%s", cellString), reset)

			fmt.Printf(cellString)
		}
		fmt.Print("\n")
	}
}

// randomIncrementMatrix will, for every cell, with a certain percentage, either increment or not
func randomIncrementMatrix(matrix [][]int) {

	// TODO: better idea would be to pick a random element from the matrix and increment it
	// https://stackoverflow.com/a/33994787

	// TODO: do this concurrently, 1 thread per row?
	for row := range matrix {
		for col := range matrix[row] {
			rand := rand.Float32()
			if rand < ChanceOfIncrement {

				matrix[row][col] = matrix[row][col] + IncrementAmount
			}
		}
	}
}

func main() {

	// Not going to bother handling screen resizing
	fd := int(os.Stdin.Fd())
	tw, th, _ := terminal.GetSize(fd)

	rows := th * 3 / 4
	cols := tw / (CellWidth + CellBuffer)

	// Create a matrix of integers, same dimensions as the terminal
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
	}

	tm.Clear()             // Clear current screen
	fmt.Print("\033[?25l") // hide cursor
	for {
		// By moving cursor to top-left position we ensure that console output
		// will be overwritten each time, instead of adding new.
		tm.MoveCursor(1, 1)

		// Display
		printMatrix(matrix)

		// TODO
		// Distribute some of your value to neighbours
		// Requires keeping track of non-empty cells

		// Random Increment
		randomIncrementMatrix(matrix)

		tm.Flush() // Call it every time at the end of rendering
		time.Sleep(time.Second / UpdatesPerSecond)
	}

}
