package querybuilder

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database/querybuilder/types"
)

func getPtrs(dest reflect.Value, typeOf reflect.Type) []interface{} {
	ptrs := make([]interface{}, 0)
	for i := 0; i < dest.NumField(); i++ {
		dd := reflect.Indirect(dest.Field(i))
		if typeOf.Field(i).Tag.Get("scan") == "ignore" {
			continue
		}
		if dd.Kind() == reflect.Struct {
			if isInbuiltType(dd) && dd.Interface().(types.InbuiltType).IgnoreScan() {
				continue
			}
			ptrs = append(ptrs, getPtrs(dest.Field(i), typeOf.Field(i).Type)...)
			continue
		}
		ptrs = append(ptrs, dest.Field(i).Addr().Interface())
	}
	return ptrs
}

// getIntoStruct scans rows into slice of struct
func getIntoStruct(rows *sql.Rows, dest interface{}) error {
	v := reflect.ValueOf(dest)
	direct := reflect.Indirect(v)

	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("Destination not pointer")
	}
	if direct.Kind() != reflect.Slice {
		return fmt.Errorf("Destination not slice")
	}

	base := v.Elem().Type().Elem()

	for rows.Next() {
		ptrs := make([]interface{}, 0)
		vp := reflect.New(base)
		vpInd := vp.Elem()

		ptrs = append(ptrs, getPtrs(vpInd, vpInd.Type())...)

		err := rows.Scan(ptrs...)
		if err != nil {
			return err
		}

		direct.Set(reflect.Append(direct, reflect.Indirect(vp)))
	}
	return nil
}

// getIntoVar scans row into slice of single variable
func getIntoVar(rows *sql.Rows, dest interface{}) error {
	v := reflect.ValueOf(dest)
	direct := reflect.Indirect(v)

	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("Destination not pointer")
	}
	if direct.Kind() != reflect.Slice {
		return fmt.Errorf("Destination not slice")
	}

	base := v.Elem().Type().Elem()

	for rows.Next() {
		vp := reflect.New(base)
		vpInd := vp.Elem()
		err := rows.Scan(vpInd.Addr().Interface())
		if err != nil {
			return err
		}

		direct.Set(reflect.Append(direct, reflect.Indirect(vp)))
	}
	return nil
}
