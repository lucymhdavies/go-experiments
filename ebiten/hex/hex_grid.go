package main

import (
	"math"

	"github.com/hajimehoshi/ebiten"
)

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

func (g *HexGrid) FindNearestTile(x, y float64) (nearestTile *HexTile, mouseOverTile bool) {

	distance := -1.0

	for _, tile := range g.tiles {
		tile.highlighted = false

		tileDistanceX := x - tile.screenMidX
		tileDistanceY := y - tile.screenMidY

		tileDistance := math.Sqrt(
			math.Pow(tileDistanceX, 2) +
				math.Pow(tileDistanceY, 2))

		if distance == -1.0 || tileDistance < distance {
			distance = tileDistance
			nearestTile = tile
		}
	}

	sizeX, sizeY := hexImage.Size()

	// if the mouse is actually over the nearest tile
	// TODO: this currently results in not highlighting tiles
	// when the mouse is near a 3-tile boundary
	if distance < float64(sizeX)/2 && distance <= float64(sizeY)/2 {
		return nearestTile, true
	}

	return nearestTile, false
}

// TODO: get tile neighbours
