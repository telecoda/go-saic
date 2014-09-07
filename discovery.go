package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"code.google.com/p/go-uuid/uuid"

	"github.com/telecoda/go-saic/db"
	"github.com/telecoda/go-saic/models"
)

func DiscoverImages(sourceImagesPath string) error {

	log.Println("Starting image discovery")
	log.Printf("Searching in directory:%v", sourceImagesPath)
	fmt.Printf("[")

	myWalkFunc := func(path string, fileInfo os.FileInfo, err error) error {

		// filter out image files only
		if fileInfo != nil {
			filename := strings.ToLower(fileInfo.Name())
			if strings.HasSuffix(filename, ".jpg") ||
				strings.HasSuffix(filename, ".png") ||
				strings.HasSuffix(filename, ".gif") {
				fmt.Printf(".")
				sourceImage := new(models.ImageDetail)
				sourceImage.Id = uuid.New()
				sourceImage.FilePath = path
				sourceImage.Filename = fileInfo.Name()
				sourceImage.Size = fileInfo.Size()

				// save to db as we go...
				db.DiscoveredImagesColl.SaveImage(*sourceImage)

			}

		}
		return nil
	}

	filepath.Walk(sourceImagesPath, myWalkFunc)

	fmt.Println("]")
	log.Println("Ending image discovery.")

	return nil
}
