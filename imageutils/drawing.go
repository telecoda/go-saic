package imageutils

import (
	"image"
	"image/color"
	"image/draw"
)

// Create a new image with grid lines drawn over it
func DrawGrid(sourceImage image.Image, tileWidth int, tileHeight int) image.Image {

	lineWidth := 1
	// convert sourceImage to RGBA image
	bounds := sourceImage.Bounds()
	gridImage := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(gridImage, gridImage.Bounds(), sourceImage, bounds.Min, draw.Src)

	lineColour := color.RGBA{0, 0, 0, 255}

	// draw horizontal lines
	for y := 0; y < bounds.Dy(); y += tileHeight {

		lineBounds := image.Rect(0, y, bounds.Dx(), y+lineWidth)
		//lineBounds := &image.Rectangle{Min: {X: 0, Y: 0}, Max: {X: 160, Y: 5}}
		draw.Draw(gridImage, lineBounds, &image.Uniform{lineColour}, image.ZP, draw.Src)

	}

	// draw vertical lines
	for x := 0; x < bounds.Dx(); x += tileWidth {

		lineBounds := image.Rect(x, 0, x+lineWidth, bounds.Dy())
		//lineBounds := &image.Rectangle{Min: {X: 0, Y: 0}, Max: {X: 160, Y: 5}}
		draw.Draw(gridImage, lineBounds, &image.Uniform{lineColour}, image.ZP, draw.Src)

	}

	return gridImage
}

// Create a new image with tiles filled in
func DrawTiles(sourceImage image.Image, tileWidth int, tileHeight int) image.Image {

	// convert sourceImage to RGBA image
	bounds := sourceImage.Bounds()
	tiledImage := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(tiledImage, tiledImage.Bounds(), sourceImage, bounds.Min, draw.Src)

	// draw tiles
	// draw vertical lines
	for x := 0; x < bounds.Dx(); x += tileWidth {
		for y := 0; y < bounds.Dy(); y += tileHeight {

			// pick a colour from within current tile
			tileColour := tiledImage.At(x, y)

			tileBounds := image.Rect(x, y, x+tileWidth, y+tileHeight)
			//lineBounds := &image.Rectangle{Min: {X: 0, Y: 0}, Max: {X: 160, Y: 5}}
			draw.Draw(tiledImage, tileBounds, &image.Uniform{tileColour}, image.ZP, draw.Src)

		}
	}

	return tiledImage
}
