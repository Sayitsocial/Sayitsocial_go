package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"

	"github.com/google/uuid"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/auth"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/event"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/event/bridge/eventhost"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/event/categories"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/orgdata"
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
	ctx := context.Background()
	tx, err := database.GetConn().BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	modelAuth := auth.Initialize(tx)
	defer modelAuth.Close()

	uid := uuid.New().String()

	err = modelAuth.Create(auth.Auth{
		UID:        uid,
		Username:   u.Email,
		Password:   u.Password,
		TypeOfUser: helpers.AuthTypeVol,
	})

	if err != nil {
		tx.Rollback()
		return err
	}

	modelData := voldata.Initialize(tx)
	defer modelData.Close()

	err = modelData.Create(voldata.VolData{
		VolunteerID:  uid,
		DisplayName:  fmt.Sprintf("%s %s", u.FirstName, u.LastName),
		ContactEmail: u.Email,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
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
	ctx := context.Background()
	tx, err := database.GetConn().BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	modelAuth := auth.Initialize(tx)
	defer modelAuth.Close()

	uid := uuid.New().String()

	err = modelAuth.Create(auth.Auth{
		UID:        uid,
		Username:   o.Email,
		Password:   o.Password,
		TypeOfUser: helpers.AuthTypeOrg,
	})

	if err != nil {
		tx.Rollback()
		return err
	}

	modelData := orgdata.Initialize(tx)
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
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
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

// CastToModel converts request struct to model struct
func (e eventGetReq) CastToModel() (event.Event, error) {
	if e.EventID == "" && e.Name == "" && e.Category == 0 && e.StartTime == 0 && e.HostTime == 0 {
		return event.Event{}, errors.New("Requires one parameter")
	}
	return event.Event{
		EventID:   e.EventID,
		Name:      e.Name,
		HostTime:  e.HostTime,
		StartTime: e.StartTime,
		Category: categories.EventCategory{
			GeneratedID: e.Category,
		},
	}, nil
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

func (e eventPostReq) PutInDB() error {
	if e.OrganisationID == "" && e.VolunteerID == "" {
		return errors.New("One of organisation_id or volunteer_id must be present")
	}

	ctx := context.Background()
	tx, err := database.GetConn().BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	categoryModel := categories.Initialize(tx)
	defer categoryModel.Close()

	// TODO: Use count here
	if len(categoryModel.Get(categories.EventCategory{
		GeneratedID: e.Category,
	})) == 0 {
		return errors.New("Invalid category ID")
	}

	eventModel := event.Initialize(tx)
	defer eventModel.Close()

	eventID := uuid.New().String()

	err = eventModel.Create(event.Event{
		EventID:     eventID,
		Name:        e.Name,
		Description: e.Description,
		StartTime:   e.StartTime,
		HostTime:    time.Now().Unix(),
		Category: categories.EventCategory{
			GeneratedID: e.Category,
		},
	})

	if err != nil {
		tx.Rollback()
		return err
	}

	eventHostBridgeModel := eventhost.Initialize(tx)
	defer eventHostBridgeModel.Close()

	err = eventHostBridgeModel.Create(eventhost.EventHostBridge{
		Organisation: orgdata.OrgData{
			OrganisationID: e.OrganisationID,
		},
		Volunteer: voldata.VolData{
			VolunteerID: e.VolunteerID,
		},
		Event: event.Event{
			EventID: eventID,
		},
	})

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
