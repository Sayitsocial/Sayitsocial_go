package models

import (
	"database/sql"
	"fmt"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/database"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
	"github.com/google/uuid"
	"reflect"
	"regexp"
	"strconv"
)

const (
	AutoPK   = "auto"
	ManualPK = "manual"
)

func QueryBuilderGet(i interface{}, schemaName string, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	query := `SELECT `

	var searchBy = make(map[int]string)
	for i := 0; i < v.NumField(); i++ {
		row := t.Field(i).Tag.Get(helpers.RowStructTag)

		if !checkEmpty(v.Field(i)) {
			searchBy[i] = row
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
	if len(searchBy) == 0 {
		return query, nil
	}

	args := make([]interface{}, 0)
	query += " WHERE "
	var searchByCount int
	for index, rowName := range searchBy {
		searchByCount += 1
		if t.Field(index).Tag.Get("type") == "exact" {
			query += fmt.Sprintf("%v = $%s", rowName, strconv.Itoa(searchByCount))
		} else if t.Field(index).Tag.Get("type") == "like" {
			query += fmt.Sprintf("%v ILIKE $%s", rowName, strconv.Itoa(searchByCount))
		}
		if searchByCount != len(searchBy) {
			query += " AND "
		}
		args = append(args, v.Field(index).Interface())
	}

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

		if ok, typeOf := isPK(t.Field(i)); ok {
			switch typeOf {
			case AutoPK:
				continue
			case ManualPK:
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
		helpers.LogError("Destination not pointer")
		return
	}

	if direct.Kind() != reflect.Slice {
		helpers.LogError("Destination not slice")
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
		helpers.LogError(err.Error())
	}
	return ind
}

func IsTableEmpty(schemaName string, tableName string, conn *sql.DB) {
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

func IsValueExists(conn *sql.DB, key interface{}, keyname string, tableName string) (bool, int64) {
	rows, err := conn.Query(fmt.Sprintf(`SELECT generated_id FROM %s WHERE  %s=?`, tableName, keyname), key)

	if err != nil {
		helpers.LogError(err.Error())
		return false, -1
	}

	var genId int64 = -1
	for rows.Next() {
		err := rows.Scan(&genId)
		if err != nil {
			helpers.LogError(err.Error())
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
		helpers.LogError(err.Error())
		return false
	}
	if matchedInt {
		return value.IsZero()
	}

	//else check string
	matchedString, err := regexp.MatchString("string", value.Type().String())
	if err != nil {
		helpers.LogError(err.Error())
		return false
	}
	if matchedString {
		return value.String() == ""
	}

	//else check bool
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
	} else {
		return false, ""
	}
}

func GetConn(schema string, table string) *sql.DB {
	conn := database.GetConn()
	IsTableEmpty(schema, table, conn)
	return conn
}
