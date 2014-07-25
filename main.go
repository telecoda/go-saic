package main

/*
   go-saic - photomosaic creater
   by @telecoda
*/

import (
	"flag"
	"fmt"
	"github.com/telecoda/go-saic/imageutils"
	"github.com/telecoda/go-saic/models"
	"image"
	"os"
)

// command parameters
var optDiscoverImages bool
var recursiveSearch bool
var sourceDir string

var optCreateThumbnails bool
var thumbnailsDir string

var optAnalyseColours bool

var optCreateMosaic bool
var inputImagePath string
var inputImage image.Image
var outputImagePath string
var outputImageWidth int
var outputImageHeight int
var tileWidth int
var tileHeight int

// others
var horiontalTiles int
var verticalTiles int

var sourceImages []models.SourceImage

func init() {
	// discover images
	defaultSourceDir := "data" + string(os.PathSeparator) + "input" + string(os.PathSeparator) + "sourceimages"
	flag.BoolVar(&optDiscoverImages, "d", false, "search for images in source_dir")
	flag.BoolVar(&recursiveSearch, "r", false, "search image directories recursively")
	flag.StringVar(&sourceDir, "source_dir", defaultSourceDir, "directory for source images")

	// create thumbnails
	defaultThumbsDir := "data" + string(os.PathSeparator) + "output" + string(os.PathSeparator) + "thumbnail_images"
	flag.BoolVar(&optCreateThumbnails, "t", false, "create thumbnails")
	flag.StringVar(&thumbnailsDir, "thumb_dir", defaultThumbsDir, "directory to create thumbnails images in")

	// analyse colours
	flag.BoolVar(&optAnalyseColours, "c", false, "Analyse thumbnail images for most prominent colour")

	// create mosaic
	flag.BoolVar(&optCreateMosaic, "m", false, "Create a photo mosaic image")
	flag.StringVar(&inputImagePath, "f", "image.png", "path of input image (used to create mosaic from)")
	flag.IntVar(&outputImageWidth, "output_width", 1024, "default width of image to produce, height will be calculated to maintain aspect ratio")
	flag.StringVar(&outputImagePath, "o", "output.png", "path of output image")
	flag.IntVar(&tileWidth, "tile_width", 32, "width of image tiles in output image")
	flag.IntVar(&tileHeight, "tile_height", 32, "height of image tiles in output image")

}

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	// initialise request

	if optDiscoverImages {

		fmt.Println("source_dir:", sourceDir)
		sourceImages = imageutils.FindSourceImages(sourceDir)

	}

	if optCreateThumbnails {

		// create a thumbnail for each image
		imageutils.CreateThumbnailImages(sourceImages, thumbnailsDir)

	}

	if optAnalyseColours {

		// Not implemented yet..
	}

	if optCreateMosaic {

		fmt.Println("mosaic_image_path:", inputImagePath)

		inputImage, _, err := imageutils.LoadImage(inputImagePath)
		if err != nil {
			fmt.Printf("Error trying to load image:%s Error:%s", inputImagePath, err)
			return
		}

		outputImageHeight = calcRelativeImageHeight(inputImage.Bounds().Max.X, inputImage.Bounds().Max.Y, outputImageWidth)

		// create output image
		resizedImage := imageutils.ResizeImage(inputImage, uint(outputImageWidth))

		// draw tiles
		tiledImage := imageutils.DrawColouredTiles(resizedImage, tileWidth, tileHeight)

		// draw photo tiles
		photoImage := imageutils.DrawPhotoTiles(tiledImage, tileWidth, tileHeight)

		// draw a grid where mosaic tiles should be
		gridImage := imageutils.DrawGrid(photoImage, tileWidth, tileHeight)
		// save image created
		err = imageutils.SaveImage(outputImagePath, &gridImage)
		if err != nil {
			fmt.Printf("Error saving new image:%s Error:%s", outputImagePath, err)
			return
		}

	}

}

func calcRelativeImageHeight(originalWidth int, originalHeight int, targetWidth int) int {
	floatWidth := float64(originalWidth)
	floatHeight := float64(originalHeight)

	aspectRatio := float64(targetWidth) / floatWidth

	adjustedHeight := floatHeight * aspectRatio

	targetHeight := int(adjustedHeight)
	fmt.Printf("Source width:%d height:%d Target width:%d height:%d\n", originalWidth, originalHeight, targetWidth, targetHeight)
	return targetHeight
}

func calcMosaicTiles(targetWidth int, targetHeight int, tileWidth int, tileHeight int) (int, int) {

	horzTiles := targetWidth / tileWidth
	if targetWidth%tileWidth > 0 {
		horzTiles++
	}
	vertTiles := targetHeight / tileHeight
	if targetHeight%tileHeight > 0 {
		vertTiles++
	}
	fmt.Printf("Target width:%d height:%d Tile width:%d height:%d Horizontal tiles:%d Vertical tiles:%d\n", targetWidth, targetHeight, tileWidth, tileHeight, horzTiles, vertTiles)
	return horzTiles, vertTiles
}
