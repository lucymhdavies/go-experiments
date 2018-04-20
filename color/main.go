package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/mgutz/ansi"
)

const (
	CellWidth = 2
)

var (
	MaxDisplay = int(math.Pow(16, CellWidth))

	minColor int = 232
	maxColor int = 255
	midColor int = 243 // white below, black above
)

func main() {

	// cache escape codes and build strings manually
	reset := ansi.ColorCode("reset")

	for i := minColor; i <= maxColor; i++ {

		fg := "white"

		if i > 243 {
			fg = "black"
		}

		color := ansi.ColorCode(fmt.Sprintf("%s:%d", fg, i))
		fmt.Println(color, fmt.Sprintf("%d", i), reset)
	}

	for {
		cellValue := rand.Intn(MaxDisplay)

		cellColor, fg := transformValToColor(float32(cellValue))

		color := ansi.ColorCode(fmt.Sprintf("%s:%d", fg, cellColor))
		fmt.Println(color, fmt.Sprintf("%d --> %d", cellValue, cellColor), reset)

		time.Sleep(time.Second)
	}
}

func transformValToColor(cellValue float32) (int, string) {

	// https://math.stackexchange.com/a/377174
	colorCode := int(cellValue*(float32(maxColor-minColor)/float32(MaxDisplay)) + float32(minColor))

	fg := "white"

	if colorCode > midColor {
		fg = "black"
	}

	return colorCode, fg

}
