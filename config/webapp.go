package config

import (
	"net/http"
	"strings"
)

type WebApp struct {
	Base string
}

func (w *WebApp) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if strings.Contains(req.URL.Path, ".") {
		http.FileServer(http.Dir(w.Base)).ServeHTTP(rw, req)
		return
	}
	http.ServeFile(rw, req, w.Base+"/index.html")
}
