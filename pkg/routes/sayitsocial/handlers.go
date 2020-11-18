package sayitsocial

import (
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/gorilla/mux"
)

type Sayitsocial struct {
}

const baseURL = ""

// Register registers the html router for static webpages
func (s Sayitsocial) Register(r *mux.Router) {
	htmlRouter := r.PathPrefix(baseURL).Subrouter()
	htmlRouter.StrictSlash(true)

	//htmlRouter.Use(middleware.AuthMiddleware())
	r.PathPrefix("/").HandlerFunc(reactAppHandler)

}

func reactAppHandler(w http.ResponseWriter, r *http.Request) {
	ext := path.Ext(r.URL.Path)
	if ext == ".html" || ext == "" {
		helpers.LogInfo("Serving index")
		http.ServeFile(w, r, filepath.Join(helpers.StaticPath2, "index.html"))
	} else {
		isStatic, _ := path.Match("/static/*/*", r.URL.Path)
		if isStatic {
			w.Header().Add("Cache-Control", "max-age=604800000")
		}
		split := strings.Split(r.URL.Path, "/")
		helpers.LogInfo(split[len(split)-1])
		if ext == ".js" || ext == ".css" {
			http.ServeFile(w, r, filepath.Join(helpers.StaticPath2, "static/", split[len(split)-2], split[len(split)-1]))
		} else {
			http.ServeFile(w, r, filepath.Join(helpers.StaticPath2, split[len(split)-2], split[len(split)-1]))
		}
	}
}
