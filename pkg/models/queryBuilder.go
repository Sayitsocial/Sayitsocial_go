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
)

type tmpHolder struct {
	name   string
	typeOf string
	value  reflect.Value
}

// GeographyPoints is a struct that holds longitude, latitude of a location and optionally radius
// Radius can be specified to initiate a search of other points in radius from specified coordinates
type GeographyPoints struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Radius    float64 `scan:"ignore" json:"-"`
}

func getAllMembers(inte interface{}, tableName string, isOuter bool) string {
	t := reflect.TypeOf(inte)
	v := reflect.ValueOf(inte)

	var ret = ""
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct {
			if _, ok := v.Field(i).Interface().(GeographyPoints); ok {
				ret += fmt.Sprintf("ST_X(%s.%s::geometry), ST_Y(%s.%s::geometry), ", tableName, t.Field(i).Tag.Get(helpers.RowStructTag), tableName, t.Field(i).Tag.Get(helpers.RowStructTag))
			} else if foreignTable := t.Field(i).Tag.Get("fk"); foreignTable != "" {
				ret += getAllMembers(v.Field(i).Interface(), foreignTable, false)
			} else {
				ret += getAllMembers(v.Field(i).Interface(), tableName, false)
			}
			continue
		} else {
			ret += fmt.Sprintf("%s.%s ", tableName, t.Field(i).Tag.Get(helpers.RowStructTag))
		}
		ret += ","
	}
	return func() string {
		if isOuter {
			return ret[:len(ret)-1]
		}
		return ret
	}()
}

func getSearchBy(inte interface{}, tableName string, appendTableName bool) (totalSearchByCount int, searchBy []tmpHolder) {
	t := reflect.TypeOf(inte)
	v := reflect.ValueOf(inte)
	for i := 0; i < v.NumField(); i++ {
		if _, ok := v.Field(i).Interface().(GeographyPoints); !ok && v.Field(i).Kind() == reflect.Struct {
			if appendTableName {
				if foreignTable := t.Field(i).Tag.Get("fk"); foreignTable != "" {
					tmpCount, tmpSearchBy := getSearchBy(v.Field(i).Interface(), foreignTable, true)
					totalSearchByCount += tmpCount
					searchBy = append(searchBy, tmpSearchBy...)
				} else {
					tmpCount, tmpSearchBy := getSearchBy(v.Field(i).Interface(), tableName, true)
					totalSearchByCount += tmpCount
					searchBy = append(searchBy, tmpSearchBy...)
				}
				continue
			}
		}
		if !checkEmpty(v.Field(i)) {
			if _, ok := v.Field(i).Interface().(GeographyPoints); ok {
				searchBy = append(searchBy, tmpHolder{
					name: func() string {
						if appendTableName {
							return fmt.Sprintf("ST_DWithin(%s.%s, ST_MakePoint(%v,%v), %v)", tableName, t.Field(i).Tag.Get(helpers.RowStructTag), v.Field(i).FieldByName("Longitude").Interface(), v.Field(i).FieldByName("Latitude").Interface(), v.Field(i).FieldByName("Radius").Interface())
						}
						return fmt.Sprintf("ST_DWithin(%s, ST_MakePoint(%v,%v), %v)", t.Field(i).Tag.Get(helpers.RowStructTag), v.Field(i).FieldByName("Longitude").Interface(), v.Field(i).FieldByName("Latitude").Interface(), v.Field(i).FieldByName("Radius").Interface())
					}(),
					typeOf: "onlyvalue",
					value:  reflect.ValueOf(nil),
				})
				continue
			}

			if val := t.Field(i).Tag.Get("type"); val != "" {
				searchBy = append(searchBy, tmpHolder{
					name: func() string {
						if !appendTableName {
							return fmt.Sprintf("%s ", t.Field(i).Tag.Get(helpers.RowStructTag))
						}
						return fmt.Sprintf("%s.%s ", tableName, t.Field(i).Tag.Get(helpers.RowStructTag))
					}(),
					typeOf: val,
					value: func() reflect.Value {
						if v.Field(i).Kind() == reflect.Struct {
							tmpCount, tmpSearchBy := getSearchBy(v.Field(i).Interface(), "", false)
							if tmpCount > 0 {
								if tmpSearchBy[0].value.Kind() != reflect.Struct {
									return tmpSearchBy[0].value
								}
							}
							return reflect.ValueOf(nil)

						}
						return v.Field(i)
					}(),
				})
				totalSearchByCount++
			}
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
	helpers.LogInfo(searchBy)
	for searchByCount, r := range searchBy {
		if r.typeOf == "exact" {
			query += fmt.Sprintf("%v = $%d", r.name, searchByCount+1)
		} else if r.typeOf == "like" {
			query += fmt.Sprintf("%v ILIKE $%d", r.name, searchByCount+1)
		} else if r.typeOf == "onlyvalue" {
			query += fmt.Sprintf("%v", r.name)
		}
		if searchByCount < totalSearchByCount-1 {
			query += " AND "
		}
		if r.typeOf != "onlyvalue" {
			args = append(args, r.value.Interface())
		}
	}
	return args, query
}

// QueryBuilderGet generates normal get queries for non nested structures
func QueryBuilderGet(i interface{}, tableName string) (string, []interface{}) {
	query := `SELECT ` + getAllMembers(i, tableName, true)

	args, where := getArgsWhere(getSearchBy(i, tableName, true))
	query += " FROM " + tableName
	query += where

	return query, args
}

func getSructCreateArg(v reflect.Value, t reflect.StructField) interface{} {
	for k := 0; k < v.NumField(); k++ {
		if t.Type.Field(k).Tag.Get(helpers.RowStructTag) == t.Tag.Get("fr") && !checkEmpty(v.Field(k)) {
			return v.Field(k).Interface()
		}
	}
	return sql.NullString{}
}

// QueryBuilderCreate generates normal create queries for non nested structures
func QueryBuilderCreate(i interface{}, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	query := `INSERT INTO ` + tableName + "("

	var valuesCount = 0
	args := make([]interface{}, 0)

	// pqsl driver interprets function names passed in args as string and errors out
	// TODO: Eliminate args and put everything in query
	geographyValues := make(map[int]bool, 0)

	for i := 0; i < v.NumField(); i++ {
		row := t.Field(i).Tag.Get(helpers.RowStructTag)

		if checkEmpty(v.Field(i)) {
			if ok, _ := isPK(t.Field(i)); ok {
				continue
			}
		}

		if row != "" {
			_, isGeographyPoints := v.Field(i).Interface().(GeographyPoints)
			if isGeographyPoints {
				geographyValues[i] = true
			}

			if valuesCount != 0 {
				query += ", " + row
			} else {
				query += row
			}

			if !isGeographyPoints {
				if v.Field(i).Kind() == reflect.Struct {
					args = append(args, getSructCreateArg(v.Field(i), t.Field(i)))
				} else {
					args = append(args, v.Field(i).Interface())
				}
			}
			valuesCount++
		}
	}

	numbering := 0
	query += ") values("
	for i := 0; i < valuesCount; i++ {
		if _, ok := geographyValues[i]; ok {
			query += fmt.Sprintf("ST_SetSRID(ST_MakePoint(%v, %v), 4326)", v.Field(i).FieldByName("Longitude").Interface(), v.Field(i).FieldByName("Latitude").Interface())
		} else {
			numbering++
			query += "$" + strconv.Itoa(numbering)
		}

		if i < valuesCount-1 {
			query += ", "
		}
	}

	query += ")"

	return query, args
}

// QueryBuilderDelete generates normal delete queries for non nested structures
func QueryBuilderDelete(i interface{}, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	query := "DELETE FROM " + tableName + " WHERE "

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
func QueryBuilderUpdate(i interface{}, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	var searchBy int
	query := `UPDATE ` + tableName + " SET "
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
		if v.Field(i).Kind() == reflect.Struct && t.Field(i).Tag.Get("fk") != "" {
			ret += fmt.Sprintf(" INNER JOIN %s ON (%s.%s = %s.%s) ", t.Field(i).Tag.Get("fk"), tableName, t.Field(i).Tag.Get(helpers.RowStructTag), t.Field(i).Tag.Get("fk"), t.Field(i).Tag.Get("fr"))
			ret += getInnerJoin(v.Field(i).Interface(), fmt.Sprintf("%s", t.Field(i).Tag.Get("fk")))
		}
	}
	return ret
}

// QueryBuilderJoin generates get queries for nested structures with inner join support
func QueryBuilderJoin(inte interface{}, tableName string) (string, []interface{}) {
	args, where := getArgsWhere(getSearchBy(inte, tableName, true))
	query := fmt.Sprintf("SELECT %s FROM %s %s %s", getAllMembers(inte, tableName, true), tableName, getInnerJoin(inte, tableName), where)
	return query, args
}

// QueryBuilderCount generates count queries for primary key in structure
func QueryBuilderCount(inte interface{}, tableName string) (string, []interface{}) {
	t := reflect.TypeOf(inte)
	v := reflect.ValueOf(inte)

	var pk string = "*"
	for i := 0; i < v.NumField(); i++ {
		if ok, _ := isPK(t.Field(i)); ok {
			pk = t.Field(i).Tag.Get("row")
			break
		}
	}
	args, where := getArgsWhere(getSearchBy(inte, tableName, false))

	return fmt.Sprintf("SELECT COUNT(%s) FROM %s %s", pk, tableName, where), args
}

func getPtrs(dest reflect.Value, typeOf reflect.Type) []interface{} {
	ptrs := make([]interface{}, 0)
	for i := 0; i < dest.NumField(); i++ {
		dd := reflect.Indirect(dest.Field(i))
		if typeOf.Field(i).Tag.Get("scan") == "ignore" {
			continue
		}
		if dd.Kind() == reflect.Struct {
			ptrs = append(ptrs, getPtrs(dest.Field(i), typeOf.Field(i).Type)...)
			continue
		}
		ptrs = append(ptrs, dest.Field(i).Addr().Interface())
	}
	return ptrs
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

		ptrs = append(ptrs, getPtrs(vpInd, vpInd.Type())...)

		err := rows.Scan(ptrs...)
		if err != nil {
			helpers.LogError(err.Error())
		}

		direct.Set(reflect.Append(direct, reflect.Indirect(vp)))
	}
}

// GetIntoVar scans row into slice of single variable
func GetIntoVar(rows *sql.Rows, dest interface{}) {
	v := reflect.ValueOf(dest)
	direct := reflect.Indirect(v)
	base := v.Elem().Type().Elem()

	if v.Kind() != reflect.Ptr {
		helpers.LogError("Destination not pointer")
		return
	}

	for rows.Next() {
		vp := reflect.New(base)
		vpInd := vp.Elem()
		err := rows.Scan(vpInd.Addr().Interface())
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
