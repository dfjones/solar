package image_gif

import (
	"github.com/nfnt/resize"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"log"
	"os"
)

var instance *GifGenerator

type Config struct {
	width     uint
	height    uint
	maxCount  uint
	delay     int
	loopCount int
	dir       string
	fileName  string
}

type GifGenerator struct {
	SubmitChan chan string
	Conf       Config
	gifd       gif.GIF
}

func init() {
	instance = NewGenerator()
}

func GetInstance() *GifGenerator {
	return instance
}

func NewGenerator() *GifGenerator {
	c := NewConfig()
	g := &GifGenerator{SubmitChan: make(chan string, 300),
		Conf: *c, gifd: gif.GIF{make([]*image.Paletted, 0), make([]int, 0), c.loopCount}}
	go run(g)
	return g
}

func NewConfig() *Config {
	return &Config{width: 648, height: 365, maxCount: 288, delay: 50, loopCount: 100, dir: "/data/solar/gif/", fileName: "timelapse.gif"}
}

func (g *GifGenerator) OpenGif() (*os.File, error) {
	return os.Open(g.Conf.dir + g.Conf.fileName)
}

func (g *GifGenerator) Submit(fileName string) {
	g.SubmitChan <- fileName
}

func run(g *GifGenerator) {
	for i := range g.SubmitChan {
		g.add(i)
		g.render()
	}
}

func (g *GifGenerator) add(jpegName string) {
	log.Println("Add gif:", jpegName)
	file, err := os.Open(jpegName)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	jpg, err := jpeg.Decode(file)
	if err != nil {
		log.Println("Error decoding file", err)
		return
	}

	m := resize.Resize(g.Conf.width, g.Conf.height, jpg, resize.Bicubic)

	bounds := m.Bounds()

	p := image.NewPaletted(bounds, palette.Plan9)

	sr := m.Bounds()
	draw.Draw(p, sr, m, sr.Min, draw.Src)

	if uint(len(g.gifd.Image)) > g.Conf.maxCount {
		g.gifd.Image = g.gifd.Image[1:]
	} else {
		g.gifd.Delay = append(g.gifd.Delay, g.Conf.delay)
	}

	g.gifd.Image = append(g.gifd.Image, p)
}

func (g *GifGenerator) render() {
	log.Println("Rendering gif...")
	file, err := os.Create(g.Conf.dir + g.Conf.fileName)
	if err != nil {
		log.Println("Error creating gif output file", err)
		return
	}
	defer file.Close()

	log.Println("gif data: ", g.gifd)
	err = gif.EncodeAll(file, &g.gifd)
	if err != nil {
		log.Println("Error encoding gif", err)
	}
}
