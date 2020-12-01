// Basically a massive clusterfuck of I dont know what
// Should've used hardcoded queries :''(

package querybuilder

import (
	"fmt"
	"reflect"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
)

// QueryBuilderGet generates normal get queries for non nested structures
func QueryBuilderGet(i interface{}, tableName string) (string, []interface{}) {
	args, where := getArgsWhere(getSearchBy(i, tableName, true, false), 1)
	return fmt.Sprintf("SELECT %s FROM %s %s %s %s", getAllMembers(i, tableName, false), tableName, where, getOrderBy(i), getLimit(i)), args
}

// QueryBuilderCreate generates normal create queries for non nested structures
func QueryBuilderCreate(i interface{}, schema string, tableName string) (string, []interface{}) {
	members := getAllMembers(i, "", true)
	index := 1
	values, args := getValuesCount(i, &index, members)
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", schema+"."+tableName, members, values), args
}

// QueryBuilderDelete generates normal delete queries for non nested structures
func QueryBuilderDelete(i interface{}, tableName string) (string, []interface{}) {
	args, q := getArgsWhere(getSearchBy(i, tableName, false, false), 1)
	query := fmt.Sprintf("DELETE FROM %s %s", tableName, q)
	return query, args
}

// QueryBuilderUpdate generates normal update queries for non nested structures
func QueryBuilderUpdate(i interface{}, schema string, tableName string) (string, []interface{}) {
	members := getAllMembers(i, "", true)
	index := 1
	_, args := getValuesCount(i, &index, members)
	q, l := generateUpdateQuery(members)
	tmp, where := getArgsWhere(getSearchBy(i, tableName, false, true), l)
	args = append(args, tmp...)
	return fmt.Sprintf("UPDATE %s SET %s %s", schema+"."+tableName, q, where), args
}

// QueryBuilderJoin generates get queries for nested structures with inner join support
func QueryBuilderJoin(inte interface{}, tableName string) (string, []interface{}) {
	args, where := getArgsWhere(getSearchBy(inte, tableName, true, false), 1)
	helpers.LogInfo(getLimit(inte))
	query := fmt.Sprintf("SELECT %s FROM %s %s %s %s %s", getAllMembers(inte, tableName, false), tableName, getInnerJoin(inte, tableName), where, getOrderBy(inte), getLimit(inte))
	return query, args
}

// QueryBuilderCount generates count queries for primary key in structure
func QueryBuilderCount(inte interface{}, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(inte)
	v := reflect.ValueOf(inte)

	var pk string = "*"
	for i := 0; i < v.NumField(); i++ {
		if ok, _ := isPK(t.Field(i)); ok {
			pk = t.Field(i).Tag.Get(Row)
			break
		}
	}
	args, where := getArgsWhere(getSearchBy(inte, tableName, false, false), 1)

	return fmt.Sprintf("SELECT COUNT(%s) FROM %s %s", pk, tableName, where), args
}
