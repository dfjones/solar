package main

import (
	"github.com/dfjones/solar/solar-server/image-storage"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"io/ioutil"
	"log"
	"mime/multipart"
)

type ImageForm struct {
	ImageUpload *multipart.FileHeader `form:"image" binding:"required"`
}

func imagePostFile(imageForm ImageForm) {
	imageFile, err := imageForm.ImageUpload.Open()
	if err != nil {
		return
	}
	defer imageFile.Close()
	image_storage.Store(imageFile)
}

func imageGet(log *log.Logger) []byte {
	file, err := image_storage.GetMostRecentImageFile()
	if err != nil {
		log.Println("error: ", err)
		return nil
	}
	if file == nil {
		return nil
	}
	defer file.Close()
	ret, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("error: ", err)
		return nil
	}
	return ret
}

func main() {

	m := martini.Classic()

	m.Post("/images", binding.MultipartForm(ImageForm{}), imagePostFile)
	m.Get("/images", imageGet)

	m.Run()
}
