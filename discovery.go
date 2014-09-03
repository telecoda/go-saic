package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"code.google.com/p/go-uuid/uuid"
)

func DiscoverImages(request DiscoveryRequest) (DiscoveryResponse, error) {

	log.Println("Starting image discovery")
	log.Printf("Searching in directory:%v", request.sourceImagesPath)
	fmt.Printf("[")

	response := new(DiscoveryResponse)

	response.imagesDiscovered = make([]ImageDetail, 0)

	myWalkFunc := func(path string, fileInfo os.FileInfo, err error) error {

		// filter out image files only
		if fileInfo != nil {
			filename := strings.ToLower(fileInfo.Name())
			if strings.HasSuffix(filename, ".jpg") ||
				strings.HasSuffix(filename, ".png") ||
				strings.HasSuffix(filename, ".gif") {
				fmt.Printf(".")
				sourceImage := new(ImageDetail)
				sourceImage.Id = uuid.New()
				sourceImage.FilePath = path
				sourceImage.Filename = fileInfo.Name()
				sourceImage.Size = fileInfo.Size()
				response.imagesDiscovered = append(response.imagesDiscovered, *sourceImage)

				// save to db as we go...
				discoveredImagesColl.saveImage(*sourceImage)

			}

		}
		return nil
	}

	filepath.Walk(request.sourceImagesPath, myWalkFunc)

	fmt.Println("]")
	log.Println("Ending image discovery.")

	return *response, nil
}
