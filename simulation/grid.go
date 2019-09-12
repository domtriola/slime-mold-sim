package simulation

import (
	"math/rand"
)

// Space represents a discrete location
type Space struct {
	scent uint16
}

// Organism represents one creature that occupies one space
type Organism struct {
	direction int
}

// Grid represents the simulation area
type Grid struct {
	width  int
	height int
	rows   [][]*Space
}

// OrganismGrid is a two dimensional array containing all organisms at their
// corresponding location on the grid
type OrganismGrid struct {
	width  int
	height int
	rows   [][]*Organism
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
			row = append(row, &space)
		}
		grid.rows = append(grid.rows, row)
	}
}

func (organisms *OrganismGrid) initialize() {
	for y := 0; y < options["height"]; y++ {
		row := []*Organism{}
		for x := 0; x < options["width"]; x++ {
			// TODO: make dynamic
			if rand.Intn(10) == 1 {
				organism := Organism{
					direction: rand.Intn(360),
				}
				row = append(row, &organism)
			} else {
				row = append(row, nil)
			}
		}
		organisms.rows = append(organisms.rows, row)
	}
}
