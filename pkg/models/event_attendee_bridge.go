package models

import "github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"

// swagger:model
type EventAttendeeBridge struct {
	GeneratedID string  `sorm:"generated_id,pk_manual" json:"generated_id"`
	Volunteer   VolData `sorm:"volunteer_id,ft_volunteer.volunteer,fk_volunteer_id" json:"volunteer"`
	Event       Event   `sorm:"event_id,ft_events.events,fk_event_id" json:"event"`
}

func (EventAttendeeBridge) GetTableName() (string, string) {
	return helpers.DbSchemaEvents, "event_attendee_bridge"
}
