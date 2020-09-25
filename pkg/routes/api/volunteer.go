package api

import (
	"encoding/json"
	"net/http"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/orgdata"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/common"
)

// Signup details for Volunteer
//
//swagger:parameters createVolunteer
type volCreReq struct {

	// First name of user
	// required: true
	// in: query
	FirstName string `schema:"first_name,required" json:"first_name"`

	// Last name of user
	// required: true
	// in: query
	LastName string `schema:"last_name,required" json:"last_name"`

	// Email of user
	// required: true
	// in: query
	Email string `schema:"email,required" json:"email"`

	// Password of user
	// required: true
	// in: query
	Password string `schema:"password,required" json:"password"`
}

// swagger:route POST /api/vol/create volunteer createVolunteer
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
		common.WriteError(err.Error(), w)
		return
	}

	helpers.LogInfo(req)

	err = req.PutInDB()
	if err != nil {
		helpers.LogError(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		common.WriteError(err.Error(), w)
		return
	}
	err = common.WriteSuccess(w)
	if err != nil {
		helpers.LogError(err.Error())
	}
}

//swagger:parameters getVolunteer
type volGetReq struct {

	// Organisation ID
	// in: query
	VolunteerID string `schema:"organisation_id" json:"organisation_id"`

	// Name of organisation
	// in: query
	DisplayName string `schema:"display_name" json:"display_name"`
}

// swagger:route GET /api/vol/get volunteer getVolunteer
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
func volGetHandler(w http.ResponseWriter, r *http.Request) {
	var req orgGetReq
	err := decoder.Decode(&req, r.URL.Query())
	if err != nil {
		helpers.LogError(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		common.WriteError(err.Error(), w)
		return
	}
	data, err := req.CastToModel()
	if err != nil {
		helpers.LogError(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		common.WriteError(err.Error(), w)
		return
	}

	model := orgdata.Initialize(nil)
	defer model.Close()

	err = json.NewEncoder(w).Encode(model.Get(data))
	if err != nil {
		helpers.LogError(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		common.WriteError(err.Error(), w)
		return
	}
	err = common.WriteSuccess(w)
	if err != nil {
		helpers.LogError(err.Error())
	}
}
