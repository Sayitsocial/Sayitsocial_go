package orgdata

import (
	"database/sql"
	"encoding/json"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models"
)

const (
	tableName = "organisation"
	schema    = helpers.DbSchemaOrg
)

// swagger:model
type OrgData struct {
	OrganisationID string                 `row:"organisation_id" type:"exact" pk:"manual" json:"organisation_id"`
	DisplayName    string                 `row:"display_name" type:"like" json:"display_name"`
	Locality       string                 `row:"locality" type:"like" json:"locality,omitempty"`
	RegistrationNo string                 `row:"registration_no" type:"exact" json:"registration_no,omitempty"`
	ContactEmail   string                 `row:"contact_email" type:"like" json:"contact_email,omitempty"`
	ContactPhone   string                 `row:"contact_phone" type:"like" json:"contact_phone"`
	Desc           string                 `row:"description" type:"like" json:"desc,omitempty"`
	Owner          string                 `row:"owner" type:"like" json:"owner,omitempty"`
	Achievements   string                 `row:"achievements" type:"like" json:"achievements,omitempty"`
	TypeOfOrg      int                    `row:"type_of_org" type:"like" json:"type_of_org"`
	Location       models.GeographyPoints `row:"location" type:"onlyvalue" json:"location"`
	Short          bool                   `scan:"ignore" json:"-"`
}

type Model struct {
	trans *sql.Tx
	conn  *sql.DB
}

func (o *OrgData) MarshalJSON() ([]byte, error) {
	type tmp OrgData
	//cat := &e.Category
	helpers.LogInfo(o.Short)
	if o.Short {
		o.RegistrationNo = ""
		o.ContactEmail = ""
		o.ContactPhone = ""
		o.Desc = ""
		o.Owner = ""
		o.Achievements = ""
		o.Locality = ""
	}
	return json.Marshal(&struct {
		*tmp
	}{
		(*tmp)(o),
	})
}

// Initialize returns model of db with active connection
func Initialize(tx *sql.Tx) *Model {
	if tx != nil {
		return &Model{
			trans: tx,
		}
	}
	return &Model{
		conn: models.GetConn(schema, tableName),
	}
}

// Close closes the connection to db
// Model should not be used after close is called
func (a Model) Close() {
	err := a.conn.Close()
	if err != nil {
		helpers.LogError(err.Error())
	}
}

// Create creates a value in database
func (a Model) Create(data OrgData) error {
	query, args := models.QueryBuilderCreate(data, schema+"."+tableName)

	var err error
	if a.trans != nil {
		_, err = a.trans.Exec(query, args...)
	} else {
		_, err = a.conn.Exec(query, args...)
	}
	return err
}

// Get data from db into slice of struct
// Searches by the member provided in input struct
func (a Model) Get(data OrgData) (orgData []OrgData) {
	query, args := models.QueryBuilderGet(data, schema+"."+tableName)
	helpers.LogInfo(query)
	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}
	models.GetIntoStruct(row, &orgData)
	if data.Short {
		for i := range orgData {
			orgData[i].Short = true
		}
	}
	return
}
