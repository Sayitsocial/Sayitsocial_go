package api

import (
	"encoding/json"
	"net/http"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/event"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/event/bridge/eventattendee"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/event/bridge/eventhost"
)

// Event details
//
//swagger:parameters getEventHost
type eventHostReq struct {

	// Generated ID
	// in: query
	GeneratedID string `schema:"generated_id" json:"generated_id"`

	// ID of host of event (org)
	// in: query
	OrganisationID string `schema:"organisation_id" json:"organisation_id"`

	// ID of host of event (user)
	// in: query
	VolunteerID string `schema:"volunteer_id" json:"volunteer_id"`

	// ID of event
	// in: query
	EventID string `schema:"event_id" json:"event_id"`
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

// Event details
//
//swagger:parameters getEventAttendee
type eventAttendeeReq struct {

	// Generated ID
	// in: query
	GeneratedID string `schema:"generated_id" json:"generated_id"`

	// ID of host of event (user)
	// in: query
	VolunteerID string `schema:"volunteer_id" json:"volunteer_id"`

	// ID of event
	// in: query
	EventID string `schema:"event_id" json:"event_id"`
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

// Event details
//
//swagger:parameters getEvent
type eventGetReq struct {

	// Event id of event
	// in: query
	EventID string `schema:"event_id" json:"event_id"`

	// Name of event
	// in: query
	Name string `schema:"name" json:"name"`

	// Start time of event [unix timestamp]
	// in: query
	StartTime int64 `schema:"start_time" json:"start_time"`

	// Host time of event [unix timestamp]
	// in: query
	HostTime int64 `schema:"host_time" json:"host_time"`

	// Type of category [Refer to event_category]
	// in: query
	Category int `schema:"category" json:"category"`
}

// swagger:route GET /api/event/get event getEvent
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
func eventGetHandler(w http.ResponseWriter, r *http.Request) {
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

	model := event.Initialize(nil)
	defer model.Close()
	err = json.NewEncoder(w).Encode(model.Get(data))
	if err != nil {
		helpers.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Event details
//
//swagger:parameters createEvent
type eventPostReq struct {

	// ID of host of event (org)
	// in: query
	OrganisationID string `schema:"organisation_id" json:"organisation_id"`

	// ID of host of event (user)
	// in: query
	VolunteerID string `schema:"volunteer_id" json:"volunteer_id"`

	// Name of event
	// in: query
	Name string `schema:"name,required" json:"name"`

	// Name of event
	// in: query
	Description string `schema:"description,required" json:"description"`

	// Start time of event [unix timestamp]
	// in: query
	StartTime int64 `schema:"start_time,required" json:"start_time"`

	// Type of category [Refer to event_category]
	// in: query
	Category int `schema:"category,required" json:"category"`
}

// swagger:route POST /api/event/create event createEvent
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
func eventCreateHandler(w http.ResponseWriter, r *http.Request) {
	model := event.Initialize(nil)
	defer model.Close()
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
