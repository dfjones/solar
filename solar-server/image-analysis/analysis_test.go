package image_analysis

import (
	"image/color"
	"testing"
)

func TestAvgColor(t *testing.T) {
	file := "../test-images/1.jpeg"
	analyze(file)
	ai := AnalysisCache.Last()
	if ai, ok := ai.(*AnalyzedImage); ok {
		expected := color.RGBA{100, 113, 90, 255}
		if ai.AverageColor != expected {
			t.Error("Expected color: ", expected, "Got: ", ai.AverageColor)
		}
	} else {
		t.Error("Could not convert result to color.RGBA!")
	}
}
