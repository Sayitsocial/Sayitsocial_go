package api

import (
	"encoding/json"
	"net/http"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/orgdata"
)

// Signup details for Organisation
//
//swagger:parameters createOrganisation
type orgCreReq struct {

	// Email of Organisation
	// required: true
	// in: query
	Email string `schema:"email,required" json:"email"`

	// Password of user
	// required: true
	// in: query
	Password string `schema:"password,required" json:"password"`

	// Name of Organisation
	// required: true
	// in: query
	OrgName string `schema:"org_name,required" json:"org_name"`

	// Type of Organisation
	// required: true
	// in: query
	TypeOfOrg OrgType `schema:"org_type,required" json:"org_type"`

	// Locality of Organisation
	// required: true
	// in: query
	Locality string `schema:"locality,required" json:"locality"`

	// Owner of Organisation
	// required: true
	// in: query
	Owner string `schema:"owner,required" json:"owner"`

	// Registration Number of organisation according to ngodarpan if applicable
	// required: false
	// in: query
	RegistrationNo string `schema:"reg_no,required" json:"reg_no"`
}

// swagger:route POST /api/org/create organisation createOrganisation
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

//swagger:parameters getOrganisation
type orgGetReq struct {

	// Organisation ID
	// in: query
	OrganisationID string `schema:"organisation_id" json:"organisation_id"`

	// Name of organisation
	// in: query
	DisplayName string `schema:"display_name" json:"display_name"`

	// Owner of organisation
	// in: query
	Owner string `schema:"owner" json:"owner"`

	// Type of organisation
	// in: query
	TypeOfOrg int `schema:"type_of_org" json:"type_of_org"`
}

// swagger:route GET /api/org/get organisation getOrganisation
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
func orgGetHandler(w http.ResponseWriter, r *http.Request) {
	var req orgGetReq
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

	model := orgdata.Initialize(nil)
	defer model.Close()

	err = json.NewEncoder(w).Encode(model.Get(data))
	if err != nil {
		helpers.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}