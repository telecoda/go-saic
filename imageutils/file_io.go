package imageutils

import (
	"bufio"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
)

func LoadImage(imagePath string) (image.Image, string, error) {

	file, err := os.Open(imagePath)
	if err != nil {
		log.Printf("Error during LoadImage: %s", err)
		return nil, "", err
	}
	defer file.Close()
	loadedImage, format, err := image.Decode(file)

	return loadedImage, format, err

}

func SaveImage(imagePath string, imageToSave *image.Image) error {
	if imgFilePng, err := os.Create(imagePath); err != nil {
		log.Printf("Error saving PNG image: %s\n", err)
		return err
	} else {
		defer imgFilePng.Close()
		buffer := bufio.NewWriter(imgFilePng)
		err := png.Encode(buffer, *imageToSave)
		if err != nil {
			log.Printf("Error encoding image:%s", err)
			return err
		}
		buffer.Flush()

		return nil
	}
}
