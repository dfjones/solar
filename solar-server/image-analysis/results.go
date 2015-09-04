package image_analysis

import (
	"github.com/dfjones/solar/solar-server/lib/cappedlist"
	libcolor "github.com/dfjones/solar/solar-server/lib/color"
)

type AnalyzedImage struct {
	Path         string
	AverageColor libcolor.HSL
}

var AnalysisCache *cappedlist.CappedList = cappedlist.New(288)
