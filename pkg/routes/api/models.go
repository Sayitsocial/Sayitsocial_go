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
	// in: body
	FirstName string `schema:"first_name,required"`

	// Last name of user
	// required: true
	// in: body
	LastName string `schema:"last_name,required"`

	// Email of user
	// required: true
	// in: body
	Email string `schema:"email,required"`

	// Password of user
	// required: true
	// in: body
	Password string `schema:"password,required"`
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

type OrgType int

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
	// in: body
	Email string `schema:"email,required"`

	// Password of user
	// required: true
	// in: body
	Password string `schema:"password,required"`

	// Name of Organisation
	// required: true
	// in: body
	OrgName string `schema:"org_name,required"`

	// Type of Organisation
	// required: true
	// in: body
	TypeOfOrg OrgType `schema:"org_type,required"`

	// Locality of Organisation
	// required: true
	// in: body
	Locality string `schema:"locality,required"`

	// Owner of Organisation
	// required: true
	// in: body
	Owner string `schema:"owner,required"`

	// Registration Number of organisation according to ngodarpan if applicable
	// required: false
	// in: body
	RegistrationNo string `schema:"reg_no,required"`
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
