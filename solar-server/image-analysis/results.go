package image_analysis

import (
	"github.com/dfjones/solar/solar-server/lib/cappedlist"
	"image/color"
)

type AnalyzedImage struct {
	path         string
	averageColor color.RGBA
}

var AnalysisCache *cappedlist.CappedList = cappedlist.New(288)
