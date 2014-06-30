package main

import (
	"flag"
	"fmt"
	goimage "github.com/telecoda/go-saic/image"
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

	mosaicImage, _, err := goimage.LoadImage(mosaicImagePath)
	if err != nil {
		fmt.Printf("Error trying to load image:%s Error:%s", mosaicImagePath, err)
		return
	}

	targetHeight = calcRelativeImageHeight(mosaicImage.Bounds().Max.X, mosaicImage.Bounds().Max.Y, targetWidth)

	if searchImages {
		sourceImages = goimage.FindSourceImages(sourceDir)
	}

	if createThumbnails {
		// create a thumbnail for each image
		goimage.CreateThumbnailImages(sourceImages, thumbnailsDir)
	}

	// create target image
	resizedImage := goimage.ResizeImage(&mosaicImage, uint(targetWidth))

	// save image created
	err = goimage.SaveImage(targetImagePath, &resizedImage)
	if err != nil {
		fmt.Printf("Error saving new image:%s Error:%s", targetImagePath, err)
		return
	}

}

// calculates
func calcRelativeImageHeight(originalWidth int, originalHeight int, targetWidth int) int {
	floatWidth := float64(originalWidth)
	floatHeight := float64(originalHeight)

	aspectRatio := float64(targetWidth) / floatWidth

	adjustedHeight := floatHeight * aspectRatio

	targetHeight := int(adjustedHeight)
	fmt.Printf("Source width:%d height:%d Target width:%d height:%d\n", originalWidth, originalHeight, targetWidth, targetHeight)
	return targetHeight
}
