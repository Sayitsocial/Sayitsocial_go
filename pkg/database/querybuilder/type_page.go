package querybuilder

// Page holds limit and offset to implement pagination
type Page struct {
	Limit  int64
	Offset int64
}

func (Page) memberSearchQuery(tableName string, rowTag string) string {
	return ""
}

func (Page) memberCreateQuery(tableName string, rowTag string) string {
	return ""
}

func (Page) whereQuery(tableName string, rowTag string) tmpHolder {
	return tmpHolder{}
}

func (Page) createArgs() string {
	return ""
}

func (s Page) isEmpty() bool {
	return (s.Limit == 0)
}

func (Page) ignoreScan() bool {
	return true
}
