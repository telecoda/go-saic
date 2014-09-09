package main

/*
   go-saic - photomosaic generator
   by @telecoda
*/

import (
	"flag"
	"fmt"
	"image"
	"os"

	"github.com/telecoda/go-saic/db"
	"github.com/telecoda/go-saic/imageutils"
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
var optMosaicType string
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
	flag.StringVar(&optMosaicType, "type", "tinted", "Type of mosaic (tinted or matched)")
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

	if optMosaicType != "matched" && optMosaicType != "tinted" {
		fmt.Println("Error: -type parameter must be 'tinted' or 'matched'")
		return
	}

	db.InitDB(optClearDB, optScrubDB)

	if optListDB {
		db.ListDB()
	}

	// initialise request

	if optDiscoverImages {

		err := DiscoverImages(sourceDir)
		if err != nil {
			fmt.Printf("Error trying to discover images. Error:%s\n", err)
			return
		}

	}

	if optCreateThumbnails {

		err := imageutils.CreateThumbnails(thumbnailsDir)

		if err != nil {
			fmt.Printf("Error creating image thumbnails. Error:%s\n", err)
			return
		}

	}

	if optCreateMosaic {

		err := imageutils.CreateImageMosaic(inputImagePath, outputImagePath, outputImageWidth, tileSize, optMosaicType)
		if err != nil {
			fmt.Printf("Error creating image mosaic:%s Error:%s", outputImagePath, err)
			return
		}

	}

}
