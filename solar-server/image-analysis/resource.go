package image_analysis

import (
	"encoding/json"
	"github.com/gocraft/web"
	"image/color"
	"log"
	"path"
	"strconv"
	"strings"
	"time"
)

type imageJson struct {
	Time         time.Time
	AverageColor color.RGBA
}

func Register(r *web.Router) {
	r.Get("/analysis", getAll)
}

func getAll(w web.ResponseWriter, req *web.Request) {
	entries := AnalysisCache.All()
	encoder := json.NewEncoder(w)
	err := encoder.Encode(entries)
	if err != nil {
		log.Println("Error:", err)
	}
}

func (ai *AnalyzedImage) MarshalJSON() ([]byte, error) {
	_, file := path.Split(ai.Path)
	parts := strings.Split(file, "-")
	unixtime, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}
	return json.Marshal(imageJson{
		time.Unix(unixtime, 0),
		ai.AverageColor,
	})
}
