package image

import (
	"fmt"
	"github.com/nfnt/resize"
	"github.com/telecoda/go-saic/models"
	"image"
	"log"
	"os"
	"strings"
)

const THUMBNAIL_WIDTH = 40

func CreateThumbnailImages(sourceImages []models.SourceImage, thumbnailImagesDir string) error {
	log.Println("Starting CreateThumbnailImages.")
	fmt.Printf("[")
	defer log.Println("Ending CreateThumbnailImages.")

	// check if dir does not exist
	if _, err := os.Stat(thumbnailImagesDir); err != nil {
		err := os.Mkdir(thumbnailImagesDir, 0777)
		if err != nil {
			return err
		}
	}

	for _, sourceImage := range sourceImages {

		if sourceImage.FilePath != "" {
			_, err := createThumbnailImage(&sourceImage, thumbnailImagesDir)
			if err != nil {
				log.Printf("Error during createThumbnailImages: %s", err)
				return err
			}
			fmt.Printf(".")

		}

	}
	fmt.Println("]")

	return nil
}

func createThumbnailImage(sourceImage *models.SourceImage, thumbnailImagesDir string) (image.Image, error) {

	loadedImage, format, err := LoadImage(sourceImage.FilePath)
	if err != nil {
		log.Printf("Error during createThumbnailImage: %s", err)
		return nil, err
	}

	// update attributes
	sourceImage.Width = loadedImage.Bounds().Max.X
	sourceImage.Height = loadedImage.Bounds().Max.Y
	sourceImage.Format = format

	// resize
	thumbnailImage := resize.Resize(THUMBNAIL_WIDTH, 0, loadedImage, resize.Lanczos3)

	fullPath := thumbnailImagesDir + string(os.PathSeparator) + sourceImage.Filename
	// remove file extension
	fullPath = strings.TrimSuffix(fullPath, ".png")
	fullPath = strings.TrimSuffix(fullPath, ".jpg")
	fullPath = strings.TrimSuffix(fullPath, ".gif")

	// all saved as .png files
	fullPath += ".png"

	err = SaveImage(fullPath, &thumbnailImage)
	if err != nil {
		log.Printf("Error creating PNG thumbnail: %s\n", err)
		return nil, err
	}

	sourceImage.ThumbnailPath = fullPath

	promColor, err := findProminentColour(thumbnailImage)
	if err != nil {
		log.Printf("Error finding prominent colour: %s", err)
		return nil, err
	}

	sourceImage.ProminentColour = promColor
	return thumbnailImage, nil

}

func ResizeImage(originalImage *image.Image, newWidth int) image.Image {

	return resize.Resize(newWidth, 0, originalImage, resize.Lanczos3)

}
