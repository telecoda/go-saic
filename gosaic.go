package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const TOTAL_IMAGES = 100

const THUMBNAIL_WIDTH = 40

var scanImages bool
var recursiveScan bool
var createThumbnails bool
var imageDir string
var thumbnailsDir string
var sourceImage string
var originalImages []OriginalImage

type OriginalImage struct {
	filePath        string // complete path
	filename        string // just file
	format          string
	size            int64
	width           int
	height          int
	thumbnailPath   string
	prominentColour color.Color
}

func init() {
	flag.StringVar(&imageDir, "image_dir", "images", "images directory")
	flag.StringVar(&thumbnailsDir, "thumb_dir", "thumbnail_images", "thumbnials images directory")
	flag.BoolVar(&scanImages, "s", false, "scan images")
	flag.BoolVar(&recursiveScan, "r", false, "scan image directories recursively")
	flag.BoolVar(&createThumbnails, "t", false, "create thumbnails")
	flag.StringVar(&sourceImage, "image", "image.png", "source image")
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
	fmt.Println("sourceimage:", sourceImage)

	if scanImages {
		findOriginalImages(imageDir)
	}

	if createThumbnails {
		// create a thumbnail for each image
		createThumbnailImages(originalImages, thumbnailsDir)
	}

}

func findOriginalImages(imageDir string) {
	log.Println("Starting findOriginalImages.")

	originalImages = make([]OriginalImage, TOTAL_IMAGES)

	filepath.Walk(imageDir, myWalkFunc)

	log.Println("Ending findOriginalImages.")
}

func myWalkFunc(path string, fileInfo os.FileInfo, err error) error {

	// filter out image files only
	filename := strings.ToLower(fileInfo.Name())
	if strings.HasSuffix(filename, ".jpg") ||
		strings.HasSuffix(filename, ".png") ||
		strings.HasSuffix(filename, ".gif") {
		fmt.Printf("Image found. path: %s fileInfo:%s \n", path, fileInfo.Name())

		originalImage := new(OriginalImage)
		originalImage.filePath = path
		originalImage.filename = fileInfo.Name()
		originalImage.size = fileInfo.Size()
		originalImages = append(originalImages, *originalImage)

	}

	return nil
}
