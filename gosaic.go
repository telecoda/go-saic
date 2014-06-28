package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const TOTAL_IMAGES = 100

const THUMBNAIL_WIDTH = 80

var scanImages bool
var recursiveScan bool
var createThumbnails bool
var imageDir string
var thumbnailsDir string
var sourceImage string
var originalImages []OriginalImage

type OriginalImage struct {
	filePath      string // complete path
	filename      string // just file
	format        string
	size          int64
	width         int
	height        int
	thumbnailPath string
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

func createThumbnailImages(originalImages []OriginalImage, thumbnailImagesDir string) error {
	fmt.Println("Starting createThumbnailImages.")

	// check if dir does not exist
	if _, err := os.Stat(thumbnailImagesDir); err != nil {
		fmt.Printf("thumbnail_images dir does not exist:%s\n", thumbnailImagesDir)
		fmt.Println("Creating missing dir.")

		err := os.Mkdir(thumbnailImagesDir, 0777)
		if err != nil {
			return err
		}
	}

	for _, originalImage := range originalImages {
		//if originalImage = nil {
		if originalImage.filePath != "" {
			thumbNailImage, err := createThumbnailImage(&originalImage)
			if err != nil {
				log.Printf("Error during createThumbnailImages: %s", err)
				return err
			}
			if thumbNailImage != nil {
				fmt.Printf("Thumbnail created. Dimensions width:%d height:%d \n", thumbNailImage.Bounds().Max.X, thumbNailImage.Bounds().Max.Y)
			}

		}

		//}
	}

	fmt.Println("Ending createThumbnailImages.")

	return nil
}

func createThumbnailImage(originalImage *OriginalImage) (image.Image, error) {
	fmt.Printf("Starting createThumbnailImage:%s\n", originalImage.filename)

	filename := strings.ToLower(originalImage.filename)

	defer fmt.Printf("Ending createThumbnailImage:%s\n", filename)

	file, err := os.Open(originalImage.filePath)
	if err != nil {
		log.Printf("Error during createThumbnailImage: %s", err)
		return nil, err
	}
	defer file.Close()

	log.Printf("File opened: %s", file.Name())
	sourceImage, format, err := image.Decode(file)
	if err != nil {
		log.Printf("Error during createThumbnailImage: %s", err)
		return nil, err
	}

	// update attributes
	originalImage.width = sourceImage.Bounds().Max.X
	originalImage.height = sourceImage.Bounds().Max.Y

	fmt.Printf("Image loaded name:%s format:%s %d\n", filename, format, sourceImage.Bounds().Max.X)

	// resize
	rectangle := image.Rect(0, 0, 40, 40)
	point := image.Point{0, 0}
	thumbnailImage := image.NewRGBA(rectangle)
	draw.Draw(thumbnailImage, rectangle, sourceImage, point, 1)
	/*if imgFileJpg, err := os.Create("red-thumb.jpg"); err != nil {
		defer imgFileJpg.Close()
		err := jpeg.Encode(bufio.NewWriter(imgFileJpg), thumbnailImage, &jpeg.Options{jpeg.DefaultQuality})
		if err != nil {
			log.Printf("Error encoding image:%s", err)
			return nil, err
		}
	}*/

	fullPath := thumbnailsDir + string(os.PathSeparator) + originalImage.filename
	// remove file extension
	fullPath = strings.TrimRight(fullPath, ".png")
	fullPath = strings.TrimRight(fullPath, ".jpg")
	fullPath = strings.TrimRight(fullPath, ".gif")

	fullPath += ".png"

	if imgFilePng, err := os.Create(fullPath); err != nil {
		fmt.Println("Png error: ", err)
	} else {
		defer imgFilePng.Close()
		buffer := bufio.NewWriter(imgFilePng)
		err := png.Encode(buffer, thumbnailImage)
		if err != nil {
			log.Printf("Error encoding image:%s", err)
			return nil, err
		}
		buffer.Flush()
	}

	originalImage.thumbnailPath = fullPath

	return thumbnailImage, nil

}
