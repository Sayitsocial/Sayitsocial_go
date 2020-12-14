package api

import (
	"encoding/json"
	"net/http"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database/querybuilder"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/common"
)

// Signup details for Organisation
//
//swagger:model
type orgCreReq struct {

	// Email of Organisation
	// required: true
	Email string `schema:"email,required" json:"email"`

	// Password of user
	// required: true
	Password string `schema:"password,required" json:"password"`

	// Name of Organisation
	// required: true
	OrgName string `schema:"org_name,required" json:"org_name"`

	// Type of Organisation
	// required: true
	TypeOfOrg int `schema:"org_type,required" json:"org_type,string"`

	// Owner of Organisation
	// required: true
	Owner string `schema:"owner,required" json:"owner"`

	// Registration Number of organisation according to ngodarpan if applicable
	// required: false
	RegistrationNo string `schema:"reg_no,required" json:"reg_no"`

	// Location in [Longitude, Latitude]
	// minItems: 2
	// maxItems: 2
	Location querybuilder.GeographyPoints `schema:"location" json:"location"`
}

type orgCreModel struct {
	Organisation orgCreReq
}

// swagger:route POST /api/org/create organisation createOrganisation
//
// Create a new organisation
//
// This will show create a new organisation.
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
	err := readAndUnmarshal(r, &req)

	if err != nil {
		helpers.LogError("Error in GET parameters : " + err.Error())
		common.WriteError(err.Error(), http.StatusBadRequest, w)
		return
	}

	err = req.PutInDB()
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusInternalServerError, w)
		return
	}
	err = common.WriteSuccess(w)
	if err != nil {
		helpers.LogError(err.Error())
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
	TypeOfOrg int `schema:"type_of_org" json:"type_of_org,string"`

	// Location in [Longitude, Latitude, Radius]
	// in: query
	// minItems: 3
	// maxItems: 3
	Location querybuilder.GeographyPoints `schema:"location" json:"location"`

	// Sort results by [followers, ASC/DESC]
	// in: query
	SortBy querybuilder.SortBy `schema:"sortby" json:"sortby"`

	// Get short results
	// in: query
	Short bool `schema:"short" json:"short"`
}

// swagger:model
type orgDataShort struct {
	OrganisationID string                       `json:"organisation_id"`
	DisplayName    string                       `json:"display_name"`
	TypeOfOrg      int                          `json:"type_of_org"`
	Location       querybuilder.GeographyPoints `json:"location"`
	Followers      uint64                       `json:"follower_count"`
}

// swagger:response orgResponse
type orgResponse struct {
	org models.OrgData
}

// This response will be returned if "short" is true
// Status code will be 200
// swagger:response orgResponseShort
type orgShortResponse struct {
	org orgDataShort
}

// swagger:route GET /api/org/get organisation getOrganisation
//
// Get details of an organisation
//
// This will details of an organisation.
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
//       200: orgResponse
//		 201: orgResponseShort
func orgGetHandler(w http.ResponseWriter, r *http.Request) {
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

	err = json.NewEncoder(w).Encode(x.(*[]models.OrgData))
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusInternalServerError, w)
		return
	}
}
