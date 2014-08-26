package image_resource

import (
	"github.com/dfjones/solar/solar-server/image-storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Register(r *mux.Router) {
	r.HandleFunc("/images", postImage).Methods("POST")
	r.HandleFunc("/images/latest", getLatest).Methods("GET")
	r.HandleFunc("/images/{id:[0-9]+}", getByIndex).Methods("GET")
}

func getLatest(w http.ResponseWriter, req *http.Request) {
	file, err := image_storage.GetMostRecentImageFile()
	if err != nil {
		log.Println("error: ", err)
		return
	}
	serveAndClose(file, w, req)
}

func getByIndex(w http.ResponseWriter, req *http.Request) {
	index, err := strconv.Atoi(mux.Vars(req)["index"])
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

func postImage(w http.ResponseWriter, r *http.Request) {
	imageFile, _, err := r.FormFile("image")
	if err != nil {
		return
	}
	defer imageFile.Close()
	image_storage.Store(imageFile)
}
