// Basically a massive clusterfuck of I dont know what
// Should've used hardcoded queries :''(

package models

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
	geographyPoints = "GeographyPoints"
	sortBy          = "SortBy"
)

type tmpHolder struct {
	name   string
	typeOf string
	value  reflect.Value
}

type inbuiltType interface {
	memberSearchQuery(tableName string, rowTag string) string
	memberCreateQuery(tableName string, rowTag string) string
	selectQuery(tableName string, rowTag string) tmpHolder
	isEmpty() bool
	createArgs() string
	ignoreScan() bool
}

// GeographyPoints is a struct that holds longitude, latitude of a location and optionally radius
// Radius can be specified to initiate a search of other points in radius from specified coordinates
type GeographyPoints struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
	Radius    string `scan:"ignore" json:"-"`
}

func (GeographyPoints) memberSearchQuery(tableName string, rowTag string) string {
	return fmt.Sprintf("ST_X(%s.%s::geometry), ST_Y(%s.%s::geometry)", tableName, rowTag, tableName, rowTag)
}

func (g GeographyPoints) memberCreateQuery(tableName string, rowTag string) string {
	return rowTag
}

func (g GeographyPoints) selectQuery(tableName string, rowTag string) tmpHolder {
	return tmpHolder{
		name:   fmt.Sprintf("ST_DWithin(%s.%s, ST_MakePoint(%v,%v), %v)", tableName, rowTag, g.Longitude, g.Latitude, g.Radius),
		typeOf: "onlyvalue",
		value:  reflect.ValueOf(g),
	}

}

func (g GeographyPoints) createArgs() string {
	return fmt.Sprintf("ST_SetSRID(ST_MakePoint(%v,%v), 4326)", g.Longitude, g.Latitude)
}

func (g GeographyPoints) isEmpty() bool {
	return (g.Latitude == "" || g.Longitude == "")
}

func (GeographyPoints) ignoreScan() bool {
	return false
}

type SortBy struct {
	Column string `json:"column"`
	Mode   string `json:"mode"`
}

func (SortBy) memberSearchQuery(tableName string, rowTag string) string {
	return ""
}

func (SortBy) memberCreateQuery(tableName string, rowTag string) string {
	return ""
}

func (SortBy) selectQuery(tableName string, rowTag string) tmpHolder {
	return tmpHolder{}
}

func (SortBy) createArgs() string {
	return ""
}

func (s SortBy) isEmpty() bool {
	return (s.Column == "")
}

func (SortBy) ignoreScan() bool {
	return true
}

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

func (Page) selectQuery(tableName string, rowTag string) tmpHolder {
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

func getAllMembers(inte interface{}, tableName string, isCreate bool) string {
	t := reflect.TypeOf(inte)
	v := reflect.ValueOf(inte)

	var ret = ""
	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Tag.Get(helpers.RowStructTag) != "" {
			if isCreate && checkEmpty(v.Field(i)) {
				continue
			}
			ret += ","
			if v.Field(i).Kind() == reflect.Struct {
				if isInbuiltType(v.Field(i)) {
					if isCreate {
						ret += v.Field(i).Interface().(inbuiltType).memberCreateQuery(tableName, t.Field(i).Tag.Get(helpers.RowStructTag))
						continue
					}
					ret += v.Field(i).Interface().(inbuiltType).memberSearchQuery(tableName, t.Field(i).Tag.Get(helpers.RowStructTag))
					continue
				}
				if !isCreate {
					if foreignTable := t.Field(i).Tag.Get("fk"); foreignTable != "" {
						ret += getAllMembers(v.Field(i).Interface(), foreignTable, isCreate)
						continue
					}
					ret += getAllMembers(v.Field(i).Interface(), tableName, isCreate)
					continue
				}
				for j := 0; j < v.Field(i).NumField(); j++ {
					if t.Field(i).Type.Field(j).Tag.Get("pk") != "" && !checkEmpty(v.Field(i).Field(j)) {
						ret += fmt.Sprintf("%s", t.Field(i).Tag.Get(helpers.RowStructTag))
					}
				}
			} else {
				if tableName != "" {
					ret += fmt.Sprintf("%s.%s", tableName, t.Field(i).Tag.Get(helpers.RowStructTag))
					continue
				}
				ret += fmt.Sprintf("%s", t.Field(i).Tag.Get(helpers.RowStructTag))
			}
		}
	}
	return strings.Trim(ret, ",")
}

func cleanTmpHolders(t *[]tmpHolder) {
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
					searchBy = append(searchBy, v.Field(i).Interface().(inbuiltType).selectQuery(tableName, t.Field(i).Tag.Get(helpers.RowStructTag)))
					continue
				}
				var tmpSearchBy []tmpHolder
				switch foreignTable := t.Field(i).Tag.Get("fk"); foreignTable {
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
		if val := t.Field(i).Tag.Get("type"); val != "" && !checkEmpty(v.Field(i)) {
			searchBy = append(searchBy, tmpHolder{
				name: func() string {
					if !appendTableName {
						return fmt.Sprintf("%s ", t.Field(i).Tag.Get(helpers.RowStructTag))
					}
					return fmt.Sprintf("%s.%s ", tableName, t.Field(i).Tag.Get(helpers.RowStructTag))
				}(),
				typeOf: val,
				value:  v.Field(i),
			})
		}
	}
	cleanTmpHolders(&searchBy)
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
			} else if r.typeOf == "onlyvalue" {
				query += fmt.Sprintf("%v", r.name)
			}
			if index < len(searchBy)-1 {
				query += " AND "
			}
			if r.typeOf != "onlyvalue" {
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
			_, ok1 := v.Field(i).Interface().(Page)
			helpers.LogInfo(ok1)
			if val, ok := v.Field(i).Interface().(Page); ok && !checkEmpty(v.Field(i)) {
				return fmt.Sprintf(" LIMIT %d OFFSET %d", val.Limit, val.Offset)
			}
		}
	}
	return ""
}

// QueryBuilderGet generates normal get queries for non nested structures
func QueryBuilderGet(i interface{}, tableName string) (string, []interface{}) {
	args, where := getArgsWhere(getSearchBy(i, tableName, true, false), 1)
	return fmt.Sprintf("SELECT %s FROM %s %s %s %s", getAllMembers(i, tableName, false), tableName, where, getOrderBy(i), getLimit(i)), args
}

func getStructCreateArg(v reflect.Value, t reflect.StructField) interface{} {
	for k := 0; k < v.NumField(); k++ {
		if t.Type.Field(k).Tag.Get(helpers.RowStructTag) == t.Tag.Get("fr") && !checkEmpty(v.Field(k)) {
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
		tag := t.Field(i).Tag.Get(helpers.RowStructTag)
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
					if val.Type().Field(i).Tag.Get("pk") != "" {
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

func generateUpdateQuery(query string) (ret string, length int) {
	split := strings.Split(query, ",")
	for i, e := range split {
		ret += fmt.Sprintf("%s = $%d,", e, i+1)
	}
	return strings.Trim(ret, ","), len(split)
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
			pk = t.Field(i).Tag.Get(helpers.RowStructTag)
			break
		}
	}
	args, where := getArgsWhere(getSearchBy(inte, tableName, false, false), 1)

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
			if isInbuiltType(dd) && dd.Interface().(inbuiltType).ignoreScan() {
				continue
			}
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

func isInbuiltType(v reflect.Value) bool {
	if _, ok := v.Interface().(inbuiltType); ok && v.Kind() == reflect.Struct {
		switch v.Interface().(type) {
		case GeographyPoints, SortBy, Page:
			return true
		}
	}
	return false
}

// GetConn returns connection to tables
// Also check if table exists
func GetConn(schema string, table string) *sql.DB {
	conn := database.GetConn()
	isTableExist(schema, table, conn)
	return conn
}
