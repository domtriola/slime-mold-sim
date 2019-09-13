package simulation

import (
	"image"
	"image/color"
	"image/gif"
	"log"
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
	anim := gif.GIF{LoopCount: options["nFrames"]}

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
			if len(space.organisms) > 0 {
				img.SetColorIndex(x, y, 1)
			}
		}
	}

	return img
}

func drawNextFrame(grid Grid, anim *gif.GIF, pal color.Palette) {
	img := createImage(grid, pal)
	anim.Delay = append(anim.Delay, options["delay"])
	anim.Image = append(anim.Image, img)
}
