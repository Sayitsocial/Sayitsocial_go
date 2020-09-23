package api

import (
	"fmt"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/auth"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/orgdata"
	"github.com/google/uuid"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/voldata"
)

type request interface {
	PutInDB() error
	Validate() bool
}

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

func (u volCreReq) PutInDB() error {
	helpers.LogInfo("here")
	modelAuth := auth.Initialize()
	defer modelAuth.Close()

	uid := uuid.New().String()

	err := modelAuth.Create(auth.Auth{
		UID:        uid,
		Username:   u.Email,
		Password:   u.Password,
		TypeOfUser: helpers.AuthTypeVol,
	})

	if err != nil {
		return err
	}

	modelData := voldata.Initialize()
	defer modelData.Close()

	err = modelData.Create(voldata.VolData{
		VolunteerID:  uid,
		DisplayName:  fmt.Sprintf("%s %s", u.FirstName, u.LastName),
		ContactEmail: u.Email,
	})
	return err
}

// OrgType is type of organisation
type OrgType int

// enums for org types
const (
	NGO     OrgType = 0
	Company OrgType = 1
	Social  OrgType = 2
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

func (o orgCreReq) PutInDB() error {
	modelAuth := auth.Initialize()
	defer modelAuth.Close()

	uid := uuid.New().String()

	err := modelAuth.Create(auth.Auth{
		UID:        uid,
		Username:   o.Email,
		Password:   o.Password,
		TypeOfUser: helpers.AuthTypeOrg,
	})

	if err != nil {
		return err
	}

	modelData := orgdata.Initialize()
	defer modelData.Close()

	err = modelData.Create(orgdata.OrgData{
		OrganisationID: uid,
		DisplayName:    o.OrgName,
		Locality:       o.Locality,
		RegistrationNo: o.RegistrationNo,
		ContactEmail:   o.Email,
		Owner:          o.Owner,
		TypeOfOrg:      int(o.TypeOfOrg),
	})
	return err
}
