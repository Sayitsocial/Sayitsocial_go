// Basically a massive clusterfuck of I dont know what
// Should've used hardcoded queries :''(

package querybuilder

import (
	"fmt"
	"reflect"
)

// queryBuilderCreate generates normal create queries for non nested structures
func queryBuilderCreate(i interface{}, schema string, tableName string) (string, []interface{}) {
	members := getAllMembers(i, "", true)
	index := 1
	values, args := getValuesCount(i, &index, members)
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", schema+"."+tableName, members, values), args
}

// queryBuilderDelete generates normal delete queries for non nested structures
func queryBuilderDelete(i interface{}, schema string, tableName string) (string, []interface{}) {
	args, q := getArgsWhere(getSearchBy(i, tableName, false, false), 0)
	query := fmt.Sprintf("DELETE FROM %s %s", schema+"."+tableName, q)
	return query, args
}

// queryBuilderUpdate generates normal update queries for non nested structures
func queryBuilderUpdate(i interface{}, schema string, tableName string) (string, []interface{}) {
	members := getAllMembers(i, "", true)
	index := 1
	_, args := getValuesCount(i, &index, members)
	q, l := generateUpdateQuery(members)
	tmp, where := getArgsWhere(getSearchBy(i, tableName, false, true), l)
	args = append(args, tmp...)
	return fmt.Sprintf("UPDATE %s SET %s %s", schema+"."+tableName, q, where), args
}

// queryBuilderJoin generates get queries for nested structures with inner join support
func queryBuilderJoin(inte interface{}, schema string, tableName string) (string, []interface{}) {
	args, where := getArgsWhere(getSearchBy(inte, tableName, true, false), 0)
	query := fmt.Sprintf("SELECT %s FROM %s %s %s %s %s", getAllMembers(inte, schema+"."+tableName, false), schema+"."+tableName, getInnerJoin(inte, tableName), where, getOrderBy(inte), getLimit(inte))
	return query, args
}

// queryBuilderCount generates count queries for primary key in structure
func queryBuilderCount(inte interface{}, schema string, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(inte)
	v := reflect.ValueOf(inte)

	var pk string = "*"
	for i := 0; i < v.NumField(); i++ {
		if ok, _ := isPK(t.Field(i)); ok {
			pk = t.Field(i).Tag.Get(Row)
			break
		}
	}
	args, where := getArgsWhere(getSearchBy(inte, schema+"."+tableName, false, false), 1)

	return fmt.Sprintf("SELECT COUNT(%s) FROM %s %s", pk, schema+"."+tableName, where), args
}
