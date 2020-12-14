package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database/querybuilder"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models"

	"github.com/google/uuid"
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
	tx, err := querybuilder.GetTransaction(ctx, nil)
	if err != nil {
		return err
	}

	modelAuth, err := querybuilder.Initialize(nil, tx)
	if err != nil {
		helpers.LogError(err.Error())
	}

	uid := uuid.New().String()

	err = modelAuth.Create(models.Auth{
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

	modelData, err := querybuilder.Initialize(nil, tx)
	if err != nil {
		helpers.LogError(err.Error())
	}

	err = modelData.Create(models.VolData{
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
	model, err := querybuilder.Initialize(nil, nil)
	if err != nil {
		helpers.LogError(err.Error())
	}
	defer model.Close()
	uid := uuid.New().String()
	err = model.Create(models.Followers{
		GeneratedID:    uid,
		OrganisationID: f.OrganisationID,
		Volunteer: models.VolData{
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
	tx, err := querybuilder.GetTransaction(ctx, nil)
	if err != nil {
		return err
	}

	modelAuth, err := querybuilder.Initialize(nil, tx)
	if err != nil {
		helpers.LogError(err.Error())
	}

	uid := uuid.New().String()

	err = modelAuth.Create(models.Auth{
		UID:        uid,
		Username:   o.Email,
		Password:   o.Password,
		TypeOfUser: helpers.AuthTypeOrg,
	})

	if err != nil {
		tx.Rollback()
		return err
	}

	modelData, err := querybuilder.Initialize(nil, tx)
	if err != nil {
		helpers.LogError(err.Error())
	}

	err = modelData.Create(models.OrgData{
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
	tx, err := querybuilder.GetTransaction(ctx, nil)
	if err != nil {
		return err
	}

	categoryModel, err := querybuilder.Initialize(nil, nil)
	if err != nil {
		helpers.LogError(err.Error())
	}
	defer categoryModel.Close()

	// TODO: Use count here
	x, err := categoryModel.Get(models.EventCategory{
		GeneratedID: e.Category,
	})

	if err != nil {
		return err
	}

	if len(*x.(*[]models.EventCategory)) == 0 {
		return errors.New("Invalid category ID")
	}

	if e.Location.Latitude != "" && e.Location.Longitude != "" && e.Location.Radius != "" {
		return errors.New("Invalid location [Should be Longitude, Latitude, Radius]")
	}

	eventModel, err := querybuilder.Initialize(nil, tx)
	if err != nil {
		helpers.LogError(err.Error())
	}

	eventID := uuid.New().String()

	err = eventModel.Create(models.Event{
		EventID:     eventID,
		Name:        e.Name,
		Description: e.Description,
		StartTime:   e.StartTime,
		HostTime:    time.Now().Unix(),
		Category: models.EventCategory{
			GeneratedID: e.Category,
		},
		TypeOfEvent: e.TypeOfEvent,
		Location:    e.Location,
	})

	if err != nil {
		tx.Rollback()
		return err
	}

	eventHostBridgeModel, err := querybuilder.Initialize(nil, tx)
	if err != nil {
		helpers.LogError(err.Error())
	}

	err = eventHostBridgeModel.Create(models.EventHostBridge{
		GeneratedID: uuid.New().String(),
		Organisation: models.OrgData{
			OrganisationID: e.OrganisationID,
		},
		Volunteer: models.VolData{
			VolunteerID: e.VolunteerID,
		},
		Event: models.Event{
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
func (e eventGetReq) CastToModel() (models.Event, error) {
	if e.EventID == "" && e.Name == "" && e.Category == 0 && e.StartTime == 0 && e.HostTime == 0 {
		return models.Event{}, errors.New("Requires one parameter")
	}
	if e.Location.Latitude != "" && e.Location.Longitude != "" && e.Location.Radius != "" {
		return models.Event{}, errors.New("Invalid location [Should be Longitude, Latitude, Radius]")
	}

	return models.Event{
		EventID:   e.EventID,
		Name:      e.Name,
		HostTime:  e.HostTime,
		StartTime: e.StartTime,
		Category: models.EventCategory{
			GeneratedID: e.Category,
		},
		TypeOfEvent: e.TypeOfEvent,
		Location:    e.Location,
		Short:       e.Short,
		// BUG: Gorilla decoder cant parse arrays properly sometimes
		SortBy: e.SortBy,
		Page: querybuilder.Page{
			Limit:  helpers.MaxPage,
			Offset: helpers.MaxPage * e.Page,
		},
	}, nil
}

// CastToModel converts request struct to model struct
func (e eventHostReq) CastToModel() (models.EventHostBridge, error) {
	if e.GeneratedID == "" && e.OrganisationID == "" && e.VolunteerID == "" && e.EventID == "" {
		return models.EventHostBridge{}, errors.New("Requires one parameter")
	}
	return models.EventHostBridge{
		GeneratedID: e.GeneratedID,
		Organisation: models.OrgData{
			OrganisationID: e.OrganisationID,
		},
		Volunteer: models.VolData{
			VolunteerID: e.VolunteerID,
		},
		Event: models.Event{
			EventID: e.EventID,
		},
	}, nil
}

// CastToModel converts request struct to model struct
func (e eventAttendeeReq) CastToModel() (models.EventAttendeeBridge, error) {
	if e.GeneratedID == "" && e.VolunteerID == "" && e.EventID == "" {
		return models.EventAttendeeBridge{}, errors.New("Requires one parameter")
	}
	return models.EventAttendeeBridge{
		GeneratedID: e.GeneratedID,
		Volunteer: models.VolData{
			VolunteerID: e.VolunteerID,
		},
		Event: models.Event{
			EventID: e.EventID,
		},
	}, nil
}

// CastToModel converts request struct to model struct
func (e orgGetReq) CastToModel() (models.OrgData, error) {
	return models.OrgData{
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
func (e volGetReq) CastToModel() (models.VolData, error) {
	return models.VolData{
		VolunteerID: e.VolunteerID,
		DisplayName: e.DisplayName,
	}, nil
}

func (f followerReq) RemoveFromDB() error {
	if f.OrganisationID == "" || f.VolunteerID == "" {
		return errors.New("All parameters are required")
	}

	model, err := querybuilder.Initialize(nil, nil)
	if err != nil {
		helpers.LogError(err.Error())
	}
	defer model.Close()

	return model.Delete(models.Followers{
		OrganisationID: f.OrganisationID,
		Volunteer: models.VolData{
			VolunteerID: f.VolunteerID,
		},
	})
}
