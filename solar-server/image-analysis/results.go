package image_analysis

import (
	"github.com/dfjones/solar/solar-server/lib/cappedlist"
	"image/color"
)

type AnalyzedImage struct {
	Path         string
	AverageColor color.RGBA
}

var AnalysisCache *cappedlist.CappedList = cappedlist.New(288)
