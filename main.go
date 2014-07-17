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
	flag.StringVar(&sourceDir, "source_dir", "images", "directory for source images")
	flag.StringVar(&thumbnailsDir, "thumb_dir", "thumbnail_images", "directory to produce thumbnails images in")
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

func main() {
	flag.Parse()
	fmt.Println("searchimages:", searchImages)
	if searchImages {
		fmt.Println("search recursively:", recursiveSearch)
	}

	fmt.Println("source_dir:", sourceDir)
	fmt.Println("mosaic_image_path:", mosaicImagePath)

	mosaicImage, _, err := imageutils.LoadImage(mosaicImagePath)
	if err != nil {
		fmt.Printf("Error trying to load image:%s Error:%s", mosaicImagePath, err)
		return
	}

	targetHeight = calcRelativeImageHeight(mosaicImage.Bounds().Max.X, mosaicImage.Bounds().Max.Y, targetWidth)

	if searchImages {
		sourceImages = imageutils.FindSourceImages(sourceDir)
	}

	if createThumbnails {
		// create a thumbnail for each image
		imageutils.CreateThumbnailImages(sourceImages, thumbnailsDir)
	}

	// create target image
	resizedImage := imageutils.ResizeImage(mosaicImage, uint(targetWidth))

	// draw a grid where mosaic tiles should be
	gridImage := imageutils.DrawGrid(resizedImage, tileWidth, tileHeight)
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
