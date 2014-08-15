package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"sync"
)

type ImageForm struct {
	ImageUpload *multipart.FileHeader `form:"image" binding:"required"`
}

type LatestImage struct {
	bytes []byte
	mutex *sync.Mutex
}

var latestImage *LatestImage = &LatestImage{nil, &sync.Mutex{}}

func imagePostTemp(res http.ResponseWriter, log *log.Logger, imageForm ImageForm, errors binding.Errors) string {
	if errors != nil {
		log.Println("errors %s", errors)
		return "error"
	}

	file, err := imageForm.ImageUpload.Open()
	if err != nil {
		log.Println("error: %s", err)
	}
	defer file.Close()

	tempFile, err := ioutil.TempFile("/tmp", "solar")

	io.Copy(tempFile, file)

	return tempFile.Name()
}

func imagePostMem(imageForm ImageForm) {
	file, err := imageForm.ImageUpload.Open()
	if err != nil {
		return
	}
	defer file.Close()

	latestImage.mutex.Lock()
	defer latestImage.mutex.Unlock()
	latestImage.bytes, err = ioutil.ReadAll(file)
}

func imageGet(res http.ResponseWriter) []byte {
	latestImage.mutex.Lock()
	defer latestImage.mutex.Unlock()
	return latestImage.bytes
}

func main() {

	m := martini.Classic()

	m.Post("/images", binding.MultipartForm(ImageForm{}), imagePostMem)
	m.Get("/images", imageGet)

	m.Run()
}
