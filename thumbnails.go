package main

import (
	"bufio"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

func createThumbnailImages(originalImages []OriginalImage, thumbnailImagesDir string) error {
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

	for _, originalImage := range originalImages {

		if originalImage.filePath != "" {
			thumbNailImage, err := createThumbnailImage(&originalImage)
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

func createThumbnailImage(originalImage *OriginalImage) (image.Image, error) {
	log.Printf("Starting createThumbnailImage:%s\n", originalImage.filename)

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
	thumbnailImage := resize.Resize(THUMBNAIL_WIDTH, 0, sourceImage, resize.Lanczos3)

	fullPath := thumbnailsDir + string(os.PathSeparator) + originalImage.filename
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

	originalImage.thumbnailPath = fullPath

	return thumbnailImage, nil

}
