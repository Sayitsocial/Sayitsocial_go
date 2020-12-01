package querybuilder

type SortBy struct {
	Column string `json:"column"`
	Mode   string `json:"mode"`
}

func (SortBy) memberSearchQuery(tableName string, rowTag string) string {
	return ""
}

func (SortBy) memberCreateQuery(tableName string, rowTag string) string {
	return ""
}

func (SortBy) whereQuery(tableName string, rowTag string) tmpHolder {
	return tmpHolder{}
}

func (SortBy) createArgs() string {
	return ""
}

func (s SortBy) isEmpty() bool {
	return (s.Column == "")
}

func (SortBy) ignoreScan() bool {
	return true
}
