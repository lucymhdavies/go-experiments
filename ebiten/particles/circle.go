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
)

type Color struct {
	R, G, B, A float64
}

type Circle struct {
	Pos    r3.Vector
	Radius float64
	Color  Color
}

func NewCircle(x, y, r float64) Circle {
	return Circle{
		Pos: r3.Vector{
			X: x,
			Y: y,
		},
		Radius: r,
		Color:  colorWhite,
	}
}

var (
	dot *ebiten.Image
)

func (c Circle) Draw(screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(c.Color.R, c.Color.G, c.Color.B, c.Color.A)
	op.GeoM.Reset()
	op.GeoM.Translate(c.Pos.X, c.Pos.Y)

	// TODO: draw circle of radius...
	// See: https://github.com/ae6rt/golang-examples/blob/master/goeg/src/shaper2/shapes/shapes.go

	screen.DrawImage(dot, op)

	return nil
}
