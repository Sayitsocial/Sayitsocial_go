package querybuilder

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database"
	"github.com/Sayitsocial/Sayitsocial_go/pkg/database/querybuilder/types"
)

const (
	regexInt  = "int"
	regexStr  = "string"
	regexBool = "bool"

	ForeignTable = "ft"
	ForeignKey   = "fk"

	PrimaryKey = "pk"

	Row = "row"

	TypeOfSearch = "type"

	indexPlaceholder = "`i"
)

type valueType int

const inbuiltValueType valueType = 0
const normalValueType valueType = 1

type colHolder struct {
	col       string
	table     string
	valueType valueType
	value     interface{}
	primary   string
	isForeign bool
}

type foreignHolder struct {
	col          string
	table        string
	key          string
	foreignTable string
	value        interface{}
}

func generateColHolder(i interface{}, tableName string, isForeign bool) (holder []colHolder, foreign []foreignHolder) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	for i := 0; i < v.NumField(); i++ {
		if col := t.Field(i).Tag.Get(Row); col != "" {
			if v.Field(i).Kind() == reflect.Struct {
				if isInbuiltType(v.Field(i)) {
					holder = append(holder, colHolder{
						col:       col,
						table:     tableName,
						valueType: inbuiltValueType,
						value:     v.Field(i).Interface(),
						primary:   t.Field(i).Tag.Get(PrimaryKey),
						isForeign: isForeign,
					})
					continue
				}

				tName := tableName
				foreignTable := t.Field(i).Tag.Get(ForeignTable)
				foreignKey := t.Field(i).Tag.Get(ForeignKey)

				if foreignTable != "" {
					tName = foreignTable
				}
				tHolder, tfHolder := generateColHolder(v.Field(i).Interface(), tName, func() bool {
					return tName != tableName
				}())

				if foreignTable != "" && foreignKey != "" {
					for _, e := range tHolder {
						if e.col == foreignKey {
							foreign = append(foreign, foreignHolder{
								col:          col,
								table:        tableName,
								key:          foreignKey,
								value:        e.value,
								foreignTable: foreignTable,
							})
							break
						}
					}
				}
				holder = append(holder, tHolder...)
				foreign = append(foreign, tfHolder...)
				continue
			}
			holder = append(holder, colHolder{
				col:       col,
				table:     tableName,
				valueType: normalValueType,
				value:     v.Field(i).Interface(),
				primary:   t.Field(i).Tag.Get(PrimaryKey),
				isForeign: isForeign,
			})
		}
	}
	return
}

func checkEmpty(i interface{}) bool {
	value := reflect.ValueOf(i)
	if isInbuiltType(value) {
		return i.(types.InbuiltType).IsEmpty()
	}

	if value.Kind() == reflect.Struct {
		for i := 0; i < value.NumField(); i++ {
			if !checkEmpty(value.Field(i)) {
				return false
			}
		}
		return true
	}

	// Checks int
	matchedInt, err := regexp.MatchString(regexInt, value.Type().String())
	if err != nil {
		return false
	}
	if matchedInt {
		return value.IsZero()
	}

	// else check string
	matchedString, err := regexp.MatchString(regexStr, value.Type().String())
	if err != nil {
		return false
	}
	if matchedString {
		return value.String() == ""
	}

	// else check bool
	matchedBool, err := regexp.MatchString(regexBool, value.Type().String())
	if err != nil {
		return false
	}
	if matchedBool {
		// Bool cant be search factor
		return true
	}

	return !value.IsValid()
}

func isPK(pk string) bool {
	return pk != ""
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
	_, ok := v.Interface().(types.InbuiltType)
	return ok
}

// GetConn returns connection to tables
func GetConn() *sql.DB {
	return database.GetConn()
}

func getAllMembers(cols []colHolder, foreignCols []foreignHolder, isCreate bool) (string, []int) {
	indices := make([]int, 0)
	ret := ""
	for i, col := range cols {
		if isCreate && checkEmpty(col.value) {
			continue
		}

		switch isCreate {
		case true:
			if col.valueType == inbuiltValueType {
				ret = fmt.Sprintf("%s,%s", ret, col.value.(types.InbuiltType).CreateQuery(col.table, col.col))
				indices = append(indices, i)
				continue
			}
			if col.isForeign {
				continue
			}
			ret = fmt.Sprintf("%s,%s", ret, col.col)
		case false:
			if col.valueType == inbuiltValueType {
				ret = fmt.Sprintf("%s,%s", ret, col.value.(types.InbuiltType).SearchQuery(col.table, col.col))
				continue
			}
			ret = fmt.Sprintf("%s,%s.%s", ret, col.table, col.col)
		}

		if isCreate {
			indices = append(indices, i)
		}
	}

	if isCreate {
		for _, foreignCol := range foreignCols {
			ret = fmt.Sprintf("%s,%s", ret, foreignCol.col)
		}
	}
	return strings.Trim(ret, ","), indices
}

func getWhere(cols []colHolder, foreignCols []foreignHolder, appendTableName bool, forcePK bool) (query string, args []interface{}) {
	for _, col := range cols {
		if !checkEmpty(col.value) && !col.isForeign {
			if forcePK && !isPK(col.primary) {
				continue
			}
			if appendTableName {
				if col.valueType == inbuiltValueType {
					tquery, targs := col.value.(types.InbuiltType).WhereQuery(col.table, col.col)
					if tquery != "" {
						query = fmt.Sprintf("%s %s AND", query, tquery)
					}
					args = append(args, targs...)
					continue
				}
				query = fmt.Sprintf("%s %s.%s=$%s AND", query, col.table, col.col, indexPlaceholder)
				args = append(args, col.value)
				continue
			}
			query = fmt.Sprintf("%s %s=$%s AND", query, col.col, indexPlaceholder)
			args = append(args, col.value)
		}
	}

	// for _, foreignCol := range foreignCols {
	// 	query = fmt.Sprintf("%s")
	// }

	query = strings.Trim(strings.Trim(query, "AND"), " ")
	if query != "" {
		return "WHERE " + query, args
	}
	return "", nil
}

func getInsertValues(indices []int, colHolder []colHolder, foreign []foreignHolder) (string, []interface{}) {
	query := ""
	args := make([]interface{}, 0)
	for _, col := range indices {
		if c := colHolder[col]; c.valueType == inbuiltValueType {
			q, a := c.value.(types.InbuiltType).CreateArgs(indexPlaceholder)
			query = fmt.Sprintf("%s,%s", query, q)
			args = append(args, a...)
		} else {
			query = fmt.Sprintf("%s,$%s", query, indexPlaceholder)
			args = append(args, c.value)
		}
	}

	for _, foreignCol := range foreign {
		query = fmt.Sprintf("%s,$%s", query, indexPlaceholder)
		args = append(args, foreignCol.value)
	}
	return strings.Trim(query, ","), args
}

func generateUpdateQuery(cols []colHolder, foreign []foreignHolder) (ret string, args []interface{}) {
	for _, col := range cols {
		if !col.isForeign && !checkEmpty(col.value) {
			if col.valueType == inbuiltValueType {
				q, a := col.value.(types.InbuiltType).CreateArgs(indexPlaceholder)
				ret = fmt.Sprintf("%s,%s=%s", ret, col.col, q)
				args = append(args, a...)
				continue
			}
			ret = fmt.Sprintf("%s,%s=$%s", ret, col.col, indexPlaceholder)
			args = append(args, col.value)
		}
	}

	for _, foreignCol := range foreign {
		if !checkEmpty(foreignCol.value) {
			ret = fmt.Sprintf("%s,%s=$%s", ret, foreignCol.col, indexPlaceholder)
			args = append(args, foreignCol.value)
		}
	}
	return strings.Trim(ret, ","), args
}

func getInnerJoin(foreign []foreignHolder) (ret string) {
	for _, foreignCols := range foreign {
		ret = fmt.Sprintf("INNER JOIN %s ON (%s.%s=%s.%s)", foreignCols.foreignTable, foreignCols.table, foreignCols.col, foreignCols.foreignTable, foreignCols.key)
	}
	return
}

func indexifyQuery(s string) string {
	exp := regexp.MustCompile(indexPlaceholder)
	index := 0
	return string(exp.ReplaceAllFunc([]byte(s), func(s []byte) []byte {
		index++
		return []byte(strconv.Itoa(index))
	}))
}
