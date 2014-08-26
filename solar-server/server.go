package main

import (
	"github.com/dfjones/solar/solar-server/image-resource"
	"github.com/gocraft/web"
	"net/http"
)

type Context struct {
}

func main() {

	router := web.New(Context{}).
		Middleware(web.LoggerMiddleware).
		Middleware(web.StaticMiddleware("./public")).
		Middleware(web.StaticMiddleware("/gopath/src/app/public"))

	image_resource.Register(router)

	http.ListenAndServe(":3000", router)
}
