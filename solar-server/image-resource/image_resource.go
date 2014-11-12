package image_resource

import (
	"github.com/dfjones/solar/solar-server/image-analysis"
	"github.com/dfjones/solar/solar-server/image-gif"
	"github.com/dfjones/solar/solar-server/image-storage"
	"github.com/dfjones/solar/solar-server/lib/request"
	"github.com/gocraft/web"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

func Register(r *web.Router) {
	r.Post("/images", postImage)
	r.Get("/images/latest", getLatest)
	r.Get("/images/:index", getByIndex)
	r.Get("/images/perm/:link", getByPermalink)
}

func getLatest(w web.ResponseWriter, req *web.Request) {
	file, err := image_storage.GetMostRecentImageFile()
	if err != nil {
		log.Println("error: ", err)
		return
	}
	request.ServeAndClose(file, w, req.Request)
}

func getByIndex(w web.ResponseWriter, req *web.Request) {
	index, err := strconv.Atoi(req.PathParams["index"])
	if err != nil {
		return
	}
	file, err := image_storage.GetByIndex(index)
	if err != nil {
		log.Println("error: ", err)
		return
	}
	request.ServeAndClose(file, w, req.Request)
}

func getByPermalink(w web.ResponseWriter, req *web.Request) {
	link := req.PathParams["link"]
	// make sure we strip any escape/directory characters and get just a file name
	path := filepath.Base(link)
	if file, err := image_storage.GetByName(path); err == nil {
		request.ServeAndClose(file, w, req.Request)
	} else {
		log.Println("error: ", err)
	}
}

func postImage(w web.ResponseWriter, r *web.Request) {
	imageFile, _, err := r.FormFile("image")
	if err != nil {
		return
	}
	defer imageFile.Close()
	file, err := image_storage.Store(imageFile)
	if err != nil {
		log.Println("Error storing image: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	image_analysis.Analyze(file)
	image_gif.GetInstance().Submit(file)

	w.WriteHeader(http.StatusOK)
}
