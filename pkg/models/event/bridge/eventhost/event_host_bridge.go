package eventhost

import (
	"database/sql"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/event"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/orgdata"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/voldata"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models"
)

const (
	tableName = "event_host_bridge"
	schema    = "public"
)

type EventHostBridge struct {
	GeneratedID  string          `row:"generated_id" type:"exact" pk:"manual" json:"generated_id"`
	Organisation orgdata.OrgData `row:"organisation_id" type:"exact" json:"organisation_id" fk:"organisation.organisation" fr:"organisation"`
	Volunteer    voldata.VolData `row:"volunteer_id" type:"exact" json:"volunteer_id" fk:"volunteer.volunteer" fr:"volunteer"`
	Event        event.Event     `row:"event_id" type:"exact" json:"event_id" fk:"public.events" fr:"event"`
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
func (a Model) Create(data EventHostBridge) error {
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
func (a Model) Get(data EventHostBridge) (eventHostBridge []EventHostBridge) {
	query, args := models.QueryBuilderJoin(data, schema+"."+tableName)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}
	models.GetIntoStruct(row, &eventHostBridge)
	return
}

// Count gets count of rows corresponsing to provided search params
func (a Model) Count(data EventHostBridge) (count []int) {
	query, args := models.QueryBuilderCount(data, schema+"."+tableName)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}

	models.GetIntoVar(row, &count)
	return
}
