package imageutils

import (
	"fmt"
	"github.com/telecoda/go-saic/models"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const TOTAL_IMAGES = 100

var sourceImages []models.SourceImage

func FindSourceImages(imageDir string) []models.SourceImage {
	log.Println("Starting findSourceImages.")
	log.Printf("Searching in directory:%v", imageDir)
	fmt.Printf("[")

	sourceImages = make([]models.SourceImage, TOTAL_IMAGES)

	filepath.Walk(imageDir, myWalkFunc)

	fmt.Println("]")
	log.Println("Ending findSourceImages.")

	return sourceImages
}

func myWalkFunc(path string, fileInfo os.FileInfo, err error) error {

	// filter out image files only
	if fileInfo != nil {
		filename := strings.ToLower(fileInfo.Name())
		if strings.HasSuffix(filename, ".jpg") ||
			strings.HasSuffix(filename, ".png") ||
			strings.HasSuffix(filename, ".gif") {
			fmt.Printf(".")
			sourceImage := new(models.SourceImage)
			sourceImage.FilePath = path
			sourceImage.Filename = fileInfo.Name()
			sourceImage.Size = fileInfo.Size()
			sourceImages = append(sourceImages, *sourceImage)

		}

	}
	return nil
}
