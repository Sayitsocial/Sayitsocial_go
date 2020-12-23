package types

// InbuiltType is an overloadable interface for custom sql types
type InbuiltType interface {

	// Custom query to be replaced while searching of database columns in SELECT operation
	SearchQuery(tableName string, rowTag string) string

	// Custom query to be replaced while searching of database columns in INSERT operation
	CreateQuery(tableName string, rowTag string) string

	// Returns custom arguments for create query
	CreateArgs(indexPlaceholder string) (string, []interface{})

	// Returns custom holder to parse where queries
	WhereQuery(tableName string, rowTag string) (string, []interface{})

	// Checks if struct is empty or not
	IsEmpty() bool

	// True if struct should be ignore while scanning values from DB
	IgnoreScan() bool
}
