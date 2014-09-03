package main

import (
	"encoding/json"
	"fmt"
	//"github.com/nfnt/resize"
	"github.com/disintegration/imaging"
	"github.com/telecoda/go-saic/imageutils"
	"image"
	"log"
	"os"
	"strings"
)

const THUMBNAIL_WIDTH = 40

// create a thumbnail for each image in db
func CreateThumbnails(thumbnailImagesDir string) error {

	log.Println("Starting CreateThumbnails.")
	fmt.Printf("[")
	defer log.Println("Ending CreateThumbnails.")

	// check if dir does not exist
	if _, err := os.Stat(thumbnailImagesDir); err != nil {
		err := os.Mkdir(thumbnailImagesDir, 0777)
		if err != nil {
			return err
		}
	}

	// get a list of discoveredImages from the DB
	discoveredImagesColl.dbCol.ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {

		var inputImage ImageDetail

		if json.Unmarshal(docContent, &inputImage) != nil {
			fmt.Println("cannot deserialize!")
			return false
		}

		// create a thumbnail for this image
		request := ThumbnailRequest{
			InputImage:    inputImage,
			Width:         THUMBNAIL_WIDTH,
			ThumbnailsDir: thumbnailImagesDir,
		}

		outputImage, err := createThumbnailImage(request)
		if err != nil {
			fmt.Printf("Error: problem creating thumbnail - %v \n", err)
			return true
		}

		fmt.Print(".")
		thumbnailImagesColl.saveImage(outputImage.ThumbnailImage)

		return true  // move on to the next document OR
		return false // do not move on to the next document
	})

	fmt.Println("]")
	return nil
}
func createThumbnailImage(request ThumbnailRequest) (*ThumbnailResponse, error) {

	/* Note: for now all thumbnails will be square
	   as it just makes everything much easier...
	*/
	loadedImage, format, err := imageutils.LoadImage(request.InputImage.FilePath)
	if err != nil {
		log.Printf("Error during createThumbnailImage: %s", err)
		return nil, err
	}

	// update attributes
	request.InputImage.Width = loadedImage.Bounds().Max.X
	request.InputImage.Height = loadedImage.Bounds().Max.Y
	request.InputImage.Format = format

	// crop image to centre square
	// take size of smallest dimension
	var squareSize int
	if request.InputImage.Width <= request.InputImage.Height {
		squareSize = request.InputImage.Width
	} else {
		squareSize = request.InputImage.Height
	}

	croppedImage := imaging.CropCenter(loadedImage, squareSize, squareSize)

	// resize
	thumbnailImage := ResizeImage(croppedImage, THUMBNAIL_WIDTH, THUMBNAIL_WIDTH)

	var fullPath string = request.ThumbnailsDir + string(os.PathSeparator) + request.InputImage.Id
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
			Id:       request.InputImage.Id,
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

func ResizeImage(originalImage image.Image, newWidth int, newHeight int) image.Image {

	//return resize.Resize(newWidth, 0, originalImage, resize.Lanczos3)
	return imaging.Resize(originalImage, newWidth, newHeight, imaging.BSpline)

}
