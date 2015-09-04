package image_analysis

import (
	"image"
	"image/jpeg"
	"log"
	"os"
	"runtime"
	"time"

	libcolor "github.com/dfjones/solar/solar-server/lib/color"
)

type threadResult struct {
	hbins       []int
	saturations []float64
	lightness   []float64
}

func newThreadResult() threadResult {
	return threadResult{
		hbins:       make([]int, int(hMax)),
		saturations: make([]float64, int(hMax)),
		lightness:   make([]float64, int(hMax)),
	}
}

var analyzeChan = make(chan string, 300)

var colorMax = float64(0xFFFF)
var eightMax = float64(0xFF)
var hMax = float64(360)

func init() {
	go analyzer()
}

// Analyze the given file
func Analyze(fileName string) {
	analyzeChan <- fileName
}

func diffMs(start, end time.Time) int64 {
	return (end.UnixNano() - start.UnixNano()) / int64(1e6)
}

func analyzer() {
	for f := range analyzeChan {
		start := time.Now()
		analyze(f)
		log.Println("Analysis time: ", diffMs(start, time.Now()), "ms")
	}
}

func analyze(fileName string) {
	decode := time.Now()
	log.Println("Decode start: ", fileName)
	file, err := os.Open(fileName)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	img, err := jpeg.Decode(file)
	log.Println("Decode finish in:", diffMs(decode, time.Now()), " ms")
	if err != nil {
		log.Println("Error decoding file:", err)
		return
	}
	calcAvg := time.Now()
	avg := AvgColor(img)
	log.Println("Avg Color:", avg, diffMs(calcAvg, time.Now()), "ms")
	AnalysisCache.Add(&AnalyzedImage{
		fileName,
		avg,
	})
}

// AvgColor computes the average HSL values for a given image
func AvgColor(img image.Image) libcolor.HSL {
	bounds := img.Bounds()
	min := bounds.Min
	max := bounds.Max
	cores := runtime.NumCPU()

	resultChan := make(chan threadResult)
	for i := 0; i < cores; i++ {
		go func(i int) {
			tr := newThreadResult()
			for y := min.Y + i; y < max.Y; y += cores {
				for x := min.X; x < max.X; x++ {
					p := img.At(x, y)
					pr, pg, pb, _ := p.RGBA()
					h, s, l := libcolor.RGBToHSL(uint8(pr), uint8(pg), uint8(pb))
					hbinIndex := round(h * hMax)
					tr.hbins[hbinIndex]++
					tr.saturations[hbinIndex] += s
					tr.lightness[hbinIndex] += l
				}
			}
			resultChan <- tr
		}(i)
	}

	merged := newThreadResult()
	for i := 0; i < cores; i++ {
		tr := <-resultChan
		for h := 0; h <= int(hMax); h++ {
			merged.hbins[h] += tr.hbins[h]
			merged.saturations[h] += tr.saturations[h]
			merged.lightness[h] += tr.lightness[h]
		}
	}

	maxCount := merged.hbins[0]
	hue := 0
	for i := 1; i < int(hMax); i++ {
		if merged.hbins[i] > maxCount {
			maxCount = merged.hbins[i]
			hue = i
		}
	}
	return libcolor.HSL{
		H: float64(hue),
		S: merged.saturations[hue] / float64(maxCount),
		L: merged.lightness[hue] / float64(maxCount),
	}
}

func sum(a []uint64) uint64 {
	var s uint64
	for i := 0; i < len(a); i++ {
		s += a[i]
	}
	return s
}

func round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}

func cVal(p uint64) uint8 {
	return uint8((float64(p) / float64(colorMax)) * eightMax)
}
