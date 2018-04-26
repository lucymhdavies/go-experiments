package main

import (
	"image/color"
	"math"
)

type Color = color.RGBA

const (
	// Early version of transformValToColor introduced a weird rendering error
	// Keeping it available as an option, because it looks cool
	colorWeird = iota

	// An HSV spiral type thing
	// i.e. start at black, rotate around hue, and gradually increase saturation
	// TODO
	colorSpiral

	// normal, single color mode
	colorNormal
)

// transformValToColor returns a byte 4-tuple, denoting a color code
func Float32ToColor(cellValue float32) Color {

	switch ColorMode {
	case colorWeird:
		return float32ToColorWeird(cellValue)

	default:
		return float32ToColorNormal(cellValue)
	}

}

func float32ToColorWeird(cellValue float32) Color {
	// https://math.stackexchange.com/a/377174
	color := byte(cellValue * (255 / float32(MaxDisplay)))

	// make sure it's within range
	color = byte(math.Max(float64(color), 0))
	color = byte(math.Min(float64(color), 255))

	return Color{
		R: 0,
		G: byte(color),
		B: byte(color),
		A: 255,
	}
}

func float32ToColorNormal(cellValue float32) Color {
	// https://math.stackexchange.com/a/377174
	color := float32(cellValue * (255 / float32(MaxDisplay)))

	// make sure it's within range
	color = float32(math.Max(float64(color), 0))
	color = float32(math.Min(float64(color), 255))

	return Color{
		R: 0,
		G: byte(color),
		B: byte(color),
		A: 255,
	}
}
