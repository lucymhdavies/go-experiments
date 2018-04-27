package main

import (
	"image/color"
	"math"
)

type Color = color.RGBA

// Color Modes
const (
	colorWeird = iota
	colorSpiral
	colorNormal
)

// transformValToColor returns a byte 4-tuple, denoting a color code
func Float32ToColor(cellValue float32) Color {

	switch ColorMode {

	case colorWeird:
		return float32ToColorWeird(cellValue)

	case colorSpiral:
		return float32ToColorSpiral(cellValue)

	default:
		return float32ToColorNormal(cellValue)

	}

}

// Early version of transformValToColor introduced a weird rendering error
// Keeping it available as an option, because it looks cool
func float32ToColorWeird(cellValue float32) Color {

	// Originally, this "worked" on desktop, due to byte overflow.
	// In GopherJS, this just behaves like Normal, probably because of how
	// Javascript only has a single Number type.
	//
	// Simulate the original byte overflow with a modulo
	cellValue = float32(math.Mod(float64(cellValue), CycleValue))

	// https://math.stackexchange.com/a/377174
	color := byte(cellValue * (255 / float32(CycleValue)))

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

// An HSV spiral type thing
// i.e. start at black, rotate around hue, and gradually increase value
func float32ToColorSpiral(cellValue float32) Color {

	// First, need to be > 0
	cellValue = float32(math.Max(float64(cellValue), 0))

	// mod by MaxValue?

	// Now map 0-Max to 0-1
	// MaxDisplay in this case is the value at which we get a new cycle
	// Modulo 1, so we're always between 0 and 1
	hue := math.Mod(float64(cellValue*(1/float32(CycleValue))), 1)

	// Do the same for Val, but it takes 5 times as long to cycle
	val := math.Mod(float64(cellValue*(SpiralValueCycleRatio/float32(CycleValue))), 1)

	// TODO: do the same for sat, but take 25 times as long?
	// also, 1-sat

	colorHSV := ColorHSV{
		S: 1, // keep static
		H: hue,
		V: val,
	}

	color := colorHSV.RGB()

	return Color{
		R: color.R,
		G: color.G,
		B: color.B,
		A: 255,
	}
}

// normal, single color mode
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

// From https://play.golang.org/p/9q5yBNDh3W
// No idea who the original author is
type ColorHSV struct {
	H, S, V float64
}

func (c *ColorHSV) RGB() *Color {
	var r, g, b float64
	if c.S == 0 { //HSV from 0 to 1
		r = c.V * 255
		g = c.V * 255
		b = c.V * 255
	} else {
		h := c.H * 6
		if h == 6 {
			h = 0
		} //H must be < 1
		i := math.Floor(h) //Or ... var_i = floor( var_h )
		v1 := c.V * (1 - c.S)
		v2 := c.V * (1 - c.S*(h-i))
		v3 := c.V * (1 - c.S*(1-(h-i)))

		if i == 0 {
			r = c.V
			g = v3
			b = v1
		} else if i == 1 {
			r = v2
			g = c.V
			b = v1
		} else if i == 2 {
			r = v1
			g = c.V
			b = v3
		} else if i == 3 {
			r = v1
			g = v2
			b = c.V
		} else if i == 4 {
			r = v3
			g = v1
			b = c.V
		} else {
			r = c.V
			g = v1
			b = v2
		}

		r = r * 255 //RGB results from 0 to 255
		g = g * 255
		b = b * 255
	}
	rgb := &Color{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255,
	}
	return rgb

}
