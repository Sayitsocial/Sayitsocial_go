package models

import "github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"

// swagger:model
type EventHostBridge struct {
	GeneratedID  string  `row:"generated_id" type:"exact" pk:"manual" json:"generated_id"`
	Organisation OrgData `row:"organisation_id" type:"exact" json:"organisation" ft:"organisation.organisation" fk:"organisation_id"`
	Volunteer    VolData `row:"volunteer_id" type:"exact" json:"volunteer" ft:"volunteer.volunteer" fk:"volunteer_id"`
	Event        Event   `row:"event_id" type:"exact" json:"event" ft:"events.events" fk:"event_id"`
}

func (EventHostBridge) GetTableName() (string, string) {
	return helpers.DbSchemaEvents, "event_host_bridge"
}
