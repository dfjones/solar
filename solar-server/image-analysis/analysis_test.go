package image_analysis

import (
	"image/color"
	"testing"
)

func TestAvgColor(t *testing.T) {
	analyze("../test-images/1.jpeg")
	compareAvgTo(t, color.RGBA{100, 113, 90, 255})
	analyze("../test-images/2.jpeg")
	compareAvgTo(t, color.RGBA{138, 141, 142, 255})
}

func compareAvgTo(t *testing.T, expected color.RGBA) {
	ai := AnalysisCache.Last()
	if ai, ok := ai.(*AnalyzedImage); ok {
		if ai.AverageColor != expected {
			t.Error("Expected color: ", expected, "Got: ", ai.AverageColor)
		}
	} else {
		t.Error("Could not convert result to color.RGBA!")
	}
}

func BenchmarkAnalyze(b *testing.B) {
	for i := 0; i < b.N; i++ {
		analyze("../test-images/1.jpeg")
	}
}
