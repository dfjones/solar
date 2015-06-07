package image_analysis

import (
	libcolor "github.com/dfjones/solar/solar-server/lib/color"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

var analyzeChan chan string = make(chan string, 300)

var colorMax float64 = float64(0xFFFF)
var eightMax float64 = float64(0xFF)
var hMax float64 = float64(360)

func init() {
	go analyzer()
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
	avg := AvgColor(img)
	log.Println("Avg Color:", avg, diffMs(calcAvg, time.Now()), "ms")
	AnalysisCache.Add(&AnalyzedImage{
		fileName,
		avg,
	})
}

func AvgColor(img image.Image) color.HSL {
	bounds := img.Bounds()
	min := bounds.Min
	max := bounds.Max
	pixels := uint64(bounds.Size().X * bounds.Size().Y)
	cores := runtime.NumCPU()
	pr := make([]uint64, cores)
	pg := make([]uint64, cores)
	pb := make([]uint64, cores)
	var wg sync.WaitGroup
	for i := 0; i < cores; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			hbins := make([]float64, int(hMax))
			var as, al float64
			for y := min.Y + i; y < max.Y; y += cores {
				for x := min.X; x < max.X; x++ {
					p := img.At(x, y)
					pr, pg, pb, _ := p.RGBA()
					h, s, l := libcolor.RGBToHSL(pr, pg, pb)
					hbinIndex := round(h * hMax)
					hbins[hbinIndex]++
					as += s
					al += l
				}
			}
			pr[i] = r
			pg[i] = g
			pb[i] = b
		}(i)
	}
	wg.Wait()
	r := sum(pr) / pixels
	g := sum(pg) / pixels
	b := sum(pb) / pixels
	//log.Println("r g b p", r, g, b, pixels)
	return color.RGBA{
		cVal(r),
		cVal(g),
		cVal(b),
		uint8(0xFF),
	}
}

func sum(a []uint64) uint64 {
	var s uint64
	for i := 0; i < len(a); i++ {
		s += a[i]
	}
	return s
}

func round(v float64) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}

func cVal(p uint64) uint8 {
	return uint8((float64(p) / float64(colorMax)) * eightMax)
}
