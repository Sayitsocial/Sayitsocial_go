package querybuilder

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
)

const (
	RegexInt  = "int"
	RegexStr  = "string"
	RegexBool = "bool"

	ForeignTable = "ft"
	ForeignKey   = "fk"

	PrimaryKey = "pk"

	Row = "row"

	TypeOfSearch = "type"
)

type inbuiltType interface {

	// Custom query to be replaced while searching of database columns in SELECT operation
	memberSearchQuery(tableName string, rowTag string) string

	// Custom query to be replaced while searching of database columns in INSERT operation
	memberCreateQuery(tableName string, rowTag string) string

	// Returns custom holder to parse where queries
	whereQuery(tableName string, rowTag string) tmpHolder

	// Checks if struct is empty or not
	isEmpty() bool

	// Returns custom arguments for create query
	createArgs() string

	// True if struct should be ignore while scanning values from DB
	ignoreScan() bool
}

// isTableExist runs migrations if table is non existent
func isTableExist(schemaName string, tableName string, conn *sql.DB) {
	rows, err := conn.Query(`SELECT EXISTS (SELECT 1 FROM pg_tables WHERE schemaname = $1 AND tablename  = $2);`, schemaName, tableName)
	var exists bool
	if err != nil {
		for rows.Next() {
			err := rows.Scan(&exists)
			if err != nil {
				helpers.LogError(err.Error())
			}
		}
	}

	if !exists || err != nil {
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
	if isInbuiltType(value) {
		return value.Interface().(inbuiltType).isEmpty()
	}
	// Checks int
	matchedInt, err := regexp.MatchString(RegexInt, value.Type().String())
	if err != nil {
		helpers.LogError(err.Error())
		return false
	}
	if matchedInt {
		return value.IsZero()
	}

	// else check string
	matchedString, err := regexp.MatchString(RegexStr, value.Type().String())
	if err != nil {
		helpers.LogError(err.Error())
		return false
	}
	if matchedString {
		return value.String() == ""
	}

	// else check bool
	matchedBool, err := regexp.MatchString(RegexBool, value.Type().String())
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
	if field.Tag.Get(PrimaryKey) != "" {
		return true, field.Tag.Get(PrimaryKey)
	}
	return false, ""
}

func getFkTable(field reflect.StructField) (bool, string) {
	if field.Tag.Get(ForeignTable) != "" {
		return true, field.Tag.Get(ForeignTable)
	}
	return false, ""
}

func getForeignRow(field reflect.StructField) string {
	if field.Tag.Get(ForeignKey) != "" {
		return field.Tag.Get(ForeignKey)
	}
	return field.Tag.Get(Row)
}

func isInbuiltType(v reflect.Value) bool {
	_, ok := v.Interface().(inbuiltType)
	return ok && v.Kind() == reflect.Struct
}

// GetConn returns connection to tables
// Also check if table exists
func GetConn(schema string, table string) *sql.DB {
	conn := database.GetConn()
	isTableExist(schema, table, conn)
	return conn
}

func getAllMembers(inte interface{}, tableName string, isCreate bool) string {
	t := reflect.TypeOf(inte)
	v := reflect.ValueOf(inte)

	var ret = ""
	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Tag.Get(Row) != "" {
			if isCreate && checkEmpty(v.Field(i)) {
				continue
			}
			ret += ","
			if v.Field(i).Kind() == reflect.Struct {
				if isInbuiltType(v.Field(i)) {
					if isCreate {
						ret += v.Field(i).Interface().(inbuiltType).memberCreateQuery(tableName, t.Field(i).Tag.Get(Row))
						continue
					}
					ret += v.Field(i).Interface().(inbuiltType).memberSearchQuery(tableName, t.Field(i).Tag.Get(Row))
					continue
				}
				if !isCreate {
					if foreignTable := t.Field(i).Tag.Get(ForeignTable); foreignTable != "" {
						ret += getAllMembers(v.Field(i).Interface(), foreignTable, isCreate)
						continue
					}
					ret += getAllMembers(v.Field(i).Interface(), tableName, isCreate)
					continue
				}
				for j := 0; j < v.Field(i).NumField(); j++ {
					if t.Field(i).Type.Field(j).Tag.Get(PrimaryKey) != "" && !checkEmpty(v.Field(i).Field(j)) {
						ret += fmt.Sprintf("%s", t.Field(i).Tag.Get(Row))
					}
				}
			} else {
				if tableName != "" {
					ret += fmt.Sprintf("%s.%s", tableName, t.Field(i).Tag.Get(Row))
					continue
				}
				ret += fmt.Sprintf("%s", t.Field(i).Tag.Get(Row))
			}
		}
	}
	return strings.Trim(ret, ",")
}

func cleanTmpHolder(t *[]tmpHolder) {
	for i, a := range *t {
		if a.name == "" && a.typeOf == "" {
			(*t)[i] = (*t)[len(*t)-1]
			*t = (*t)[:len(*t)-1]
		}
	}
}

func getSearchBy(inte interface{}, tableName string, appendTableName bool, forcePK bool) (searchBy []tmpHolder) {
	t := reflect.TypeOf(inte)
	v := reflect.ValueOf(inte)
	for i := 0; i < v.NumField(); i++ {
		if forcePK {
			if pk, _ := isPK(t.Field(i)); !pk || checkEmpty(v.Field(i)) {
				continue
			}
		}
		if v.Field(i).Kind() == reflect.Struct {
			if appendTableName {
				if isInbuiltType(v.Field(i)) && !checkEmpty(v.Field(i)) {
					searchBy = append(searchBy, v.Field(i).Interface().(inbuiltType).whereQuery(tableName, t.Field(i).Tag.Get(Row)))
					continue
				}
				var tmpSearchBy []tmpHolder
				switch foreignTable := t.Field(i).Tag.Get(ForeignTable); foreignTable {
				case "":
					tmpSearchBy = getSearchBy(v.Field(i).Interface(), foreignTable, appendTableName, forcePK)
					break
				default:
					tmpSearchBy = getSearchBy(v.Field(i).Interface(), tableName, appendTableName, forcePK)
					break
				}
				searchBy = append(searchBy, tmpSearchBy...)

				continue
			}
		}
		if val := t.Field(i).Tag.Get(TypeOfSearch); val != "" && !checkEmpty(v.Field(i)) {
			searchBy = append(searchBy, tmpHolder{
				name: func() string {
					if !appendTableName {
						return fmt.Sprintf("%s ", t.Field(i).Tag.Get(Row))
					}
					return fmt.Sprintf("%s.%s ", tableName, t.Field(i).Tag.Get(Row))
				}(),
				typeOf: val,
				value:  v.Field(i),
			})
		}
	}
	cleanTmpHolder(&searchBy)
	return
}

func getArgsWhere(searchBy []tmpHolder, index int) ([]interface{}, string) {
	args := make([]interface{}, 0)
	if len(searchBy) == 0 {
		return nil, ""
	}
	query := " WHERE "
	for _, r := range searchBy {
		if r.name != "" && r.name != " " {
			if r.typeOf == "exact" {
				query += fmt.Sprintf("%v = $%d", r.name, index+1)
			} else if r.typeOf == "like" {
				query += fmt.Sprintf("%v ILIKE $%d", r.name, index+1)
			} else if r.typeOf == "onlyname" {
				query += fmt.Sprintf("%v", r.name)
			}
			if index < len(searchBy)-1 {
				query += " AND "
			}
			if r.typeOf != "onlyname" {
				args = append(args, r.value.Interface())
			}
		}
	}
	return args, strings.TrimSuffix(query, "AND ")
}

func getOrderBy(inte interface{}) (orderQuery string) {
	v := reflect.ValueOf(inte)

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct && isInbuiltType(v.Field(i)) {
			if val, ok := v.Field(i).Interface().(SortBy); ok && !checkEmpty(v.Field(i)) {
				return fmt.Sprintf(" ORDER BY %s %s", val.Column, val.Mode)
			}
		}
	}
	return ""
}

func getLimit(inte interface{}) (limitQuery string) {
	v := reflect.ValueOf(inte)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct && isInbuiltType(v.Field(i)) {
			if val, ok := v.Field(i).Interface().(Page); ok && !checkEmpty(v.Field(i)) {
				return fmt.Sprintf(" LIMIT %d OFFSET %d", val.Limit, val.Offset)
			}
		}
	}
	return ""
}

func getStructCreateArg(v reflect.Value, t reflect.StructField) interface{} {
	for k := 0; k < v.NumField(); k++ {
		if t.Type.Field(k).Tag.Get(Row) == t.Tag.Get("fk") && !checkEmpty(v.Field(k)) {
			return v.Field(k).Interface()
		}
	}
	return sql.NullString{}
}

func getFieldsByTagMap(inte interface{}) (fieldsByTag map[string]reflect.Value) {
	fieldsByTag = make(map[string]reflect.Value)
	v := reflect.ValueOf(inte)
	t := reflect.TypeOf(inte)

	for i := 0; i < v.NumField(); i++ {
		tag := t.Field(i).Tag.Get(Row)
		if tag != "" {
			fieldsByTag[tag] = v.Field(i)
		}
	}
	return
}

func getValuesCount(inte interface{}, index *int, message string) (ret string, args []interface{}) {
	tagMap := getFieldsByTagMap(inte)
	values := strings.Split(message, ",")
	for _, i := range values {
		if val, ok := tagMap[i]; ok {
			if isInbuiltType(val) && !checkEmpty(val) {
				ret += fmt.Sprintf("%s,", val.Interface().(inbuiltType).createArgs())
				continue
			}
			if val.Kind() == reflect.Struct {
				for i := 0; i < val.NumField(); i++ {
					if val.Type().Field(i).Tag.Get(PrimaryKey) != "" {
						ret += fmt.Sprintf("$%d,", *index)
						args = append(args, val.Field(i).Interface())
						*index++
						break
					}
				}
				continue
			}
			ret += fmt.Sprintf("$%d,", *index)
			args = append(args, val.Interface())
			*index++
		}
	}
	return strings.Trim(ret, ","), args
}

func generateUpdateQuery(query string) (ret string, length int) {
	split := strings.Split(query, ",")
	for i, e := range split {
		ret += fmt.Sprintf("%s = $%d,", e, i+1)
	}
	return strings.Trim(ret, ","), len(split)
}

func getInnerJoin(inte interface{}, tableName string) string {
	t := reflect.TypeOf(inte)
	v := reflect.ValueOf(inte)
	var ret string
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct && t.Field(i).Tag.Get(ForeignTable) != "" {
			ret += fmt.Sprintf(" INNER JOIN %s ON (%s.%s = %s.%s) ", t.Field(i).Tag.Get(ForeignTable), tableName, t.Field(i).Tag.Get(Row), t.Field(i).Tag.Get(ForeignTable), t.Field(i).Tag.Get(ForeignKey))
			ret += getInnerJoin(v.Field(i).Interface(), fmt.Sprintf("%s", t.Field(i).Tag.Get(ForeignTable)))
		}
	}
	return ret
}
