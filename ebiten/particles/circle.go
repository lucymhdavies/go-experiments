package main

import (
	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"
)

// type Vector struct {
// 	X, Y float64
// }

var (
	colorWhite = Color{
		R: 1.0,
		G: 1.0,
		B: 1.0,
		A: 0.5,
	}

	colorGreen = Color{
		R: 0.0,
		G: 1.0,
		B: 0.0,
		A: 1.0,
	}
)

type Color struct {
	R, G, B, A float64
}

type Circle struct {
	Pos     r3.Vector
	Radius  float64
	Color   Color
	Visible bool
}

func NewCircle(x, y, r float64) Circle {
	return Circle{
		Pos: r3.Vector{
			X: x,
			Y: y,
		},
		Radius:  r,
		Color:   colorWhite,
		Visible: true,
	}
}

var (
	dot *ebiten.Image
)

func (c Circle) Draw(screen *ebiten.Image) error {
	if !c.Visible {
		return nil
	}

	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(c.Color.R, c.Color.G, c.Color.B, c.Color.A)
	op.GeoM.Reset()

	if c.Radius > 1 {
		op.GeoM.Scale(c.Radius*2+1, c.Radius*2+1)
		op.GeoM.Translate(c.Pos.X-c.Radius, c.Pos.Y-c.Radius)
	} else {
		// No scaling. Single dot. Just draw it
		op.GeoM.Translate(c.Pos.X, c.Pos.Y)
	}

	// TODO: draw circle of radius...
	// See: https://github.com/ae6rt/golang-examples/blob/master/goeg/src/shaper2/shapes/shapes.go

	screen.DrawImage(dot, op)

	return nil
}
