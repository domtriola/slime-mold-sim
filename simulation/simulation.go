package simulation

import (
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

var options = map[string]int{
	"width":             500,
	"height":            500,
	"nFrames":           500,
	"loopCount":         1000,
	"delay":             2,
	"sensorDegree":      45,
	"sensorDistance":    9,
	"scentSpreadFactor": 3,
}

// Build creates the simulation as a GIF
func Build(urlOptions map[string]interface{}) (name string) {
	setOptions(urlOptions)

	name = "tmp/image.gif"
	animate(name)
	return name
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

	grid.initialize()
	for i := 0; i < options["nFrames"]; i++ {
		drawNextFrame(grid, &anim, pal)
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
				img.SetColorIndex(x, y, 1)
			} else {
				img.SetColorIndex(x, y, 0)
			}
		}
	}

	return img
}

func drawNextFrame(grid Grid, anim *gif.GIF, pal color.Palette) {
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

	for _, row := range grid.rows {
		for _, space := range row {
			if space.organism != nil {
				organism := space.organism

				destinationSpace := grid.rows[organism.nextDiscreteXPos][organism.nextDiscreteYPos]

				if destinationSpace.organism == nil {
					space.organism = nil
					destinationSpace.organism = organism
				} else if destinationSpace.organism.id != organism.id {
					organism.direction = float64(rand.Intn(360))
				}
			}
		}
	}

	img := createImage(grid, pal)
	anim.Delay = append(anim.Delay, options["delay"])
	anim.Image = append(anim.Image, img)
}
