package image

import (
	"image"
	"image/color"
)

func findProminentColour(myImage image.Image) (color.RGBA, error) {

	var totalRed uint64
	var totalGreen uint64
	var totalBlue uint64
	var totalPixels uint64

	totalRed = 0
	totalGreen = 0
	totalBlue = 0

	var rect = myImage.Bounds()

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			colour := myImage.At(x, y)
			r, g, b, _ := colour.RGBA()
			totalRed += uint64(r)
			totalGreen += uint64(g)
			totalBlue += uint64(b)
		}
	}

	totalPixels = uint64(rect.Max.Y * rect.Max.X)

	averageRed := totalRed / totalPixels
	averageGreen := totalGreen / totalPixels
	averageBlue := totalBlue / totalPixels

	averageColour := color.RGBA{R: uint8(averageRed), G: uint8(averageGreen), B: uint8(averageBlue), A: 255}

	return averageColour, nil
}
