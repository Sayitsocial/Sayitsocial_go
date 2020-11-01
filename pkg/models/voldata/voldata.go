package voldata

import (
	"database/sql"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models"
)

const (
	tableName = helpers.DbTableVolunteerData
	schema    = helpers.DbSchemaVol
)

// swagger:model
type VolData struct {
	VolunteerID  string `row:"volunteer_id" type:"exact" pk:"manual" json:"volunteer_id"`
	DisplayName  string `row:"display_name" type:"like" json:"display_name"`
	ContactEmail string `row:"contact_email" type:"like" json:"contact_email"`
	ContactPhone string `row:"contact_phone" type:"like" json:"contact_phone"`
	Bio          string `row:"bio" type:"like" json:"bio"`
	Joined       int64  `row:"joined" type:"exact" json:"joined"`
}

type Model struct {
	trans *sql.Tx
	conn  *sql.DB
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
func (a Model) Create(data VolData) error {
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
func (a Model) Get(data VolData) (volData []VolData) {
	query, args := models.QueryBuilderGet(data, schema+"."+tableName)

	helpers.LogInfo(query)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}

	models.GetIntoStruct(row, &volData)
	return
}
