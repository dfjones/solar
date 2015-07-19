package main

import (
	"net/http"
	"runtime"

	"github.com/dfjones/solar/solar-server/image-analysis"
	"github.com/dfjones/solar/solar-server/image-gif"
	"github.com/dfjones/solar/solar-server/image-resource"
	"github.com/dfjones/solar/solar-server/image-storage"
	"github.com/gocraft/web"
)

type Context struct {
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	router := web.New(Context{}).
		Middleware(web.LoggerMiddleware).
		Middleware(web.StaticMiddleware("./public")).
		Middleware(web.StaticMiddleware("/gopath/src/app/public"))

	image_resource.Register(router)
	image_analysis.Register(router)
	image_gif.Register(router)

	analyzePersisted()

	http.ListenAndServe(":3000", router)
}

func analyzePersisted() {
	paths := image_storage.GetAllPaths()
	for _, p := range paths {
		image_analysis.Analyze(p)
		image_gif.GetInstance().Submit(p)
	}
}
