package imageutils

import (
	"image"
	"image/color"
)

type Analyzer interface {
	// function will analyze the pixels in the targetRect to find the most prominent colour
	FindColourInTile(sourceImage image.RGBA, targetRect image.Rectangle) color.RGBA
}
