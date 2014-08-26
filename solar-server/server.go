package main

import (
	"github.com/dfjones/solar/solar-server/image-resource"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {

	r := mux.NewRouter()

	image_resource.Register(r)

	http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, r))
	http.ListenAndServe(":3000", nil)
}
