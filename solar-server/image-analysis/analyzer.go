package image_analysis

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"runtime"
	"time"
)

var analyzeChan chan string = make(chan string, 300)

var colorMax float64 = float64(0xFFFF)
var eightMax float64 = float64(0xFF)

func init() {
	for i := 0; i < runtime.NumCPU()-1; i++ {
		go analyzer()
	}
}

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
	avg := avgColor(img)
	log.Println("Avg Color in:", diffMs(calcAvg, time.Now()), " ms")
	log.Println("Avg Color:", avg)
	AnalysisCache.Add(&AnalyzedImage{
		fileName,
		avg,
	})
}

func avgColor(img image.Image) color.RGBA {
	bounds := img.Bounds()
	min := bounds.Min
	max := bounds.Max
	var r, g, b float64
	pixels := float64(bounds.Size().X * bounds.Size().Y)
	for x := min.X; x < max.X; x++ {
		for y := min.Y; y < max.Y; y++ {
			color := img.At(x, y)
			cr, cg, cb, _ := color.RGBA()
			r += float64(cr) / pixels
			g += float64(cg) / pixels
			b += float64(cb) / pixels
		}
	}
	log.Println("r g b p", r, g, b, pixels)
	return color.RGBA{
		cVal(r),
		cVal(g),
		cVal(b),
		uint8(0xFF),
	}
}

func cVal(p float64) uint8 {
	return uint8((p / colorMax) * eightMax)
}
