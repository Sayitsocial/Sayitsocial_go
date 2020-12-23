package types

import (
	"fmt"
)

// GeographyPoints is a struct that holds longitude, latitude of a location and optionally radius
// Radius can be specified to initiate a search of other points in radius from specified coordinates
type GeographyPoints struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
	Radius    string `scan:"ignore" json:"-"`
}

func (GeographyPoints) memberSearchQuery(tableName string, rowTag string) string {
	return fmt.Sprintf("ST_X(%s.%s::geometry),ST_Y(%s.%s::geometry)", tableName, rowTag, tableName, rowTag)
}

func (g GeographyPoints) memberCreateQuery(tableName string, rowTag string) string {
	return rowTag
}

func (g GeographyPoints) whereQuery(tableName string, rowTag string) (string, []interface{}) {
	return fmt.Sprintf("ST_DWithin(%s.%s,ST_MakePoint(%v,%v),%v)", tableName, rowTag, g.Longitude, g.Latitude, g.Radius), make([]interface{}, 0)
}

func (g GeographyPoints) createArgs(i string) (string, []interface{}) {
	return fmt.Sprintf("ST_SetSRID(ST_MakePoint(%v,%v),4326)", g.Longitude, g.Latitude), make([]interface{}, 0)
}

func (g GeographyPoints) isEmpty() bool {
	return (g.Latitude == "" || g.Longitude == "")
}

func (GeographyPoints) ignoreScan() bool {
	return false
}
