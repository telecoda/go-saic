package main

import (
	"image"
	"image/color"
)

// discovery
type DiscoveryRequest struct {
	sourceImagesPath string
}

type DiscoveryResponse struct {
	imagesDiscovered []ImageDetail
}

type ImageDetail struct {
	Id       string // optional Id
	FilePath string // complete path
	Filename string // just file
	Format   string
	Size     int64
	Width    int
	Height   int
	// these are the average values of colours
	// to help find similar imagess
	Red   int
	Blue  int
	Green int
}

// thumbnailing
type ThumbnailRequest struct {
	InputImage    ImageDetail
	Width         int
	ThumbnailsDir string
}

type ThumbnailResponse struct {
	ThumbnailImage ImageDetail
}

// prominent colour analysis
type ProminentColourRequest struct {
	AnalysisImage ImageDetail
	WithinRect    image.Rectangle
}

type ProminentColourResponse struct {
	ProminentColour color.Color
}

// mosaic creation
type MosaicCreateRequest struct {
	AnalysisImage ImageDetail
}

type MosaicCreateResponse struct {
	MosaicImage image.Image
}
