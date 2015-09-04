package image_analysis

import (
	"math"
	"testing"

	libcolor "github.com/dfjones/solar/solar-server/lib/color"
)

func TestAvgColor(t *testing.T) {
	analyze("../test-images/1.jpeg")
	compareAvgTo(t, libcolor.HSL{0.422, 0.132, 0.459})
	analyze("../test-images/2.jpeg")
	compareAvgTo(t, libcolor.HSL{0, 0.021, 0.601})
	analyze("../test-images/3.jpeg")
	compareAvgTo(t, libcolor.HSL{0.063, 0.554, 0.566})
	analyze("../test-images/4.jpeg")
	compareAvgTo(t, libcolor.HSL{0.075, 0.323, 0.382})
}

func compareAvgTo(t *testing.T, expected libcolor.HSL) {
	ai := AnalysisCache.Last()
	if ai, ok := ai.(*AnalyzedImage); ok {
		if !closeEnough(ai.AverageColor, expected) {
			t.Error("Expected color: ", expected, "Got: ", ai.AverageColor)
		}
	} else {
		t.Error("Could not convert result to color.RGBA!")
	}
}

func closeEnough(a, b libcolor.HSL) bool {
	maxDelta := 0.001
	diffs := []float64{
		math.Abs(a.H - b.H),
		math.Abs(a.S - b.S),
		math.Abs(a.L - b.L),
	}
	for i := 0; i < 3; i++ {
		if diffs[i] > maxDelta {
			return false
		}
	}
	return true
}

func BenchmarkAnalyze(b *testing.B) {
	for i := 0; i < b.N; i++ {
		analyze("../test-images/1.jpeg")
	}
}
