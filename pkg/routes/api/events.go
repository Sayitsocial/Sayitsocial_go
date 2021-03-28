package api

import (
	"encoding/json"
	"net/http"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database/querybuilder"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/database/querybuilder/types"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/common"
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

// swagger:response eventHostResponse
type eventHostResp struct {
	// in: query
	eventHost models.EventHostBridge
}

// swagger:route GET /api/event/host event getEventHost
//
// Get hosts of event
//
// This will show hosts of an event.
// Atleast one param is required
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
//       200: eventHostResponse
func eventHostBridge(w http.ResponseWriter, r *http.Request) {
	var req eventHostReq
	err := readAndUnmarshal(r, &req)
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusBadRequest, w)
		return
	}
	data, err := req.CastToModel()
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusBadRequest, w)
		return
	}

	model, err := querybuilder.Initialize(nil, nil)
	if err != nil {
		helpers.LogError(err.Error())
	}
	defer model.Close()

	x, err := model.Get(data)
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusInternalServerError, w)
	}

	err = json.NewEncoder(w).Encode(x.(*[]models.EventHostBridge))
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusInternalServerError, w)
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

// swagger:response eventAttendeeResponse
type eventAttendeeResponse struct {
	// in: query
	eventAttendee models.EventAttendeeBridge
}

// swagger:route GET /api/event/attendee event getEventAttendee
//
// Get attendees of event
//
// This will show attendees of an event.
// Atleast one param is required
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
//		 - JWTAuth: []
//
//     Responses:
//       200: eventAttendeeResponse
func eventAttendeeBridge(w http.ResponseWriter, r *http.Request) {
	var req eventAttendeeReq
	err := readAndUnmarshal(r, &req)
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusBadRequest, w)
		return
	}
	data, err := req.CastToModel()
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusBadRequest, w)
		return
	}

	model, err := querybuilder.Initialize(nil, nil)
	if err != nil {
		helpers.LogError(err.Error())
	}
	defer model.Close()

	x, err := model.Get(data)
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusInternalServerError, w)
	}

	err = json.NewEncoder(w).Encode(x.(*[]models.EventAttendeeBridge))
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusInternalServerError, w)
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

	// Type of event [0 - Virtual, 1 - Physical]
	// in: query
	TypeOfEvent int64 `schema:"type_of_event" json:"type_of_event"`

	// Location in [Longitude, Latitude, Radius]
	// in: query
	// minItems: 3
	// maxItems: 3
	Location types.GeographyPoints `schema:"location" json:"location"`

	// Sort results by [trending_index/ST_XMin(location)]
	// in: query
	SortBy string `schema:"sortby" json:"sortby"`

	// Pagination
	// in: query
	Page int64 `schema:"page" json:"page"`

	// Get short results
	// in: query
	Short bool `schema:"short" json:"short"`
}

// swagger:model
type eventShort struct {
	EventID       string                `json:"event_id"`
	Name          string                `json:"name"`
	Location      types.GeographyPoints `json:"location"`
	TypeOfEvent   int64                 `json:"type_of_event"`
	TrendingIndex int64                 `json:"trending_index"`
}

// swagger:response eventResponse
type eventResponse struct {
	// in: query
	event models.Event
}

// This response will be returned if "short" is true
// Status code will be 200
// swagger:response eventShortResponse
type eventShortResponse struct {
	// in: query
	event eventShort
}

// swagger:route GET /api/event/get event getEvent
//
// Get details of event
//
//
// This will show details of event
// Atleast one param is required
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
//       - 200: eventResponse
//		 - 201: eventShortResponse
func eventGetHandler(w http.ResponseWriter, r *http.Request) {
	var req eventGetReq
	err := readAndUnmarshal(r, &req)
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusBadRequest, w)
		return
	}

	data, err := req.CastToModel()
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusBadRequest, w)
		return
	}

	model, err := querybuilder.Initialize(nil, nil)
	if err != nil {
		helpers.LogError(err.Error())
	}
	defer model.Close()
	x, err := model.Page(req.Page, helpers.MaxPage).Order(req.SortBy, false).Get(data)
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusInternalServerError, w)
	}

	err = json.NewEncoder(w).Encode(x.(*[]models.Event))
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusInternalServerError, w)
		return
	}
}

// Event details
//
//swagger:parameters getEventAll
type eventGetAllReq struct {
	// Sort results by [trending_index/ST_XMin(location)]
	// in: query
	SortBy string `schema:"sortby" json:"sortby"`

	// Pagination
	// in: query
	Page int64 `schema:"page" json:"page"`
}

// swagger:route GET /api/event/get/all event getEventAll
//
// Get all events
//
//
// This will show details of event
// Atleast one param is required
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
//       - 200: eventResponse
//		 - 201: eventShortResponse
func eventGetAllHandler(w http.ResponseWriter, r *http.Request) {
	var req eventGetAllReq
	err := readAndUnmarshal(r, &req)
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusBadRequest, w)
		return
	}
	model, err := querybuilder.Initialize(nil, nil)
	if err != nil {
		helpers.LogError(err.Error())
	}
	defer model.Close()
	data, err := req.CastToModel()
	x, err := model.Page(req.Page, helpers.MaxPage).Order(req.SortBy, false).Get(data)
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusInternalServerError, w)
	}

	err = json.NewEncoder(w).Encode(x.(*[]models.Event))
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusInternalServerError, w)
		return
	}
}

//swagger:model
type eventPostReq struct {

	// ID of host of event (org)
	OrganisationID string `schema:"organisation_id" json:"organisation_id"`

	// ID of host of event (user)
	VolunteerID string `schema:"volunteer_id" json:"volunteer_id"`

	// Name of event
	Name string `schema:"name,required" json:"name"`

	// Name of event
	Description string `schema:"description,required" json:"description"`

	// Start time of event [unix timestamp]
	StartTime int64 `schema:"start_time,required" json:"start_time"`

	// Type of category [Refer to event_category]
	Category int `schema:"category,required" json:"category"`

	// Type of category [0 - Virtual, 1 - Physical]
	TypeOfEvent int64 `schema:"type_of_event,required" json:"type_of_event"`

	// Location in [Longitude, Latitude]
	// minItems: 2
	// maxItems: 2
	Location types.GeographyPoints `schema:"location" json:"location"`
}

//swagger:parameters createEvent
type eventPostModel struct {
	//in: query
	EventsPostModel eventPostReq
}

// swagger:route POST /api/event/create event createEvent
//
// Create a new event
//
// This will show create a new event.
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
	model, err := querybuilder.Initialize(nil, nil)
	if err != nil {
		helpers.LogError(err.Error())
	}
	defer model.Close()

	var req eventPostReq
	err = readAndUnmarshal(r, &req)
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusBadRequest, w)
		return
	}
	err = req.PutInDB()
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusBadRequest, w)
		return
	}
	err = common.WriteSuccess(w)
	if err != nil {
		helpers.LogError(err.Error())
	}
}
