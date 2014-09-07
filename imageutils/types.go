package imageutils

import (
	"image"
	"image/color"

	"github.com/telecoda/go-saic/models"
)

type Analyzer interface {
	// function will analyze the pixels in the targetRect to find the most prominent colour
	FindColourInTile(sourceImage image.RGBA, targetRect image.Rectangle) color.RGBA
}

// thumbnailing
type ThumbnailRequest struct {
	InputImage    models.ImageDetail
	Width         int
	ThumbnailsDir string
}

type ThumbnailResponse struct {
	ThumbnailImage models.ImageDetail
}
