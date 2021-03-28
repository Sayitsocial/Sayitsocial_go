package types

import (
	"fmt"
)

// GeographyPoints is a struct that holds longitude, latitude of a location and optionally radius
// Radius can be specified to initiate a search of other points in radius from specified coordinates
type SormBool struct {
	Value        bool
	IsValueEmpty bool `scan:"ignore" json:"-"`
}

func (SormBool) SearchQuery(tableName string, rowTag string) string {
	return fmt.Sprintf("%s.%s", tableName, rowTag)
}

func (b SormBool) CreateQuery(rowTag string) string {
	return rowTag
}

func (b SormBool) WhereQuery(tableName string, rowTag string, indexPlaceholder string) (string, []interface{}) {
	args := make([]interface{}, 0)
	args = append(args, b.Value)
	return fmt.Sprintf("%s.%s=$%s", tableName, rowTag, indexPlaceholder), args
}

func (b SormBool) CreateArgs(indexPlaceholder string) (string, []interface{}) {
	args := make([]interface{}, 0)
	args = append(args, b.Value)
	return fmt.Sprintf("$%s", indexPlaceholder), args
}

func (b SormBool) IsEmpty() bool {
	return b.IsValueEmpty
}

func (SormBool) IgnoreScan() bool {
	return false
}
