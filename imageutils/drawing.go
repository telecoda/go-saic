package imageutils

import (
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/disintegration/imaging"

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

func drawPhotoTiles(sourceImage image.Image, imageTiles *[][]models.ImageTile, tileWidth int) image.Image {

	// convert sourceImage to RGBA image
	bounds := sourceImage.Bounds()
	photoImage := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(photoImage, photoImage.Bounds(), sourceImage, bounds.Min, draw.Src)

	for _, tiles := range *imageTiles {
		for _, tile := range tiles {

			// draw image using first tile discovered
			if tile.SimilarImages != nil {
				//
				for _, similarImage := range *tile.SimilarImages {
					tileImage, _, err := LoadImage(similarImage.FilePath)
					if err != nil {
						panic("Error loading image")
					}
					// resize image to tile size
					resizedImage := ResizeImage(tileImage, tileWidth, tileWidth)
					draw.Draw(photoImage, tile.Rect, resizedImage, tileImage.Bounds().Min, draw.Src)
					break
				}

			}

		}
	}

	return photoImage
}

func drawTintedPhotoTiles(sourceImage image.Image, imageTiles *[][]models.ImageTile, tileWidth int) image.Image {

	// convert sourceImage to RGBA image
	bounds := sourceImage.Bounds()
	photoImage := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	//draw.Draw(photoImage, photoImage.Bounds(), sourceImage, bounds.Min, draw.Src)

	for _, tiles := range *imageTiles {
		for _, tile := range tiles {

			// draw image using thumbnail
			if tile.ThumbnailImage != nil {
				//
				tileImage, _, err := LoadImage(tile.ThumbnailImage.FilePath)
				if err != nil {
					panic("Error loading image")
				}
				// resize image to tile size
				resizedImage := ResizeImage(tileImage, tileWidth, tileWidth)

				// create grayscale version
				grayscaleImage := imaging.Grayscale(resizedImage)
				draw.Draw(photoImage, tile.Rect, grayscaleImage, tileImage.Bounds().Min, draw.Src)

				tintedImage := image.NewRGBA(image.Rect(0, 0, tile.Rect.Dx(), tile.Rect.Dy()))
				draw.Draw(tintedImage, image.Rect(0, 0, tile.Rect.Dx(), tile.Rect.Dy()), &image.Uniform{tile.ProminentColour}, image.ZP, draw.Src)
				// 50% alpha
				//tile.ProminentColour.A = 128
				// draw a tinted square over the top to match the original colour
				photoImage = imaging.Overlay(photoImage, tintedImage, tile.Rect.Min, 0.5)

			}

		}
	}

	return photoImage
}
