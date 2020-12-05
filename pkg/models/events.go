package models

import (
	"encoding/json"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database/querybuilder"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
)

// swagger:model
type Event struct {
	EventID       string                       `row:"event_id" type:"exact" json:"event_id" pk:"manual"`
	Name          string                       `row:"name" type:"like" json:"name"`
	Description   string                       `row:"description" type:"like" json:"description,omitempty"`
	StartTime     int64                        `row:"start_time" type:"exact" json:"start_time,omitempty"`
	HostTime      int64                        `row:"host_time" type:"exact" json:"host_time,omitempty"`
	Location      querybuilder.GeographyPoints `row:"location" type:"onlyvalue" json:"location"`
	TypeOfEvent   int64                        `row:"type_of_event" type:"exact" json:"type_of_event"`
	Category      EventCategory                `row:"category" type:"exact" ft:"events.event_category" fk:"generated_id" json:"-"`
	TrendingIndex int64                        `row:"trending_index" type:"exact" json:"trending_index"`
	SortBy        querybuilder.SortBy          `type:"sort" json:"-"`
	Short         bool                         `scan:"ignore" json:"-"`
	Page          querybuilder.Page            `json:"page"`
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
