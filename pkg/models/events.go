package models

import (
	"encoding/json"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database/querybuilder/types"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
)

// swagger:model
type Event struct {
	EventID       string                `sorm:"event_id,pk_manual" json:"event_id" `
	Name          string                `sorm:"name" json:"name"`
	Description   string                `sorm:"description" json:"description,omitempty"`
	StartTime     int64                 `sorm:"start_time"  json:"start_time,omitempty"`
	HostTime      int64                 `sorm:"host_time"  json:"host_time,omitempty"`
	Location      types.GeographyPoints `sorm:"location"  json:"location"`
	TypeOfEvent   int64                 `sorm:"type_of_event"  json:"type_of_event"`
	Category      EventCategory         `sorm:"category,ft_events.event_category,fk_generated_id" json:"-"`
	TrendingIndex int64                 `sorm:"trending_index" json:"trending_index"`
	Short         bool                  `sorm:"short,ignore" json:"-"`
}

func (Event) GetTableName() (string, string) {
	return helpers.DbSchemaEvents, "events"
}

func (e *Event) MarshalJSON() ([]byte, error) {
	type tmp Event
	cat := &e.Category
	if e.Short {
		e.Description = ""
		e.StartTime = 0
		e.HostTime = 0
		cat = nil

	}
	return json.Marshal(&struct {
		*tmp
		*EventCategory `json:"category,omitempty"`
	}{
		(*tmp)(e),
		cat,
	})
}
