package main

/*
   go-saic - photomosaic creater
   by @telecoda
*/

import (
	"flag"
	"fmt"
	"github.com/telecoda/go-saic/imageutils"
	"image"
	"os"
)

// command parameters
var optClearDB bool
var optListDB bool
var optScrubDB bool
var optDiscoverImages bool
var recursiveSearch bool
var sourceDir string

var optCreateThumbnails bool
var thumbnailsDir string

var optCreateMosaic bool
var inputImagePath string
var inputImage image.Image
var outputImagePath string
var outputImageWidth int
var outputImageHeight int
var tileSize int

// others
var horiontalTiles int
var verticalTiles int

func init() {
	// db
	flag.BoolVar(&optClearDB, "X", false, "clear image db")
	flag.BoolVar(&optListDB, "l", false, "list image db content")
	flag.BoolVar(&optScrubDB, "R", false, "repair & compact db")

	// discover images
	defaultSourceDir := "data" + string(os.PathSeparator) + "input" + string(os.PathSeparator) + "sourceimages"
	flag.BoolVar(&optDiscoverImages, "d", false, "search for images in source_dir")
	flag.BoolVar(&recursiveSearch, "r", false, "search image directories recursively")
	flag.StringVar(&sourceDir, "source_dir", defaultSourceDir, "directory for source images")

	// create thumbnails
	defaultThumbsDir := "data" + string(os.PathSeparator) + "output" + string(os.PathSeparator) + "thumbnail_images"
	flag.BoolVar(&optCreateThumbnails, "t", false, "create thumbnails")
	flag.StringVar(&thumbnailsDir, "thumb_dir", defaultThumbsDir, "directory to create thumbnails images in")

	// create mosaic
	flag.BoolVar(&optCreateMosaic, "m", false, "Create a photo mosaic image")
	flag.StringVar(&inputImagePath, "f", "image.png", "path of input image (used to create mosaic from)")
	flag.IntVar(&outputImageWidth, "output_width", 1024, "default width of image to produce, height will be calculated to maintain aspect ratio")
	flag.StringVar(&outputImagePath, "o", "output.png", "path of output image")
	flag.IntVar(&tileSize, "tile_size", 32, "size of image tiles in output image, width & height are the same")

}

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {

	flag.Parse()

	initDB()

	if optListDB {
		listDB()
	}

	// initialise request

	if optDiscoverImages {

		request := &DiscoveryRequest{
			sourceImagesPath: sourceDir,
		}

		_, err := DiscoverImages(*request)
		if err != nil {
			fmt.Printf("Error trying to discover images. Error:%s\n", err)
			return
		}

	}

	if optCreateThumbnails {

		err := CreateThumbnails(thumbnailsDir)

		if err != nil {
			fmt.Printf("Error creating image thumbnails. Error:%s\n", err)
			return
		}

	}

	if optCreateMosaic {

		fmt.Println("input_image_path:", inputImagePath)

		inputImage, _, err := imageutils.LoadImage(inputImagePath)
		if err != nil {
			fmt.Printf("Error trying to load image:%s Error:%s", inputImagePath, err)
			return
		}

		outputImageHeight = calcRelativeImageHeight(inputImage.Bounds().Max.X, inputImage.Bounds().Max.Y, outputImageWidth)

		// create output image
		resizedImage := ResizeImage(inputImage, outputImageWidth, outputImageHeight)

		// draw tiles
		tiledImage := imageutils.DrawColouredTiles(resizedImage, tileSize, tileSize)

		// draw photo tiles
		photoImage := imageutils.DrawPhotoTiles(tiledImage, tileSize, tileSize)

		// draw a grid where mosaic tiles should be
		gridImage := imageutils.DrawGrid(photoImage, tileSize, tileSize)
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

func calcMosaicTiles(targetWidth int, targetHeight int, tileSize int) (int, int) {

	horzTiles := targetWidth / tileSize
	if targetWidth%tileSize > 0 {
		horzTiles++
	}
	vertTiles := targetHeight / tileSize
	if targetHeight%tileSize > 0 {
		vertTiles++
	}
	fmt.Printf("Target width:%d height:%d Tile width:%d height:%d Horizontal tiles:%d Vertical tiles:%d\n", targetWidth, targetHeight, tileSize, tileSize, horzTiles, vertTiles)
	return horzTiles, vertTiles
}
