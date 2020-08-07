package routes

import (
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/api"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/authentication"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/sayitsocial"
	"github.com/gorilla/mux"
	"net/http"
)

type App interface {
	Register(r *mux.Router)
}

var apps = []App{api.Api{}, authentication.Authentication{}, sayitsocial.Jizzberry{}}

func RegisterApps(r *mux.Router) {
	for _, i := range apps {
		i.Register(r)
	}
}

func RegisterFileServer(r *mux.Router) {
	r.PathPrefix("/static").Handler(
		http.FileServer(http.Dir(helpers.StaticPath)),
	)

	r.PathPrefix("/devhtml").Handler(http.StripPrefix("/devhtml",
		http.FileServer(http.Dir(helpers.StaticPath)),
	))
}
