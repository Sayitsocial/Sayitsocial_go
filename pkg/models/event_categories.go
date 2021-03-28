package models

import "github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"

type EventCategory struct {
	GeneratedID int    `sorm:"generated_id,pk_auto" json:"generated_id,omitempty"`
	Name        string `sorm:"name" json:"name,omitempty"`
}

func (EventCategory) GetTableName() (string, string) {
	return helpers.DbSchemaEvents, "event_category"
}
