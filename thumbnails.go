package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"github.com/telecoda/go-saic/imageutils"
	"image"
	"log"
	"os"
	"strings"
)

const THUMBNAIL_WIDTH = 40

// create a thumbnail for each image in request
func CreateThumbnails(images []ImageDetail, thumbnailImagesDir string) ([]ImageDetail, error) {

	log.Println("Starting CreateThumbnails.")
	fmt.Printf("[")
	defer log.Println("Ending CreateThumbnails.")

	// check if dir does not exist
	if _, err := os.Stat(thumbnailImagesDir); err != nil {
		err := os.Mkdir(thumbnailImagesDir, 0777)
		if err != nil {
			return nil, err
		}
	}

	// create a list of thumbnail images
	thumbnails := make([]ImageDetail, 0)
	for _, inputImage := range images {

		request := ThumbnailRequest{
			InputImage:    inputImage,
			Width:         THUMBNAIL_WIDTH,
			ThumbnailsDir: thumbnailImagesDir,
		}

		response, err := createThumbnailImage(request)
		if err != nil {
			return nil, err
		}

		thumbnails = append(thumbnails, response.ThumbnailImage)

	}

	return thumbnails, nil

}
func createThumbnailImage(request ThumbnailRequest) (*ThumbnailResponse, error) {

	loadedImage, format, err := imageutils.LoadImage(request.InputImage.FilePath)
	if err != nil {
		log.Printf("Error during createThumbnailImage: %s", err)
		return nil, err
	}

	// update attributes
	request.InputImage.Width = loadedImage.Bounds().Max.X
	request.InputImage.Height = loadedImage.Bounds().Max.Y
	request.InputImage.Format = format

	// resize
	thumbnailImage := resize.Resize(THUMBNAIL_WIDTH, 0, loadedImage, resize.Lanczos3)

	var fullPath string = request.ThumbnailsDir + string(os.PathSeparator) + request.InputImage.Filename
	// remove file extension
	fullPath = strings.TrimSuffix(fullPath, ".png")
	fullPath = strings.TrimSuffix(fullPath, ".jpg")
	fullPath = strings.TrimSuffix(fullPath, ".gif")

	// all saved as .png files
	fullPath += ".png"

	err = imageutils.SaveImage(fullPath, &thumbnailImage)
	if err != nil {
		log.Printf("Error creating PNG thumbnail: %s\n", err)
		return nil, err
	}

	response := &ThumbnailResponse{
		ThumbnailImage: ImageDetail{
			FilePath: fullPath,
			Filename: request.InputImage.Filename,
			Format:   request.InputImage.Format,
			Size:     0,
			Width:    request.InputImage.Width,
			Height:   request.InputImage.Height,
		},
	}

	/*
		promColor, err := findProminentColour(thumbnailImage)
		if err != nil {
			log.Printf("Error finding prominent colour: %s", err)
			return nil, err
		}
	*/
	//	sourceImage.ProminentColour = promColor
	return response, nil

}

func ResizeImage(originalImage image.Image, newWidth uint) image.Image {

	return resize.Resize(newWidth, 0, originalImage, resize.Lanczos3)

}
