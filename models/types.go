package models

import (
	"image"
	"image/color"
)

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

type ImageTile struct {
	X               int
	Y               int
	Rect            image.Rectangle
	ProminentColour color.RGBA
	SimilarImages   *[]ImageDetail
}
