package image_gif

import (
	"github.com/dfjones/solar/solar-server/lib/request"
	"github.com/gocraft/web"
	"log"
)

func Register(r *web.Router) {
	r.Get("/gif", getGif)
}

func getGif(w web.ResponseWriter, req *web.Request) {
	file, err := GetInstance().OpenGif()
	if err != nil {
		log.Println("error", err)
		return
	}
	request.ServeAndClose(file, w, req.Request)
}
