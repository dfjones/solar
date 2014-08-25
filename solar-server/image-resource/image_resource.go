package image_resource

import (
	"github.com/dfjones/solar/solar-server/image-storage"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

type imageForm struct {
	ImageUpload *multipart.FileHeader `form:"image" binding:"required"`
}

func Register(m *martini.ClassicMartini) {
	m.Post("/images", binding.MultipartForm(imageForm{}), postImage)
	m.Get("/images/latest", getLatest)
	m.Get("/images/:index", getByIndex)
}

func getLatest(w http.ResponseWriter, req *http.Request, log *log.Logger) {
	file, err := image_storage.GetMostRecentImageFile()
	if err != nil {
		log.Println("error: ", err)
		return
	}
	serveAndClose(file, w, req)
}

func getByIndex(w http.ResponseWriter, req *http.Request, params martini.Params) {
	index, err := strconv.Atoi(params["index"])
	if err != nil {
		return
	}
	file, err := image_storage.GetByIndex(index)
	if err != nil {
		log.Println("error: ", err)
		return
	}
	serveAndClose(file, w, req)
}

func serveAndClose(file *os.File, w http.ResponseWriter, req *http.Request) {
	if file == nil {
		return
	}
	defer file.Close()
	w.Header().Set("cache-control", "public, max-age=300")
	http.ServeContent(w, req, file.Name(), time.Now(), file)
}

func postImage(imageForm imageForm) {
	imageFile, err := imageForm.ImageUpload.Open()
	if err != nil {
		return
	}
	defer imageFile.Close()
	image_storage.Store(imageFile)
}
