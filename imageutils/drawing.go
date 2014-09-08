package imageutils

import (
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/telecoda/go-saic/models"
)

// Create a new image with grid lines drawn over it
func DrawGrid(sourceImage image.Image, tileWidth int, tileHeight int) image.Image {

	log.Println("Drawing grid over image.")

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

func drawColouredTiles(sourceImage image.Image, imageTiles *[][]models.ImageTile) image.Image {

	// convert sourceImage to RGBA image
	bounds := sourceImage.Bounds()
	colouredImage := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(colouredImage, colouredImage.Bounds(), sourceImage, bounds.Min, draw.Src)

	for _, tiles := range *imageTiles {
		for _, tile := range tiles {

			draw.Draw(colouredImage, tile.Rect, &image.Uniform{tile.ProminentColour}, image.ZP, draw.Src)

		}
	}

	return colouredImage
}

func drawPhotoTiles(sourceImage image.Image, imageTiles *[][]models.ImageTile) image.Image {

	// convert sourceImage to RGBA image
	bounds := sourceImage.Bounds()
	photoImage := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(photoImage, photoImage.Bounds(), sourceImage, bounds.Min, draw.Src)

	for _, tiles := range *imageTiles {
		for _, tile := range tiles {

			// find an image of a similar colour
			draw.Draw(photoImage, tile.Rect, &image.Uniform{tile.ProminentColour}, image.ZP, draw.Src)

		}
	}

	return photoImage
}

// Create a new image tiles consisting of photos of a similar colour
/*func DrawPhotoTiles(sourceImage image.Image, tileWidth int, tileHeight int) image.Image {

	log.Println("Drawing photo tiles.")

	// convert sourceImage to RGBA image
	bounds := sourceImage.Bounds()
	photoImage := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(photoImage, photoImage.Bounds(), sourceImage, bounds.Min, draw.Src)

	// draw tiles
	for x := 0; x < bounds.Dx(); x += tileWidth {
		for y := 0; y < bounds.Dy(); y += tileHeight {

			targetRect := image.Rectangle{
				image.Point{x, y},
				image.Point{x + tileWidth, y + tileHeight},
			}
			tileColour := FindColourInTile(*photoImage, targetRect)
			log.Printf("Tile rect:%v colour:%v\n", targetRect, tileColour)

			tileBounds := image.Rect(x, y, x+tileWidth-1, y+tileHeight-1)

			draw.Draw(photoImage, tileBounds, &image.Uniform{tileColour}, image.ZP, draw.Src)

		}
	}

	return photoImage
}*/
