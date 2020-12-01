package querybuilder

import (
	"database/sql"
	"reflect"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
)

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
