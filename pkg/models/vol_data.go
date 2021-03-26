package models

import "github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"

// swagger:model
type VolData struct {
	VolunteerID  string `sorm:"volunteer_id,pk_manual" json:"volunteer_id"`
	DisplayName  string `sorm:"display_name" json:"display_name"`
	ContactEmail string `sorm:"contact_email" json:"contact_email"`
	ContactPhone string `sorm:"contact_phone" json:"contact_phone"`
	Bio          string `sorm:"bio" json:"bio"`
	Joined       int64  `sorm:"joined" json:"joined"`
}

func (VolData) OrgData() (string, string) {
	return helpers.DbSchemaVol, "volunteer"
}
