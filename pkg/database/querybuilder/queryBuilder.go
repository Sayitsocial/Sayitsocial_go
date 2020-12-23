// Basically a massive clusterfuck of I dont know what
// Should've used hardcoded queries :''(

package querybuilder

import (
	"fmt"
)

// queryBuilderCreate generates normal create queries for non nested structures
func queryBuilderCreate(cols []colHolder, foreign []foreignHolder, config Config, schema string, tableName string) (string, []interface{}) {
	query, indices := getAllMembers(cols, foreign, true)
	values, args := getInsertValues(indices, cols, foreign)
	return indexifyQuery(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", schema+"."+tableName, query, values)), args
}

// queryBuilderDelete generates normal delete queries for non nested structures
func queryBuilderDelete(cols []colHolder, foreign []foreignHolder, config Config, schema string, tableName string) (string, []interface{}) {
	q, args := getWhere(cols, foreign, false, false)
	query := fmt.Sprintf("DELETE FROM %s %s", schema+"."+tableName, q)
	return indexifyQuery(query), args
}

// queryBuilderUpdate generates normal update queries for non nested structures
func queryBuilderUpdate(cols []colHolder, foreign []foreignHolder, config Config, schema string, tableName string) (string, []interface{}) {
	q, args := generateUpdateQuery(cols, foreign)
	where, tmp := getWhere(cols, foreign, false, true)
	args = append(args, tmp...)
	return indexifyQuery(fmt.Sprintf("UPDATE %s SET %s %s", schema+"."+tableName, q, where)), args
}

// queryBuilderJoin generates get queries for nested structures with inner join support
func queryBuilderJoin(cols []colHolder, foreign []foreignHolder, config Config, schema string, tableName string) (string, []interface{}) {
	q, _ := getAllMembers(cols, foreign, false)
	where, args := getWhere(cols, foreign, true, false)
	query := fmt.Sprintf("SELECT %s FROM %s %s %s", q, schema+"."+tableName, getInnerJoin(foreign), where)

	if config.OrderBy != "" {
		query = fmt.Sprintf("%s ORDER BY %s %s", query, config.OrderBy, func() string {
			if config.OrderDesc {
				return "DESC"
			}
			return "ASC"
		}())
	}

	if config.Limit != 0 {
		query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, config.Limit, config.Offset)
	}
	return indexifyQuery(query), args
}

// queryBuilderCount generates count queries for primary key in structure
func queryBuilderCount(cols []colHolder, foreign []foreignHolder, config Config, schema string, tableName string) (string, []interface{}) {
	pk := "*"
	for _, col := range cols {
		if isPK(col.primary) && !col.isForeign {
			pk = col.col
			break
		}
	}
	where, args := getWhere(cols, foreign, false, false)
	return indexifyQuery(fmt.Sprintf("SELECT COUNT(%s) FROM %s %s", pk, schema+"."+tableName, where)), args
}
