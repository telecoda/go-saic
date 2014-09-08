package imageutils

import (
	"image"
	"image/color"
	"log"

	"github.com/telecoda/go-saic/db"
	"github.com/telecoda/go-saic/models"
)

func CreateImageMosaic(inputImagePath string, outputImagePath string, outputImageWidth int, tileSize int) error {

	log.Println("input_image_path:", inputImagePath)

	inputImage, _, err := LoadImage(inputImagePath)
	if err != nil {
		return err
	}

	outputImageHeight := calcRelativeImageHeight(inputImage.Bounds().Max.X, inputImage.Bounds().Max.Y, outputImageWidth)

	// create output image
	resizedImage := ResizeImage(inputImage, outputImageWidth, outputImageHeight)

	// draw tiles
	//tiledImage := DrawColouredTiles(resizedImage, tileSize, tileSize)

	// how many tiles?
	imageTiles := initImageTiles(outputImageWidth, outputImageHeight, tileSize)

	//log.Printf("ImageTiles:%v", imageTiles)

	// analyse input image colours
	analysedTiles := analyseImageTileColours(resizedImage, imageTiles)

	// update tiles with details of similar images
	preparedTiles := updateSimilarColourImages(analysedTiles)

	// draw colour tiles
	//colouredImage := drawColouredTiles(resizedImage, &preparedTiles)

	// draw photo tiles
	photoImage := drawPhotoTiles(resizedImage, &preparedTiles, tileSize)

	// draw a grid where mosaic tiles should be
	gridImage := DrawGrid(photoImage, tileSize, tileSize)
	// save image created
	err = SaveImage(outputImagePath, &gridImage)
	return err
}

func calcRelativeImageHeight(originalWidth int, originalHeight int, targetWidth int) int {
	floatWidth := float64(originalWidth)
	floatHeight := float64(originalHeight)

	aspectRatio := float64(targetWidth) / floatWidth

	adjustedHeight := floatHeight * aspectRatio

	targetHeight := int(adjustedHeight)
	log.Printf("Source width:%d height:%d Target width:%d height:%d\n", originalWidth, originalHeight, targetWidth, targetHeight)
	return targetHeight
}

func calcMosaicTiles(targetWidth int, targetHeight int, tileSize int) (int, int) {

	horzTiles := targetWidth / tileSize
	if targetWidth%tileSize > 0 {
		horzTiles++
	}
	vertTiles := targetHeight / tileSize
	if targetHeight%tileSize > 0 {
		vertTiles++
	}
	log.Printf("Target width:%d height:%d Tile width:%d height:%d Horizontal tiles:%d Vertical tiles:%d\n", targetWidth, targetHeight, tileSize, tileSize, horzTiles, vertTiles)
	return horzTiles, vertTiles
}

func initImageTiles(targetWidth int, targetHeight int, tileSize int) [][]models.ImageTile {

	horzTiles, vertTiles := calcMosaicTiles(targetWidth, targetHeight, tileSize)
	log.Printf("Tiles horizontal:%d vertical:%d", horzTiles, vertTiles)
	// create a 2d array of imageTiles
	imageTiles := make([][]models.ImageTile, vertTiles)
	// Loop over the rows, allocating the slice for each row.
	for i := range imageTiles {
		imageTiles[i] = make([]models.ImageTile, horzTiles)
	}

	// populate tiles with correct co-ordinates
	for x := 0; x < horzTiles; x++ {
		for y := 0; y < vertTiles; y++ {
			currentTile := &imageTiles[x][y]
			currentTile.X = x
			currentTile.Y = y
			tileStartX := x * tileSize
			tileStartY := y * tileSize
			tileEndX := tileStartX + tileSize
			tileEndY := tileStartY + tileSize
			// crop partial tile
			if tileEndX >= targetWidth {
				tileEndX = targetWidth
			}
			// crop partial tile
			if tileEndY >= targetHeight {
				tileEndY = targetHeight
			}
			currentTile.Rect = image.Rectangle{
				image.Point{tileStartX, tileStartY},
				image.Point{tileEndX, tileEndY},
			}
		}
	}

	return imageTiles
}

func analyseImageTileColours(sourceImage image.Image, imageTiles [][]models.ImageTile) [][]models.ImageTile {

	for _, tiles := range imageTiles {
		for _, tile := range tiles {
			tile.ProminentColour = color.RGBA(FindAverageColourInTile(sourceImage, tile.Rect))

			imageTiles[tile.X][tile.Y].ProminentColour = tile.ProminentColour
		}
	}

	return imageTiles
}

func updateSimilarColourImages(imageTiles [][]models.ImageTile) [][]models.ImageTile {

	for _, tiles := range imageTiles {
		for _, tile := range tiles {

			imageTiles[tile.X][tile.Y].SimilarImages = findSimilarColourImages(tile.ProminentColour)

		}
	}

	return imageTiles
}

func findSimilarColourImages(colourToMatch color.RGBA) *[]models.ImageDetail {

	var similarImages *[]models.ImageDetail
	accuracy := 0
	// search for images with similar amount of Red
	for similarImages == nil {
		accuracy += 10
		if accuracy > 255 {
			return nil
		}

		similarImages = db.FindSimilarColourImages(int(colourToMatch.R), int(colourToMatch.G), int(colourToMatch.B), accuracy)
	}

	return similarImages
}
