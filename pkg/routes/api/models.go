package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models"

	"github.com/google/uuid"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/auth"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/event"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/event/bridge/eventattendee"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/event/bridge/eventhost"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/event/categories"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/organisation/orgdata"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/volunteer/voldata"
)

// OrgType is type of organisation
type OrgType int

// enums for org types
const (
	NGO     OrgType = 0
	Company OrgType = 1
	Social  OrgType = 2
)

func (u volCreReq) PutInDB() error {
	ctx := context.Background()
	tx, err := database.GetConn().BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	modelAuth := auth.Initialize(tx)

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

func (o orgCreReq) PutInDB() error {
	ctx := context.Background()
	tx, err := database.GetConn().BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	modelAuth := auth.Initialize(tx)

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

	err = modelData.Create(orgdata.OrgData{
		OrganisationID: uid,
		DisplayName:    o.OrgName,
		Locality:       o.Locality,
		RegistrationNo: o.RegistrationNo,
		ContactEmail:   o.Email,
		Owner:          o.Owner,
		TypeOfOrg:      int(o.TypeOfOrg),
		Location: models.GeographyPoints{
			Longitude: o.Location[0],
			Latitude:  o.Location[1],
		},
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
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

	categoryModel := categories.Initialize(nil)
	defer categoryModel.Close()

	// TODO: Use count here
	if len(categoryModel.Get(categories.EventCategory{
		GeneratedID: e.Category,
	})) == 0 {
		return errors.New("Invalid category ID")
	}

	if len(e.Location) < 2 {
		return errors.New("Invalid location data")
	}

	eventModel := event.Initialize(tx)

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
		TypeOfEvent: e.TypeOfEvent,
		Location: models.GeographyPoints{
			Longitude: e.Location[0],
			Latitude:  e.Location[1],
		},
	})

	if err != nil {
		tx.Rollback()
		return err
	}

	eventHostBridgeModel := eventhost.Initialize(tx)

	err = eventHostBridgeModel.Create(eventhost.EventHostBridge{
		GeneratedID: uuid.New().String(),
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

// CastToModel converts request struct to model struct
func (e eventGetReq) CastToModel() (event.Event, error) {
	if e.EventID == "" && e.Name == "" && e.Category == 0 && e.StartTime == 0 && e.HostTime == 0 && len(e.Location) == 0 {
		return event.Event{}, errors.New("Requires one parameter")
	}
	if len(e.Location) < 3 && len(e.Location) != 0 {
		return event.Event{}, errors.New("Invalid location [Should be Longitude, Latitude, Radius]")
	}
	return event.Event{
		EventID:   e.EventID,
		Name:      e.Name,
		HostTime:  e.HostTime,
		StartTime: e.StartTime,
		Category: categories.EventCategory{
			GeneratedID: e.Category,
		},
		TypeOfEvent: e.TypeOfEvent,
		Location: func() models.GeographyPoints {
			if len(e.Location) < 3 {
				return models.GeographyPoints{}
			}
			return models.GeographyPoints{
				Longitude: e.Location[0],
				Latitude:  e.Location[1],
				Radius:    e.Location[2],
			}
		}(),
		Short: e.Short,
	}, nil
}

// CastToModel converts request struct to model struct
func (e eventHostReq) CastToModel() (eventhost.EventHostBridge, error) {
	if e.GeneratedID == "" && e.OrganisationID == "" && e.VolunteerID == "" && e.EventID == "" {
		return eventhost.EventHostBridge{}, errors.New("Requires one parameter")
	}
	return eventhost.EventHostBridge{
		GeneratedID: e.GeneratedID,
		Organisation: orgdata.OrgData{
			OrganisationID: e.OrganisationID,
		},
		Volunteer: voldata.VolData{
			VolunteerID: e.VolunteerID,
		},
		Event: event.Event{
			EventID: e.EventID,
		},
	}, nil
}

// CastToModel converts request struct to model struct
func (e eventAttendeeReq) CastToModel() (eventattendee.EventAttendeeBridge, error) {
	if e.GeneratedID == "" && e.VolunteerID == "" {
		return eventattendee.EventAttendeeBridge{}, errors.New("Requires one parameter")
	}
	return eventattendee.EventAttendeeBridge{
		GeneratedID: e.GeneratedID,
		Volunteer: voldata.VolData{
			VolunteerID: e.VolunteerID,
		},
		Event: event.Event{
			EventID: e.EventID,
		},
	}, nil
}

// CastToModel converts request struct to model struct
func (e orgGetReq) CastToModel() (orgdata.OrgData, error) {
	return orgdata.OrgData{
		OrganisationID: e.OrganisationID,
		DisplayName:    e.DisplayName,
		Owner:          e.Owner,
		TypeOfOrg:      e.TypeOfOrg,
		Location: func() models.GeographyPoints {
			if len(e.Location) < 3 {
				return models.GeographyPoints{}
			}
			return models.GeographyPoints{
				Longitude: e.Location[0],
				Latitude:  e.Location[1],
				Radius:    e.Location[2],
			}
		}(),
		Short: e.Short,
	}, nil
}

// CastToModel converts request struct to model struct
func (e volGetReq) CastToModel() (voldata.VolData, error) {
	return voldata.VolData{
		VolunteerID: e.VolunteerID,
		DisplayName: e.DisplayName,
	}, nil
}
