package api

import (
	"net/http"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

type Api struct {
}

func (a Api) Register(r *mux.Router) {

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.StrictSlash(false)

	// apiRouter.Use(middleware.AuthMiddleware())

	apiRouter.HandleFunc("/create/vol", volCreateHandler).Methods("POST", "GET")
	apiRouter.HandleFunc("/create/org", orgCreateHandler).Methods("POST", "GET")
}

// swagger:route POST /api/create/vol user_creation createVolunteer
//
// Create a new volunteer
//
// This will show create a new volunteer.
//
//     Consumes:
//     - application/x-www-form-urlencoded
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//
//     Security:
//       cookieAuth
//
//     Responses:
//       200: successResponse
func volCreateHandler(w http.ResponseWriter, r *http.Request) {
	helpers.LogInfo("here")
	var req volCreReq
	err := decoder.Decode(&req, r.URL.Query())

	if err != nil {
		helpers.LogError("Error in GET parameters : " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	helpers.LogInfo(req)

	err = req.PutInDB()
	if err != nil {
		helpers.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// swagger:route POST /api/create/org user_creation createOrganisation
//
// Create a new organisation
//
// This will show create a new volunteer.
//
//     Consumes:
//     - application/x-www-form-urlencoded
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//
//     Security:
//       cookieAuth
//
//     Responses:
//       200: successResponse
func orgCreateHandler(w http.ResponseWriter, r *http.Request) {
	var req orgCreReq
	err := decoder.Decode(&req, r.URL.Query())

	if err != nil {
		helpers.LogError("Error in GET parameters : " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = req.PutInDB()
	if err != nil {
		helpers.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
