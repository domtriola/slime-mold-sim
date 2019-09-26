package simulation

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math/rand"
	"os"
)

const (
	x0 = 0
	y0 = 0
)

var scentDecay = 0.9
var scentSpreadFactor = 0.1

var options = map[string]int{
	"width":          500,
	"height":         500,
	"nFrames":        500,
	"loopCount":      1000,
	"delay":          2,
	"sensorDegree":   45,
	"sensorDistance": 9,
}

// Build creates the simulation as a GIF
func Build(urlOptions map[string]interface{}) (name string) {
	setOptions(urlOptions)

	name = filename()
	animate(name)
	return name
}

func filename() (name string) {
	path := "tmp"
	extension := "gif"
	return fmt.Sprintf(
		"%s/w%dh%dnF%dd%dsDe%dsDi%d.%s",
		path,
		options["width"],
		options["height"],
		options["nFrames"],
		options["delay"],
		options["sensorDegree"],
		options["sensorDistance"],
		extension,
	)
}

func setOptions(urlOptions map[string]interface{}) {
	for option := range options {
		if val, ok := urlOptions[option]; ok {
			options[option] = val.(int)
		}
	}
}

func animate(name string) {
	grid := Grid{width: options["width"], height: options["height"]}
	anim := gif.GIF{LoopCount: options["loopCount"]}

	var pal color.Palette
	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}
	pal = append(pal, black)
	pal = append(pal, white)

	palette := []color.Color{
		color.RGBA{0, 0, 0, 255},
	}
	for greyScale := 0; greyScale <= 255; greyScale += 10 {
		palette = append(palette, color.RGBA{
			uint8(greyScale),
			uint8(greyScale),
			uint8(greyScale),
			255,
		})
	}

	grid.initialize()
	for i := 0; i < options["nFrames"]; i++ {
		drawNextFrame(grid, &anim, palette)
	}

	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	gif.EncodeAll(f, &anim)
}

func createImage(grid Grid, pal color.Palette) (img *image.Paletted) {
	rect := image.Rect(x0, y0, options["width"], options["height"])
	img = image.NewPaletted(rect, pal)

	for y, row := range grid.rows {
		for x, space := range row {
			if space.organism != nil {
				img.SetColorIndex(x, y, 26)
			} else if space.scent > 0.01 {
				scentColorIndex := uint8(space.scent * 100)
				if scentColorIndex > 26 {
					scentColorIndex = 26
				}
				img.SetColorIndex(x, y, scentColorIndex)
			}
		}
	}

	return img
}

func setNextPositions(grid Grid) {
	for _, row := range grid.rows {
		for _, space := range row {
			if space.organism != nil {
				organism := space.organism

				vel := Velocity{
					direction: organism.direction,
					speed:     1,
				}
				nextX, nextY := NextPos(organism.xPos, organism.yPos, vel)
				nextDiscreteX, nextDiscreteY := FloatPosToGridPos(nextX, nextY)

				organism.xPos = nextX
				organism.yPos = nextY

				if grid.hasCoord(nextDiscreteX, nextDiscreteY) {
					organism.nextDiscreteXPos = nextDiscreteX
					organism.nextDiscreteYPos = nextDiscreteY
				} else {
					space.organism = nil
				}
			}
		}
	}
}

func moveOrganisms(grid Grid) {
	for _, row := range grid.rows {
		for _, space := range row {
			space.scent *= scentDecay
			// TODO: scent spread

			if space.organism != nil {
				organism := space.organism
				destinationSpace := grid.rows[organism.nextDiscreteXPos][organism.nextDiscreteYPos]

				// Movement Step
				// Move if possible, otherwise change directions
				moveOrganism(organism, space, destinationSpace)

				// Sensory Step
				// Find the direction with the largest scent trail and turn in that direction
				rotateOrganism(organism, grid)
			}
		}
	}
}

func moveOrganism(organism *Organism, currentSpace *Space, destinationSpace *Space) {
	if destinationSpace.organism == nil {
		// NOTE: This is a greedy approach
		// Setting this to nil before calculating the rest of the movements
		// means the top-left of the screen becomes more dense than the
		// bottom-right because it opens up movement options to the upper-left
		// zones before the lower right zones
		currentSpace.organism = nil
		destinationSpace.organism = organism
		destinationSpace.scent++
	} else if destinationSpace.organism.id != organism.id {
		organism.direction = float64(rand.Intn(360))
	}
}

func rotateOrganism(organism *Organism, grid Grid) {
	sensorDistance := float64(options["sensorDistance"])
	sensorDegree := float64(options["sensorDegree"])
	leftSensorDirection := organism.direction + sensorDegree
	if leftSensorDirection > 360 {
		leftSensorDirection -= 360
	}
	rightSensorDirection := organism.direction - sensorDegree
	if rightSensorDirection < 0 {
		rightSensorDirection += 360
	}
	leftSensorVelocity := Velocity{
		speed:     sensorDistance,
		direction: leftSensorDirection,
	}
	rightSensorVelocity := Velocity{
		speed:     sensorDistance,
		direction: rightSensorDirection,
	}
	frontSensorVelocity := Velocity{
		speed:     sensorDistance,
		direction: organism.direction,
	}
	leftSensorX, leftSensorY := NextPos(organism.xPos, organism.yPos, leftSensorVelocity)
	rightSensorX, rightSensorY := NextPos(organism.xPos, organism.yPos, rightSensorVelocity)
	frontSensorX, frontSensorY := NextPos(organism.xPos, organism.yPos, frontSensorVelocity)
	leftSensorDiscreteX, leftSensorDiscreteY := FloatPosToGridPos(leftSensorX, leftSensorY)
	rightSensorDiscreteX, rightSensorDiscreteY := FloatPosToGridPos(rightSensorX, rightSensorY)
	frontSensorDiscreteX, frontSensorDiscreteY := FloatPosToGridPos(frontSensorX, frontSensorY)

	var leftSensorSpace *Space
	var rightSensorSpace *Space
	var frontSensorSpace *Space
	if grid.hasCoord(leftSensorDiscreteX, leftSensorDiscreteY) {
		leftSensorSpace = grid.rows[leftSensorDiscreteY][leftSensorDiscreteX]
	}
	if grid.hasCoord(rightSensorDiscreteX, rightSensorDiscreteY) {
		rightSensorSpace = grid.rows[rightSensorDiscreteY][rightSensorDiscreteX]
	}
	if grid.hasCoord(frontSensorDiscreteX, frontSensorDiscreteY) {
		frontSensorSpace = grid.rows[frontSensorDiscreteY][frontSensorDiscreteX]
	}

	if leftSensorSpace == nil || rightSensorSpace == nil || frontSensorSpace == nil {
		organism.direction = float64(rand.Intn(360))
	} else {
		leftScent := leftSensorSpace.scent
		rightScent := rightSensorSpace.scent
		frontScent := frontSensorSpace.scent

		if frontScent > leftScent && frontScent > rightScent {
			// Continue in same direction
		} else if leftScent > frontScent && rightScent > frontScent {
			// Rotate randomly left or right
			toss := rand.Intn(2)
			if toss == 0 {
				organism.direction += sensorDegree
				if organism.direction > 360 {
					organism.direction -= 360
				}
			} else {
				organism.direction -= sensorDegree
				if organism.direction < 0 {
					organism.direction += 360
				}
			}
		} else if leftScent > rightScent {
			// Rotate left
			organism.direction += sensorDegree
			if organism.direction > 360 {
				organism.direction -= 360
			}
		} else if rightScent > leftScent {
			// Rotate right
			organism.direction -= sensorDegree
			if organism.direction < 0 {
				organism.direction += 360
			}
		}
		// Else continue in same direction
	}
}

func drawImage(grid Grid, anim *gif.GIF, pal color.Palette) {
	img := createImage(grid, pal)
	anim.Delay = append(anim.Delay, options["delay"])
	anim.Image = append(anim.Image, img)
}

func drawNextFrame(grid Grid, anim *gif.GIF, pal color.Palette) {
	setNextPositions(grid)
	moveOrganisms(grid)
	drawImage(grid, anim, pal)
}
