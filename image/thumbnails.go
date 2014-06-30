package image

import (
	"bufio"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/telecoda/go-saic/models"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

const THUMBNAIL_WIDTH = 40

func CreateThumbnailImages(sourceImages []models.SourceImage, thumbnailImagesDir string) error {
	log.Println("Starting createThumbnailImages.")

	// check if dir does not exist
	if _, err := os.Stat(thumbnailImagesDir); err != nil {
		log.Printf("thumbnail_images dir does not exist:%s\n", thumbnailImagesDir)
		log.Println("Creating missing dir.")

		err := os.Mkdir(thumbnailImagesDir, 0777)
		if err != nil {
			return err
		}
	}

	for _, sourceImage := range sourceImages {

		if sourceImage.FilePath != "" {
			thumbNailImage, err := createThumbnailImage(&sourceImage, thumbnailImagesDir)
			if err != nil {
				log.Printf("Error during createThumbnailImages: %s", err)
				return err
			}
			if thumbNailImage != nil {
				log.Printf("Thumbnail created. Dimensions width:%d height:%d \n", thumbNailImage.Bounds().Max.X, thumbNailImage.Bounds().Max.Y)
			}

		}

	}

	fmt.Println("Ending createThumbnailImages.")

	return nil
}

func createThumbnailImage(sourceImage *models.SourceImage, thumbnailImagesDir string) (image.Image, error) {
	log.Printf("Starting createThumbnailImage:%s\n", sourceImage.Filename)

	filename := strings.ToLower(sourceImage.Filename)

	defer fmt.Printf("Ending createThumbnailImage:%s\n", filename)

	file, err := os.Open(sourceImage.FilePath)
	if err != nil {
		log.Printf("Error during createThumbnailImage: %s", err)
		return nil, err
	}
	defer file.Close()

	log.Printf("File opened: %s", file.Name())
	loadedImage, format, err := image.Decode(file)
	if err != nil {
		log.Printf("Error during createThumbnailImage: %s", err)
		return nil, err
	}

	// update attributes
	sourceImage.Width = loadedImage.Bounds().Max.X
	sourceImage.Height = loadedImage.Bounds().Max.Y

	fmt.Printf("Image loaded name:%s format:%s %d\n", filename, format, loadedImage.Bounds().Max.X)

	// resize
	thumbnailImage := resize.Resize(THUMBNAIL_WIDTH, 0, loadedImage, resize.Lanczos3)

	fullPath := thumbnailImagesDir + string(os.PathSeparator) + sourceImage.Filename
	// remove file extension
	fullPath = strings.TrimSuffix(fullPath, ".png")
	fullPath = strings.TrimSuffix(fullPath, ".jpg")
	fullPath = strings.TrimSuffix(fullPath, ".gif")

	// all saved as .png files
	fullPath += ".png"

	if imgFilePng, err := os.Create(fullPath); err != nil {
		log.Printf("Error creating PNG thumbnail: %s\n", err)
		return nil, err
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

	sourceImage.ThumbnailPath = fullPath

	promColor, err := findProminentColour(thumbnailImage)
	if err != nil {
		log.Printf("Error finding prominent colour: %s", err)
		return nil, err
	}

	sourceImage.ProminentColour = promColor
	log.Printf("Image: %s Prominent colour: %d-%d-%d-%d", sourceImage.Filename, sourceImage.ProminentColour)

	return thumbnailImage, nil

}
