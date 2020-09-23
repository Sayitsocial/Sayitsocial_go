// Basically a massive clusterfuck of I dont know what
// Should've used hardcoded queries :''(

package models

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/google/uuid"
)

const (
	autoPK   = "auto"
	manualPK = "manual"
)

type tmpHolder struct {
	name   string
	typeOf string
	value  reflect.Value
}

func getAllMembers(inte interface{}, tableName string) string {
	t := reflect.TypeOf(inte)
	v := reflect.ValueOf(inte)

	var ret = ""
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct {
			if foreignTable := t.Field(i).Tag.Get("fk"); foreignTable != "" {
				ret += getAllMembers(v.Field(i).Interface(), foreignTable)
			} else {
				ret += getAllMembers(v.Field(i).Interface(), tableName)
			}
			continue
		} else {
			ret += fmt.Sprintf("%s.%s ", tableName, t.Field(i).Tag.Get(helpers.RowStructTag))
		}
		ret += ", "
	}
	return ret[:len(ret)-1]
}

func getSearchBy(inte interface{}, tableName string) (totalSearchByCount int, searchBy []tmpHolder) {
	t := reflect.TypeOf(inte)
	v := reflect.ValueOf(inte)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct {
			if foreignTable := t.Field(i).Tag.Get("fk"); foreignTable != "" {
				tmpCount, tmpSearchBy := getSearchBy(v.Field(i).Interface(), foreignTable)
				totalSearchByCount += tmpCount
				searchBy = append(searchBy, tmpSearchBy...)
			} else {
				tmpCount, tmpSearchBy := getSearchBy(v.Field(i).Interface(), tableName)
				totalSearchByCount += tmpCount
				searchBy = append(searchBy, tmpSearchBy...)
			}
			continue
		}
		if !checkEmpty(v.Field(i)) {
			searchBy = append(searchBy, tmpHolder{
				name:   fmt.Sprintf("%s.%s ", tableName, t.Field(i).Tag.Get(helpers.RowStructTag)),
				typeOf: "exact",
				value:  v.Field(i),
			})
			totalSearchByCount++
		}
	}
	return
}

func getArgsWhere(totalSearchByCount int, searchBy []tmpHolder) ([]interface{}, string) {
	args := make([]interface{}, 0)
	if len(searchBy) == 0 {
		return nil, ""
	}
	query := " WHERE "
	for searchByCount, r := range searchBy {
		if r.typeOf == "exact" {
			query += fmt.Sprintf("%v = $%d", r.name, searchByCount+1)
		} else if r.typeOf == "like" {
			query += fmt.Sprintf("%v ILIKE $%d", r.name, searchByCount+1)
		}
		if searchByCount < totalSearchByCount-1 {
			query += " AND "
		}
		args = append(args, r.value.Interface())
	}
	return args, query
}

// QueryBuilderGet generates normal get queries for non nested structures
func QueryBuilderGet(i interface{}, schemaName string, tableName string) (string, []interface{}) {
	query := `SELECT ` + getAllMembers(i, schemaName+"."+tableName)

	args, where := getArgsWhere(getSearchBy(i, schemaName+"."+tableName))
	query += " FROM " + fmt.Sprintf("%s.%s", schemaName, tableName)
	query += where

	return query, args
}

// QueryBuilderCreate generates normal create queries for non nested structures
func QueryBuilderCreate(i interface{}, schemaName string, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	query := `INSERT INTO ` + fmt.Sprintf("%s.%s", schemaName, tableName) + "("

	var valuesCount = 0
	args := make([]interface{}, 0)

	for i := 0; i < v.NumField(); i++ {
		row := t.Field(i).Tag.Get(helpers.RowStructTag)

		if ok, typeOf := isPK(t.Field(i)); ok {
			switch typeOf {
			case autoPK:
				continue
			case manualPK:
				val := uuid.New().String()
				if row != "" {
					if valuesCount != 0 {
						query += ", " + row
					} else {
						query += row
					}
					args = append(args, val)
					valuesCount++
				}
				continue
			}
		}

		if row != "" {
			if valuesCount != 0 {
				query += ", " + row
			} else {
				query += row
			}
			args = append(args, v.Field(i).Interface())
			valuesCount++
		}
	}

	query += ") values("
	for i := 0; i < valuesCount; i++ {
		if i < valuesCount-1 {
			query += "$" + strconv.Itoa(i+1) + ", "
		} else {
			query += "$" + strconv.Itoa(i+1)
		}
	}

	query += ")"

	return query, args
}

// QueryBuilderDelete generates normal delete queries for non nested structures
func QueryBuilderDelete(i interface{}, schemaName string, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	query := `DELETE FROM ` + fmt.Sprintf("%s.%s", schemaName, tableName) + " WHERE "

	args := make([]interface{}, 0)

	for i := 0; i < v.NumField(); i++ {

		if !checkEmpty(v.Field(i)) {
			row := t.Field(i).Tag.Get(helpers.RowStructTag)
			if row != "" {
				query += row + " = $1"
				args = append(args, v.Field(i).Interface())
				return query, args
			}
		}
	}
	return "", nil
}

// QueryBuilderUpdate generates normal update queries for non nested structures
func QueryBuilderUpdate(i interface{}, schemaName string, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	var searchBy int
	query := `UPDATE ` + fmt.Sprintf("%s.%s", schemaName, tableName) + " SET "
	args := make([]interface{}, 0)

	argsCount := 0
	for i := 0; i < v.NumField(); i++ {

		if ok, _ := isPK(t.Field(i)); ok {
			searchBy = i
			continue
		}

		row := t.Field(i).Tag.Get(helpers.RowStructTag)
		if row != "" {
			if argsCount < 1 {
				query += row + " = $" + strconv.Itoa(i+1)
			} else {
				query += " ," + row + " = $" + strconv.Itoa(i+1)
			}
			args = append(args, v.Field(i).Interface())
			argsCount++
		}

	}

	if len(args) == 0 {
		return "", nil
	}

	query += " WHERE " + t.Field(searchBy).Tag.Get(helpers.RowStructTag) + " = $" + strconv.Itoa(argsCount+1)
	args = append(args, v.Field(searchBy).Interface())

	return query, args
}

func getInnerJoin(inte interface{}, tableName string) string {
	t := reflect.TypeOf(inte)
	v := reflect.ValueOf(inte)
	var ret string
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct {
			ret += fmt.Sprintf(" INNER JOIN %s ON (%s.%s = %s.%s) ", t.Field(i).Tag.Get("fk"), tableName, t.Field(i).Tag.Get(helpers.RowStructTag), t.Field(i).Tag.Get("fk"), t.Field(i).Tag.Get("fr"))
		}
	}
	return ret
}

// QueryBuilderJoin generates get queries for nested structures with inner join support
func QueryBuilderJoin(inte interface{}, schemaName string, tableName string) (string, []interface{}) {
	query := `SELECT ` + getAllMembers(inte, schemaName+"."+tableName)

	query += " FROM " + fmt.Sprintf("%s.%s", schemaName, tableName)

	innerJoin := getInnerJoin(inte, schemaName+"."+tableName)
	query += innerJoin

	args, where := getArgsWhere(getSearchBy(inte, schemaName+"."+tableName))
	query += where

	return query, args
}

// GetIntoStruct scans rows into slice of struct
func GetIntoStruct(rows *sql.Rows, dest interface{}) {
	v := reflect.ValueOf(dest)
	direct := reflect.Indirect(v)

	if v.Kind() != reflect.Ptr {
		helpers.LogError("Destination not pointer")
		return
	}

	if direct.Kind() != reflect.Slice {
		helpers.LogError("Destination not slice")
		return
	}

	base := v.Elem().Type().Elem()

	for rows.Next() {
		ptrs := make([]interface{}, 0)
		vp := reflect.New(base)
		vpInd := vp.Elem()

		for i := 0; i < vpInd.NumField(); i++ {
			dd := reflect.Indirect(vpInd.Field(i))
			if dd.Kind() == reflect.Struct {
				for j := 0; j < dd.NumField(); j++ {
					ptrs = append(ptrs, dd.Field(j).Addr().Interface())
				}
				continue
			}
			ptrs = append(ptrs, vpInd.Field(i).Addr().Interface())
		}

		err := rows.Scan(ptrs...)
		if err != nil {
			helpers.LogError(err.Error())
		}

		direct.Set(reflect.Append(direct, reflect.Indirect(vp)))
	}
}

// isTableExist runs migrations if table is non existent
func isTableExist(schemaName string, tableName string, conn *sql.DB) {
	rows, err := conn.Query(`SELECT EXISTS (SELECT 1 FROM pg_tables WHERE schemaname = $1 AND tablename  = $2);`, schemaName, tableName)

	if err != nil {
		helpers.LogError(err.Error())
		err := database.RunMigrations()
		if err != nil {
			helpers.LogError(err.Error())
		}
		return
	}
	var exists bool
	for rows.Next() {
		err := rows.Scan(&exists)
		if err != nil {
			helpers.LogError(err.Error())
		}
	}

	if !exists {
		err := database.RunMigrations()
		if err != nil {
			helpers.LogError(err.Error())
		}
	}
}

// IsValueExists checks if value exists in table
// If it exists returns the ID of that value
func IsValueExists(conn *sql.DB, key interface{}, keyname string, tableName string) (bool, int64) {
	rows, err := conn.Query(fmt.Sprintf(`SELECT generated_id FROM %s WHERE  %s=?`, tableName, keyname), key)

	if err != nil {
		helpers.LogError(err.Error())
		return false, -1
	}

	var genID int64 = -1
	for rows.Next() {
		err := rows.Scan(&genID)
		if err != nil {
			helpers.LogError(err.Error())
		}
	}

	if genID > -1 {
		return true, genID
	}

	return false, -1
}

func checkEmpty(value reflect.Value) bool {
	// Checks int
	matchedInt, err := regexp.MatchString("int", value.Type().String())
	if err != nil {
		helpers.LogError(err.Error())
		return false
	}
	if matchedInt {
		return value.IsZero()
	}

	// else check string
	matchedString, err := regexp.MatchString("string", value.Type().String())
	if err != nil {
		helpers.LogError(err.Error())
		return false
	}
	if matchedString {
		return value.String() == ""
	}

	// else check bool
	matchedBool, err := regexp.MatchString("bool", value.Type().String())
	if err != nil {
		helpers.LogError(err.Error())
		return false
	}
	if matchedBool {
		// Bool cant be search factor
		return true
	}

	return !value.IsValid()
}

func isPK(field reflect.StructField) (bool, string) {
	if field.Tag.Get(helpers.PKStructTag) != "" {
		return true, field.Tag.Get(helpers.PKStructTag)
	}
	return false, ""
}

func getFkTable(field reflect.StructField) (bool, string) {
	if field.Tag.Get("fk") != "" {
		return true, field.Tag.Get("fk")
	}
	return false, ""
}

func getForeignRow(field reflect.StructField) string {
	if field.Tag.Get("fr") != "" {
		return field.Tag.Get("fr")
	}
	return field.Tag.Get("row")
}

// GetConn returns connection to tables
// Also check if table exists
func GetConn(schema string, table string) *sql.DB {
	conn := database.GetConn()
	isTableExist(schema, table, conn)
	return conn
}
