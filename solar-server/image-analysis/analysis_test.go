package image_analysis

import (
	"testing"

	libcolor "github.com/dfjones/solar/solar-server/lib/color"
)

func TestAvgColor(t *testing.T) {
	analyze("../test-images/1.jpeg")
	compareAvgTo(t, libcolor.HSL{100, 113, 90})
	analyze("../test-images/2.jpeg")
	compareAvgTo(t, libcolor.HSL{138, 141, 142})
}

func compareAvgTo(t *testing.T, expected libcolor.HSL) {
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
