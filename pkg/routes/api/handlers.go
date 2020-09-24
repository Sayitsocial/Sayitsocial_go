package api

import (
	"encoding/json"
	"net/http"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/event"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/event/bridge/eventattendee"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/event/bridge/eventhost"

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

	apiRouter.HandleFunc("/create/vol", volCreateHandler).Methods("POST")
	apiRouter.HandleFunc("/create/org", orgCreateHandler).Methods("POST")
	apiRouter.HandleFunc("/event", eventHandler).Methods("POST", "GET")
	apiRouter.HandleFunc("/event/host", eventHostBridge).Methods("GET")
	apiRouter.HandleFunc("/event/attendee", eventAttendeeBridge).Methods("GET")
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

// swagger:route GET /api/event/host event getEventHost
//
// Get hosts of event
// Atleast one param is required
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
func eventHostBridge(w http.ResponseWriter, r *http.Request) {
	var req eventHostReq
	err := decoder.Decode(&req, r.URL.Query())
	if err != nil {
		helpers.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := req.CastToModel()
	if err != nil {
		helpers.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := eventhost.Initialize(nil)
	defer model.Close()

	err = json.NewEncoder(w).Encode(model.Get(data))
	if err != nil {
		helpers.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/event/attendee event getEventAttendee
//
// Get attendees of event
// Atleast one param is required
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
func eventAttendeeBridge(w http.ResponseWriter, r *http.Request) {
	var req eventAttendeeReq
	err := decoder.Decode(&req, r.URL.Query())
	if err != nil {
		helpers.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := req.CastToModel()
	if err != nil {
		helpers.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := eventattendee.Initialize(nil)
	defer model.Close()

	err = json.NewEncoder(w).Encode(model.Get(data))
	if err != nil {
		helpers.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	model := event.Initialize(nil)
	defer model.Close()

	if r.Method == "GET" {
		// swagger:route GET /api/event event getEvent
		//
		// Get details of event
		// Atleast one param is required
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
		var req eventGetReq
		err := decoder.Decode(&req, r.URL.Query())
		if err != nil {
			helpers.LogError(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data, err := req.CastToModel()
		if err != nil {
			helpers.LogError(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.NewEncoder(w).Encode(model.Get(data))
		if err != nil {
			helpers.LogError(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" {
		// swagger:route POST /api/event event createEvent
		//
		// Create a new event
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
		var req eventPostReq
		err := decoder.Decode(&req, r.URL.Query())
		if err != nil {
			helpers.LogError(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = req.PutInDB()
		if err != nil {
			helpers.LogError(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
