package main

import (
	"flag"
	"fmt"
	"github.com/telecoda/go-saic/imageutils"
	"github.com/telecoda/go-saic/models"
	"image"
	"os"
)

var searchImages bool
var recursiveSearch bool
var createThumbnails bool
var sourceDir string
var thumbnailsDir string
var mosaicImagePath string
var mosaicImage image.Image
var targetImagePath string
var targetWidth int
var targetHeight int
var tileWidth int
var tileHeight int
var horiontalTiles int
var verticalTiles int
var sourceImages []models.SourceImage

func init() {
	defaultSourceDir := "data" + string(os.PathSeparator) + "input" + string(os.PathSeparator) + "sourceimages"
	flag.StringVar(&sourceDir, "source_dir", defaultSourceDir, "directory for source images")
	defaultThumbsDir := "data" + string(os.PathSeparator) + "output" + string(os.PathSeparator) + "thumbnail_images"
	flag.StringVar(&thumbnailsDir, "thumb_dir", defaultThumbsDir, "directory to produce thumbnails images in")
	flag.BoolVar(&searchImages, "s", false, "search for images")
	flag.BoolVar(&recursiveSearch, "r", false, "search image directories recursively")
	flag.BoolVar(&createThumbnails, "t", false, "create thumbnails")
	flag.StringVar(&mosaicImagePath, "mosaic_image_path", "image.png", "path of image to create a mosaic from")
	flag.IntVar(&targetWidth, "target_width", 1024, "default width of image to produce, height will be calculated to maintain aspect ratio")
	flag.StringVar(&targetImagePath, "target_image_path", "target.png", "path of mosaic image to be created")
	flag.IntVar(&tileWidth, "tile_width", 32, "default width of mosaic tile")
	flag.IntVar(&tileHeight, "tile_height", 32, "default height of mosaic tile")

}

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

/*
	The image mosaic creation process involved 3 separate steps.  Not all steps are necessary everytime
	the process is invoked.

	Step one: Source image discovery
	================================
	Prerequisites:- needs a directory containing images.

	This step must be performed at least once.  This is used to discover and catalog the images that will be
	available as a reference that can be used to create a mosaic from.

	Step two: Source image thumbnail creation
	=========================================
	Prerequisites:- needs "discovery" to have run to produce a list of images to process.

	This step must be performed at least once.  It must always be run after step one (discovery).
	The process creates smaller scaled thumbnails of the source images in a separate working directory.

	Steps one and two can be run in isolation if this is a long running task.

	Step three: Creation of a photo mosaic
	======================================
	Prerequisites:- needs "discovery" & "thumbnail" to have run.
				    need "mosaic_image" - this is the image that will be used a the basis of the photo mosaic

	This step will create a photo mosaic using a source image.

	The process will not update the source "mosaic_image" a new "target_image" will be created.

	Summary of the photo mosaic process:
	- target image is created as a copy of the source mosaic_image (this can be scaled to a different size)
	- divide target image into a number of "tiles" based upon the tile height and width parameters
	- each tile is analysed to find its prominent colour
	- each tile is replaced with a thumbnail image of a similar colour
	- repeat for all the tiles on the image
	- probably have lots of gaps in resulting image due to lack of photos
	- think of a crafty way of filling the gaps...

*/
func main() {
	flag.Parse()

	// initialise request

	// source image discovery
	// source image transformation

	// target image creation

	if searchImages {
		fmt.Println("search recursively:", recursiveSearch)
	}

	fmt.Println("mosaic_image_path:", mosaicImagePath)

	mosaicImage, _, err := imageutils.LoadImage(mosaicImagePath)
	if err != nil {
		fmt.Printf("Error trying to load image:%s Error:%s", mosaicImagePath, err)
		return
	}

	targetHeight = calcRelativeImageHeight(mosaicImage.Bounds().Max.X, mosaicImage.Bounds().Max.Y, targetWidth)

	if searchImages {
		fmt.Println("source_dir:", sourceDir)
		sourceImages = imageutils.FindSourceImages(sourceDir)
	}

	if createThumbnails {
		// create a thumbnail for each image
		imageutils.CreateThumbnailImages(sourceImages, thumbnailsDir)
	}

	// create target image
	resizedImage := imageutils.ResizeImage(mosaicImage, uint(targetWidth))

	// draw tiles
	tiledImage := imageutils.DrawColouredTiles(resizedImage, tileWidth, tileHeight)

	// draw photo tiles
	photoImage := imageutils.DrawPhotoTiles(tiledImage, tileWidth, tileHeight)

	// draw a grid where mosaic tiles should be
	gridImage := imageutils.DrawGrid(photoImage, tileWidth, tileHeight)
	// save image created
	err = imageutils.SaveImage(targetImagePath, &gridImage)
	if err != nil {
		fmt.Printf("Error saving new image:%s Error:%s", targetImagePath, err)
		return
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
