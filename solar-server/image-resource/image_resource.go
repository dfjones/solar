package image_resource

import (
	"github.com/dfjones/solar/solar-server/image-analysis"
	"github.com/dfjones/solar/solar-server/image-storage"
	"github.com/gocraft/web"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
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
	serveAndClose(file, w, req.Request)
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
	serveAndClose(file, w, req.Request)
}

func getByPermalink(w web.ResponseWriter, req *web.Request) {
	link := req.PathParams["link"]
	// make sure we strip any escape/directory characters and get just a file name
	path := filepath.Base(link)
	if file, err := image_storage.GetByName(path); err == nil {
		serveAndClose(file, w, req.Request)
	} else {
		log.Println("error: ", err)
	}
}

func serveAndClose(file *os.File, w http.ResponseWriter, req *http.Request) {
	if file == nil {
		return
	}
	defer file.Close()
	w.Header().Set("cache-control", "public, max-age=300")
	http.ServeContent(w, req, file.Name(), time.Now(), file)
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
	w.WriteHeader(http.StatusOK)
}
