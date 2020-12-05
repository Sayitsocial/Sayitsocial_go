package models

import "github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"

type EventCategory struct {
	GeneratedID int    `row:"generated_id" type:"exact" json:"generated_id,omitempty" pk:"auto"`
	Name        string `row:"name" type:"exact" json:"name,omitempty"`
}

func (EventCategory) GetTableName() (string, string) {
	return helpers.DbSchemaEvents, "event_category"
}
