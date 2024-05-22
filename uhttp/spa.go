package uhttp

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

type SpaHandler struct {
	fs   fs.FS
	path string
}

func NewSpaHandler(f fs.FS) *SpaHandler {
	o := new(SpaHandler)
	o.fs = f
	o.path = "/"
	return o
}

func (o *SpaHandler) WithPath(path string) *SpaHandler {
	o.path = strings.TrimSuffix(path, "/")
	return o
}

func (o *SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, o.path) {
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}
	fmt.Println(o.path)
	path := strings.TrimPrefix(r.URL.Path, o.path)
	fmt.Println("path:", path)
	if o.path != "/" && path != "" {
		if !strings.HasPrefix(path, "/") {
			http.NotFoundHandler().ServeHTTP(w, r)
			return
		}
	}
	path = strings.TrimPrefix(path, "/")
	if path != "" {
		fi, err := fs.Stat(o.fs, path)
		if err != nil && !os.IsNotExist(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if os.IsNotExist(err) || fi.IsDir() {
			r.URL.Path = o.path
		}
	}
	http.StripPrefix(strings.TrimSuffix(o.path, "/"), http.FileServer(http.FS(o.fs))).ServeHTTP(w, r)
}
