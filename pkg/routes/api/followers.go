package api

import (
	"net/http"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/routes/common"
)

// Add follower details
//
//swagger:parameters addFollower
type followerCreReq struct {
	// ID of organisation who is being followed
	// in: query
	// required: true
	OrganisationID string `schema:"organisation_id" json:"organisation_id"`

	// ID of volunteer who is following
	// in: query
	// required: true
	VolunteerID string `schema:"volunteer_id" json:"volunteer_id"`
}

// swagger:route POST /api/followers/add followers addFollower
//
// Add a follower to certain organisation
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
func addFollowerHandler(w http.ResponseWriter, r *http.Request) {
	var req followerCreReq
	err := readAndUnmarshal(r, &req)
	if err != nil {
		helpers.LogError(err.Error())
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

// Remove follower details
//
//swagger:parameters delFollower
type followerDelReq struct {
	// ID of organisation who is being followed
	// in: query
	// required: true
	OrganisationID string `schema:"organisation_id,required" json:"organisation_id"`

	// ID of volunteer who is following
	// in: query
	// required: true
	VolunteerID string `schema:"volunteer_id,required" json:"volunteer_id"`
}

// swagger:route POST /api/followers/remove followers delFollower
//
// Remove a follower from certain organisation
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
func removeFollowerHandler(w http.ResponseWriter, r *http.Request) {
	var req followerDelReq
	err := readAndUnmarshal(r, &req)
	if err != nil {
		helpers.LogError(err.Error())
		common.WriteError(err.Error(), http.StatusBadRequest, w)
		return
	}
	err = req.RemoveFromDB()
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
