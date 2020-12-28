package querybuilder

import (
	"testing"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database/querybuilder/types"
)

type BoolTest struct {
	Column1 string         `sorm:"column1,pk_autoinc"`
	Column2 types.SormBool `sorm:"column2"`
}

func (BoolTest) GetTableName() (string, string) {
	return "public", "test_model"
}

func Test_BoolGet(t *testing.T) {
	passQuery := `SELECT public.test_model.column1,public.test_model.column2 FROM public.test_model`
	err := queryBuilderCommon(BoolTest{
		Column2: types.SormBool{
			IsValueEmpty: true,
		},
	}, passQuery, Config{}, func(args []interface{}) bool {
		return len(args) == 0
	}, queryBuilderJoin)
	if err != nil {
		t.Error(err)
	}
}

func Test_BoolGetWithArgs(t *testing.T) {
	passQuery := `SELECT public.test_model.column1,public.test_model.column2 FROM public.test_model  WHERE public.test_model.column2=$1`
	err := queryBuilderCommon(BoolTest{
		Column2: types.SormBool{
			Value: false,
		},
	}, passQuery, Config{}, func(args []interface{}) bool {
		return len(args) == 1 && args[0].(bool) == false
	}, queryBuilderJoin)
	if err != nil {
		t.Error(err)
	}
}

func Test_BoolCreate(t *testing.T) {
	passQuery := `INSERT INTO public.test_model (column1,column2) VALUES ($1,$2)`
	err := queryBuilderCommon(BoolTest{
		Column1: "test",
		Column2: types.SormBool{
			Value: true,
		},
	}, passQuery, Config{}, func(args []interface{}) bool {
		return len(args) == 2 && args[1].(bool) == true
	}, queryBuilderCreate)
	if err != nil {
		t.Error(err)
	}
}

func Test_BoolUpdate(t *testing.T) {
	passQuery := `UPDATE public.test_model SET column1=$1,column2=$2 WHERE column1=$3`
	err := queryBuilderCommon(BoolTest{
		Column1: "test",
		Column2: types.SormBool{
			Value: false,
		},
	}, passQuery, Config{}, func(args []interface{}) bool {
		return (len(args) == 3 || (args[0] == "test" && args[1].(bool) == false) && args[2] == "test")
	}, queryBuilderUpdate)
	if err != nil {
		t.Error(err)
	}
}

func Test_BoolDelete(t *testing.T) {
	passQuery := `DELETE FROM public.test_model WHERE column2=$1`
	err := queryBuilderCommon(BoolTest{
		Column2: types.SormBool{
			Value: true,
		},
	}, passQuery, Config{}, func(args []interface{}) bool {
		if val, ok := args[0].(types.SormBool); ok {
			return (val.Value == true)
		}
		return false
	}, queryBuilderDelete)
	if err != nil {
		t.Error(err)
	}
}

func Test_BoolCount(t *testing.T) {
	passQuery := `SELECT COUNT(column1) FROM public.test_model WHERE column2=$1`
	err := queryBuilderCommon(BoolTest{
		Column2: types.SormBool{
			Value: false,
		},
	}, passQuery, Config{}, func(args []interface{}) bool {
		if val, ok := args[0].(types.SormBool); ok {
			return (val.Value == false)
		}
		return false
	}, queryBuilderCount)
	if err != nil {
		t.Error(err)
	}
}
