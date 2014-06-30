package main

import (
	"flag"
	"fmt"
	goimage "github.com/telecoda/go-saic/image"
	"github.com/telecoda/go-saic/models"
	"os"
)

var scanImages bool
var recursiveScan bool
var createThumbnails bool
var imageDir string
var thumbnailsDir string
var mosaicImage string
var sourceImages []models.SourceImage

func init() {
	flag.StringVar(&imageDir, "image_dir", "images", "images directory")
	flag.StringVar(&thumbnailsDir, "thumb_dir", "thumbnail_images", "thumbnials images directory")
	flag.BoolVar(&scanImages, "s", false, "scan images")
	flag.BoolVar(&recursiveScan, "r", false, "scan image directories recursively")
	flag.BoolVar(&createThumbnails, "t", false, "create thumbnails")
	flag.StringVar(&mosaicImage, "image", "image.png", "mosaic image")
}

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	fmt.Println("scanimages:", scanImages)
	if scanImages {
		fmt.Println("scan recursively:", recursiveScan)
	}

	fmt.Println("imagedir:", imageDir)
	fmt.Println("mosaicimage:", mosaicImage)

	if scanImages {
		sourceImages = goimage.FindSourceImages(imageDir)
	}

	if createThumbnails {
		// create a thumbnail for each image
		goimage.CreateThumbnailImages(sourceImages, thumbnailsDir)
	}

}
