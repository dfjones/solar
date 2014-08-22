package main

import (
	"github.com/dfjones/solar/solar-server/image-resource"
	"github.com/go-martini/martini"
)

func main() {

	m := martini.Classic()

	image_resource.Register(m)

	m.Use(m.Static("/gopath/src/app/public"))

	m.Run()
}
