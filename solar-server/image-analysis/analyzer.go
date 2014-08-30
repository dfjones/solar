package image_analysis

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"time"
)

var analyzeChan chan string = make(chan string)
var colorMax float64 = float64(0xFFFF)
var eightMax float64 = float64(0xFF)

func init() {
	go analyzer()
}

func Analyze(fileName string) {
	analyzeChan <- fileName
}

func analyzer() {
	for f := range analyzeChan {
		analyze(f)
	}
}

func analyze(fileName string) {
	start := time.Now()
	log.Println("Decode start...")
	file, err := os.Open(fileName)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	img, err := jpeg.Decode(file)
	log.Println("Decode finish in:", time.Now().Unix()-start.Unix(), " ms")
	if err != nil {
		log.Println("Error decoding file:", err)
		return
	}
	start = time.Now()
	avg := avgColor(img)
	log.Println("Avg Color in:", time.Now().Unix()-start.Unix(), " ms")
	log.Println("Avg Color:", avg)

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
