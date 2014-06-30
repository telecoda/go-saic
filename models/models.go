package models

import "image/color"

type SourceImage struct {
	FilePath        string // complete path
	Filename        string // just file
	Format          string
	Size            int64
	Width           int
	Height          int
	ThumbnailPath   string
	ProminentColour color.Color
}
