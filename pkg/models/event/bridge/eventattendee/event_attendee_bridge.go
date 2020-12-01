package eventattendee

import (
	"database/sql"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database/querybuilder"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/event"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/models/volunteer/voldata"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
)

const (
	tableName = "event_attendee_bridge"
	schema    = helpers.DbSchemaEvents
)

// swagger:model
type EventAttendeeBridge struct {
	GeneratedID string          `row:"generated_id" type:"exact" pk:"manual" json:"generated_id"`
	Volunteer   voldata.VolData `row:"volunteer_id" type:"exact" json:"volunteer" ft:"volunteer.volunteer" fk:"volunteer_id"`
	Event       event.Event     `row:"event_id" type:"exact" json:"event" ft:"events.events" fk:"event_id"`
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
		conn: querybuilder.GetConn(schema, tableName),
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
func (a Model) Create(data EventAttendeeBridge) error {
	query, args := querybuilder.QueryBuilderCreate(data, schema, tableName)

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
func (a Model) Get(data EventAttendeeBridge) (eventAttendeeBridge []EventAttendeeBridge) {
	query, args := querybuilder.QueryBuilderJoin(data, schema+"."+tableName)
	helpers.LogInfo(query)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}
	querybuilder.GetIntoStruct(row, &eventAttendeeBridge)
	return
}

// Count gets count of rows corresponsing to provided search params
func (a Model) Count(data EventAttendeeBridge) (count []int) {
	query, args := querybuilder.QueryBuilderCount(data, schema+"."+tableName)

	row, err := a.conn.Query(query, args...)
	if err != nil {
		helpers.LogError(err.Error())
		return
	}

	querybuilder.GetIntoVar(row, &count)
	return
}
