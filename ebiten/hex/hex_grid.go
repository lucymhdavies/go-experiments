package main

import "github.com/hajimehoshi/ebiten"

//
// Hex Grid
//

type HexGrid struct {
	width  int
	height int

	// TODO: Store this in a structure that allows searching based on coords
	tiles []*HexTile
}

func NewHexGrid(width, height int) *HexGrid {

	tiles := []*HexTile{}

	// 0,0 is center
	halfWidth := (width - 1) / 2
	halfHeight := (height - 1) / 2

	// If width/height is even, round up
	if width%2 == 0 {
		halfWidth++
	}
	if height%2 == 0 {
		halfHeight++
	}

	for x := -halfWidth; x <= halfWidth; x++ {
		for y := -halfHeight; y <= halfHeight; y++ {

			// odd rows should be 1 tile shorter than even rows
			if y%2 == 1 || y%2 == -1 {
				if x == halfWidth {
					continue
				}
			}

			tile := NewHexTile(x, y)

			tiles = append(tiles, tile)
		}
	}

	return &HexGrid{
		width:  width,
		height: height,

		tiles: tiles,
	}
}

func (g *HexGrid) Draw(screen *ebiten.Image) {

	for _, tile := range g.tiles {
		tile.Draw(screen)
	}
}
