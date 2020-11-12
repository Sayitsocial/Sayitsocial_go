package followerbridge

import (
	"database/sql"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/volunteer/voldata"
)

const (
	tableName = "follower_bridge"
	schema    = helpers.DbSchemaOrg
)

// swagger:model
type Followers struct {
	GeneratedID    string          `row:"generated_id" type:"exact" pk:"manual" json:"generated_id"`
	OrganisationID string          `row:"organisation_id" type:"exact" pk:"manual" json:"organisation_id"`
	Volunteer      voldata.VolData `row:"volunteer_id" type:"exact" fk:"volunteer.volunteer" fr:"volunteer_id" json:"volunteer"`
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
func (a Model) Create(data Followers) error {
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
func (a Model) Get(data Followers) (orgData []Followers) {
	query, args := models.QueryBuilderJoin(data, schema+"."+tableName)
	helpers.LogInfo(query)
	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}
	models.GetIntoStruct(row, &orgData)
	return
}

func (a Model) Delete(data Followers) error {
	query, args := models.QueryBuilderDelete(data, schema+"."+tableName)
	var err error
	helpers.LogInfo(query)
	if a.trans != nil {
		_, err = a.trans.Exec(query, args...)
	} else {
		_, err = a.conn.Exec(query, args...)
	}
	return err
}
