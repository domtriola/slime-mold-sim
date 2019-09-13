package simulation

import (
	"math/rand"
)

// Grid represents the simulation area
type Grid struct {
	width  int
	height int
	rows   [][]*Space
}

// Space represents a discrete location
type Space struct {
	scent     uint16
	organisms []*Organism
}

// Organism represents one creature that occupies one space
type Organism struct {
	direction int

	// Position is a float so the organisms can travel at angles other than 45
	// degrees while still taking steps of one cell at a time.
	xPos float32
	yPos float32
}

func (grid Grid) hasCoord(x, y int) bool {
	return x > 0 && x < grid.width && y > 0 && y < grid.height
}

func (grid Grid) get(x, y int) *Space {
	return grid.rows[y][x]
}

func (grid *Grid) initialize() {
	for y := 0; y < options["height"]; y++ {
		row := []*Space{}
		for x := 0; x < options["width"]; x++ {
			space := Space{}

			// TODO: make dynamic
			if rand.Intn(10) == 1 {
				organism := Organism{
					direction: rand.Intn(360),
				}

				space.organisms = append(space.organisms, &organism)
			}

			row = append(row, &space)
		}
		grid.rows = append(grid.rows, row)
	}
}
