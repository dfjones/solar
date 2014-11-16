package request

import (
	"net/http"
	"os"
	"time"
)

func ServeAndClose(file *os.File, w http.ResponseWriter, req *http.Request) {
	if file == nil {
		return
	}
	defer file.Close()
	w.Header().Set("cache-control", "public, max-age=86400")
	http.ServeContent(w, req, file.Name(), time.Now(), file)
}
