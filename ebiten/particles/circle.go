package main

import (
	"github.com/golang/geo/r3"
	"github.com/hajimehoshi/ebiten"
)

// type Vector struct {
// 	X, Y float64
// }

type Circle struct {
	Pos r3.Vector
	R   float64

	// TODO: color
}

func NewCircle(x, y, r float64) Circle {
	return Circle{
		Pos: r3.Vector{
			X: x,
			Y: y,
		},
		R: r,
	}
}

var (
	dot *ebiten.Image
)

func (c Circle) Draw(screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(200.0/255.0, 200.0/255.0, 200.0/255.0, 1)
	op.GeoM.Reset()
	op.GeoM.Translate(c.Pos.X, c.Pos.Y)

	// TODO: draw circle of radius...
	// See: https://github.com/ae6rt/golang-examples/blob/master/goeg/src/shaper2/shapes/shapes.go

	screen.DrawImage(dot, op)

	return nil
}
