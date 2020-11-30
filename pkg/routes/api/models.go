package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/organisation/followerbridge"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/organisation/orgdata"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/volunteer/voldata"
)

// enums for org types
const (
	NGO     int = 0
	Company int = 1
	Social  int = 2
)

func readAndUnmarshal(r *http.Request, req interface{}) error {
	if r.Method == "GET" {
		err := decoder.Decode(req, r.URL.Query())
		helpers.LogInfo(r.URL.Query())
		return err
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, req)
	helpers.LogInfo(req)
	return err
}

func (u volCreReq) PutInDB() error {
	if u.Email == "" || u.Password == "" || u.FirstName == "" || u.LastName == "" {
		return errors.New("No parameters should be empty")
	}
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

	helpers.LogInfo(u.FirstName)
	helpers.LogInfo(u.LastName)

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

func (f followerReq) PutInDB() error {
	if f.OrganisationID == "" || f.VolunteerID == "" {
		return errors.New("No parameters should be empty")
	}
	model := followerbridge.Initialize(nil)
	defer model.Close()
	uid := uuid.New().String()
	err := model.Create(followerbridge.Followers{
		GeneratedID:    uid,
		OrganisationID: f.OrganisationID,
		Volunteer: voldata.VolData{
			VolunteerID: f.VolunteerID,
		},
	})
	return err
}

func (o orgCreReq) PutInDB() error {
	if o.Email == "" || o.Password == "" || o.OrgName == "" || o.Owner == "" {
		return errors.New("No parameters should be empty")
	}

	if o.Location.Latitude != "" && o.Location.Longitude != "" && o.Location.Radius != "" {
		return errors.New("Invalid location [Should be Longitude, Latitude, Radius]")
	}
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
		RegistrationNo: o.RegistrationNo,
		ContactEmail:   o.Email,
		Owner:          o.Owner,
		TypeOfOrg:      int(o.TypeOfOrg),
		Location:       o.Location,
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

	if e.Location.Latitude != "" && e.Location.Longitude != "" && e.Location.Radius != "" {
		return errors.New("Invalid location [Should be Longitude, Latitude, Radius]")
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
		Location:    e.Location,
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
	if e.EventID == "" && e.Name == "" && e.Category == 0 && e.StartTime == 0 && e.HostTime == 0 {
		return event.Event{}, errors.New("Requires one parameter")
	}
	if e.Location.Latitude != "" && e.Location.Longitude != "" && e.Location.Radius != "" {
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
		Location:    e.Location,
		Short:       e.Short,
		// BUG: Gorilla decoder cant parse arrays properly sometimes
		SortBy: e.SortBy,
		Page: models.Page{
			Limit:  helpers.MaxPage,
			Offset: helpers.MaxPage * e.Page,
		},
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
	if e.GeneratedID == "" && e.VolunteerID == "" && e.EventID == "" {
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
		Location:       e.Location,
		Short:          e.Short,
		// BUG: Gorilla decoder cant parse arrays properly sometimes
		SortBy: e.SortBy,
	}, nil
}

// CastToModel converts request struct to model struct
func (e volGetReq) CastToModel() (voldata.VolData, error) {
	return voldata.VolData{
		VolunteerID: e.VolunteerID,
		DisplayName: e.DisplayName,
	}, nil
}

func (f followerReq) RemoveFromDB() error {
	if f.OrganisationID == "" || f.VolunteerID == "" {
		return errors.New("All parameters are required")
	}

	model := followerbridge.Initialize(nil)
	defer model.Close()

	return model.Delete(followerbridge.Followers{
		OrganisationID: f.OrganisationID,
		Volunteer: voldata.VolData{
			VolunteerID: f.VolunteerID,
		},
	})
}
