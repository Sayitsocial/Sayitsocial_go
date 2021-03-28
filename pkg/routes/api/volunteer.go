package api

import (
	"encoding/json"
	"net/http"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database/querybuilder"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/common"
)

// Signup details for Volunteer
//
//swagger:model
type volCreReq struct {

	// First name of user
	// required: true
	FirstName string `schema:"first_name,required" json:"first_name"`

	// Last name of user
	// required: true
	LastName string `schema:"last_name,required" json:"last_name"`

	// Email of user
	// required: true
	Email string `schema:"email,required" json:"email"`

	// Password of user
	// required: true
	Password string `schema:"password,required" json:"password"`
}

//swagger:parameters createVolunteer
type volCreModel struct {
	Volunteer volCreReq
}

// swagger:route POST /api/vol/create volunteer createVolunteer
//
// Create a new volunteer
//
// This will create a new volunteer.
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
	var req volCreReq
	err := readAndUnmarshal(r, &req)

	if err != nil {
		helpers.LogError("Error in GET parameters : " + err.Error())
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

//swagger:parameters getVolunteer
type volGetReq struct {

	// Organisation ID
	// in: query
	VolunteerID string `schema:"organisation_id" json:"organisation_id"`

	// Name of organisation
	// in: query
	DisplayName string `schema:"display_name" json:"display_name"`
}

// swagger:response volResponse
type volResponse struct {
	vol models.VolData
}

// swagger:route GET /api/vol/get volunteer getVolunteer
//
// Get details of a volunteer
//
// This will show details of a volunteer.
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
//       200: volResponse
func volGetHandler(w http.ResponseWriter, r *http.Request) {
	var req orgGetReq
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

	err = json.NewEncoder(w).Encode(x.(*[]models.VolData))
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusInternalServerError, w)
		return
	}
}
