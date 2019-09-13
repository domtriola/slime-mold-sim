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
	scent    uint16
	organism *Organism
}

// Organism represents one creature that occupies one space
type Organism struct {
	id        int
	direction float64

	// Position is a float so the organisms can travel at angles other than 45
	// degrees while still taking steps of one cell at a time.
	xPos float64
	yPos float64

	nextDiscreteXPos int
	nextDiscreteYPos int
}

func (grid Grid) hasCoord(x, y int) bool {
	return x > 0 && x < grid.width && y > 0 && y < grid.height
}

func (grid Grid) get(x, y int) *Space {
	return grid.rows[y][x]
}

func (grid *Grid) initialize() {
	orgCount := 0
	for y := 0; y < options["height"]; y++ {
		row := []*Space{}
		for x := 0; x < options["width"]; x++ {
			space := Space{}

			if rand.Intn(50) == 1 {
				organism := Organism{
					id:        orgCount,
					xPos:      float64(x),
					yPos:      float64(y),
					direction: float64(rand.Intn(360)),
				}

				space.organism = &organism
				orgCount++
			}

			row = append(row, &space)
		}
		grid.rows = append(grid.rows, row)
	}
}
