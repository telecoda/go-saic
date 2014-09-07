package imageutils

import (
	"image"
	"image/color"
	"log"

	"github.com/disintegration/imaging"
)

func FindProminentColour(myImage image.Image) color.RGBA {

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

	totalPixels = uint64((rect.Max.Y - rect.Min.Y) * (rect.Max.X - rect.Min.X))

	averageRed := totalRed / totalPixels
	averageGreen := totalGreen / totalPixels
	averageBlue := totalBlue / totalPixels

	averageColour := color.RGBA{R: uint8(averageRed), G: uint8(averageGreen), B: uint8(averageBlue), A: 255}

	return averageColour
}

func FindColourInTile(sourceImage image.Image, targetRect image.Rectangle) color.RGBA {

	var totalRed uint64
	var totalGreen uint64
	var totalBlue uint64
	var totalPixels uint64

	totalRed = 0
	totalGreen = 0
	totalBlue = 0

	for y := targetRect.Min.Y; y <= targetRect.Max.Y; y++ {
		for x := targetRect.Min.X; x <= targetRect.Max.X; x++ {
			colour := sourceImage.At(x, y)
			r, g, b, _ := colour.RGBA()
			//log.Printf("x:%d,y:%d colour:%v", x, y, colour)
			totalRed += uint64(r)
			totalGreen += uint64(g)
			totalBlue += uint64(b)
		}
	}

	totalPixels = uint64((targetRect.Max.Y - targetRect.Min.Y + 1) * (targetRect.Max.X - targetRect.Min.X + 1))

	averageRed := totalRed / totalPixels
	averageGreen := totalGreen / totalPixels
	averageBlue := totalBlue / totalPixels

	averageColour := color.RGBA{R: uint8(averageRed), G: uint8(averageGreen), B: uint8(averageBlue), A: 255}
	log.Printf("average r:%d g:%d b:%d", averageRed, averageGreen, averageBlue)
	log.Printf("total r:%d g:%d b:%d", totalRed, totalGreen, totalBlue)
	log.Printf("totalpixels:%d", totalPixels)

	return averageColour

}

func FindAverageColourInTile(sourceImage image.Image, targetRect image.Rectangle) color.NRGBA {

	// calc colour by scaling image to 1x1 pixels

	croppedImage := imaging.Crop(sourceImage, targetRect)

	// scale
	singlePixelImage := imaging.Resize(croppedImage, 1, 1, imaging.BSpline)

	averageColour := singlePixelImage.At(0, 0)

	return averageColour.(color.NRGBA)

}
