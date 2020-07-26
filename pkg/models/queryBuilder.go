package models

import (
	"database/sql"
	"fmt"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/database"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"reflect"
	"regexp"
	"strconv"
)

const component = "QueryBuilder"

func QueryBuilderGet(i interface{}, schemaName string, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	query := `SELECT `

	var searchByRow string
	var searchByIndex int
	for i := 0; i < v.NumField(); i++ {
		row := t.Field(i).Tag.Get(helpers.RowStructTag)

		if !checkEmpty(v.Field(i)) {
			searchByRow = row
			searchByIndex = i
		}

		if row != "" {
			if i < t.NumField()-1 {
				query += row + ", "
			} else {
				query += row
			}
		}
	}

	query += " FROM " + fmt.Sprintf("%s.%s", schemaName, tableName)
	if searchByRow == "" {
		return query, nil
	}

	if t.Field(searchByIndex).Tag.Get("type") == "exact" {
		query += " WHERE " + searchByRow + " = $1"
	} else if t.Field(searchByIndex).Tag.Get("type") == "like" {
		query += " WHERE " + searchByRow + " LIKE $1 COLLATE NOCASE"
	}
	args := []interface{}{v.Field(searchByIndex).Interface()}

	return query, args
}

func QueryBuilderCreate(i interface{}, schemaName string, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	query := `INSERT INTO ` + fmt.Sprintf("%s.%s", schemaName, tableName) + "("

	var valuesCount = 0
	args := make([]interface{}, 0)

	for i := 0; i < v.NumField(); i++ {
		row := t.Field(i).Tag.Get(helpers.RowStructTag)

		if isPK(t.Field(i)) {
			continue
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

func QueryBuilderUpdate(i interface{}, schemaName string, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	var searchBy int
	query := `UPDATE ` + fmt.Sprintf("%s.%s", schemaName, tableName) + " SET "
	args := make([]interface{}, 0)

	argsCount := 0
	for i := 0; i < v.NumField(); i++ {

		if isPK(t.Field(i)) {
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

	if len(args) < 0 {
		return "", nil
	}

	query += " WHERE " + t.Field(searchBy).Tag.Get(helpers.RowStructTag) + " = $" + strconv.Itoa(argsCount+1)
	args = append(args, v.Field(searchBy).Interface())

	return query, args
}

func GetIntoStruct(rows *sql.Rows, dest interface{}) {
	v := reflect.ValueOf(dest)
	direct := reflect.Indirect(v)

	if v.Kind() != reflect.Ptr {
		helpers.LogError("Destination not pointer", component)
		return
	}

	if direct.Kind() != reflect.Slice {
		helpers.LogError("Destination not slice", component)
		return
	}

	base := v.Elem().Type().Elem()
	vp := reflect.New(base)

	for rows.Next() {
		direct.Set(reflect.Append(direct, scanSingleStruct(vp, rows)))
	}
}

func scanSingleStruct(dest reflect.Value, row *sql.Rows) reflect.Value {
	numfields := reflect.Indirect(dest).NumField()
	ind := reflect.Indirect(dest)

	ptrs := make([]interface{}, numfields)

	for i := 0; i < numfields; i++ {
		ptrs[i] = ind.Field(i).Addr().Interface()
	}

	err := row.Scan(ptrs...)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
	return ind
}

func IsTableEmpty(schemaName string, tableName string, conn *sql.DB) {
	rows, err := conn.Query(`SELECT EXISTS (SELECT 1 FROM pg_tables WHERE schemaname = $1 AND tablename  = $2);`, schemaName, tableName)

	if err != nil {
		helpers.LogError(err.Error(), component)
		err := database.RunMigrations()
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		return
	}
	var exists bool
	for rows.Next() {
		err := rows.Scan(&exists)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
	}

	if !exists {
		err := database.RunMigrations()
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
	}
}

func IsValueExists(conn *sql.DB, key interface{}, keyname string, tableName string) (bool, int64) {
	rows, err := conn.Query(fmt.Sprintf(`SELECT generated_id FROM %s WHERE  %s=?`, tableName, keyname), key)

	if err != nil {
		helpers.LogError(err.Error(), component)
		return false, -1
	}

	var genId int64 = -1
	for rows.Next() {
		err := rows.Scan(&genId)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
	}

	if genId > -1 {
		return true, genId
	}

	return false, -1
}

func checkEmpty(value reflect.Value) bool {
	// Checks int
	matchedInt, err := regexp.MatchString("int", value.Type().String())
	if err != nil {
		helpers.LogError(err.Error(), component)
		return false
	}
	if matchedInt {
		return value.IsZero()
	}

	//else check string
	matchedString, err := regexp.MatchString("string", value.Type().String())
	if err != nil {
		helpers.LogError(err.Error(), component)
		return false
	}
	if matchedString {
		return value.String() == ""
	}

	//else check bool
	matchedBool, err := regexp.MatchString("bool", value.Type().String())
	if err != nil {
		helpers.LogError(err.Error(), component)
		return false
	}
	if matchedBool {
		// Bool cant be search factor
		return true
	}

	return !value.IsValid()
}

func isPK(field reflect.StructField) bool {
	return field.Tag.Get(helpers.PKStructTag) == "auto"
}

func GetConn(schema string, table string) *sql.DB {
	conn := database.GetConn()
	IsTableEmpty(schema, table, conn)
	return conn
}
