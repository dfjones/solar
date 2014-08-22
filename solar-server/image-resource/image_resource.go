package image_resource

import (
	"github.com/dfjones/solar/solar-server/image-storage"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"strconv"
)

type imageForm struct {
	ImageUpload *multipart.FileHeader `form:"image" binding:"required"`
}

func Register(m *martini.ClassicMartini) {
	m.Post("/images", binding.MultipartForm(imageForm{}), postImage)
	m.Get("/images/latest", getLatest)
	m.Get("/images/:index", getByIndex)
}

func getLatest(log *log.Logger) []byte {
	file, err := image_storage.GetMostRecentImageFile()
	if err != nil {
		log.Println("error: ", err)
		return nil
	}
	return readAndClose(file)
}

func getByIndex(params martini.Params) []byte {
	index, err := strconv.Atoi(params["index"])
	if err != nil {
		return nil
	}
	file, err := image_storage.GetByIndex(index)
	if err != nil {
		log.Println("error: ", err)
		return nil
	}
	return readAndClose(file)
}

func readAndClose(file *os.File) []byte {
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

func postImage(imageForm imageForm) {
	imageFile, err := imageForm.ImageUpload.Open()
	if err != nil {
		return
	}
	defer imageFile.Close()
	image_storage.Store(imageFile)
}
