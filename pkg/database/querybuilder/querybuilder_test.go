package querybuilder

import (
	"fmt"
	"strings"
	"testing"
)

type Dummy struct {
	Column5 string `row:"column5" type:"exact" pk:"auto"`
	Column6 string `row:"column6" type:"like"`
}

func (Dummy) GetTableName() (string, string) {
	return "public", "test_model"
}

type TestModel struct {
	Column1 string `row:"column1" type:"exact" pk:"manual"`
	Column2 string `row:"column2" type:"like"`
	Column3 int64  `row:"column3" type:"exact"`
	Column4 bool   `row:"column4" type:"exact"`
	Column7 Dummy  `row:"column5" type:"exact" ft:"public.dummy" fk:"column5"`
}

func (TestModel) GetTableName() (string, string) {
	return "public", "test_model"
}

func queryBuilderCommon(model Model, passQuery string, passArgs func([]interface{}) bool, method func(interface{}, string, string) (string, []interface{})) error {
	schema, table := model.GetTableName()
	query, args := method(model, schema, table)
	if strings.Trim(query, " ") != passQuery {
		return fmt.Errorf("Expected query %s, got %s", passQuery, query)
	}

	if !passArgs(args) {
		return fmt.Errorf("Expected args [], got %v", args)
	}
	return nil
}

func Test_QueryGetNonNested(t *testing.T) {
	passQuery := `SELECT public.test_model.column5,public.test_model.column6 FROM public.test_model`
	passArgs := func(args []interface{}) bool {
		if len(args) > 0 {
			return false
		}
		return true
	}
	err := queryBuilderCommon(Dummy{}, passQuery, passArgs, queryBuilderJoin)
	if err != nil {
		t.Error(err)
	}
}

func Test_QueryGet(t *testing.T) {
	passQuery := `SELECT public.test_model.column1,public.test_model.column2,public.test_model.column3,public.test_model.column4,public.dummy.column5,public.dummy.column6 FROM public.test_model INNER JOIN public.dummy ON (test_model.column5=public.dummy.column5)`
	passArgs := func(args []interface{}) bool {
		if len(args) > 0 {
			return false
		}
		return true
	}
	err := queryBuilderCommon(TestModel{}, passQuery, passArgs, queryBuilderJoin)
	if err != nil {
		t.Error(err)
	}
}

func Test_QueryGetWithArgs(t *testing.T) {
	passQuery := `SELECT public.test_model.column1,public.test_model.column2,public.test_model.column3,public.test_model.column4,public.dummy.column5,public.dummy.column6 FROM public.test_model INNER JOIN public.dummy ON (test_model.column5=public.dummy.column5) WHERE test_model.column1=$1`
	passArgs := func(args []interface{}) bool {
		if len(args) != 1 || args[0] != "test" {
			return false
		}
		return true
	}
	err := queryBuilderCommon(TestModel{
		Column1: "test",
	}, passQuery, passArgs, queryBuilderJoin)
	if err != nil {
		t.Error(err)
	}
}

func Test_QueryCreate(t *testing.T) {
	passQuery := `INSERT INTO public.test_model (column1,column2,column3,column5) VALUES ($1,$2,$3,$4)`
	passArgs := func(args []interface{}) bool {
		if len(args) != 4 || args[0] != "test1" && args[1] != "test2" && args[2] != 69 && args[3] != "test3" {
			return false
		}
		return true
	}
	err := queryBuilderCommon(TestModel{
		Column1: "test1",
		Column2: "test2",
		Column3: 69,
		Column4: true,
		Column7: Dummy{
			Column5: "test3",
			Column6: "test4", // Shouldn't be inserted
		},
	}, passQuery, passArgs, queryBuilderCreate)
	if err != nil {
		t.Error(err)
	}
}

func Test_QueryUpdate(t *testing.T) {
	passQuery := `UPDATE public.test_model SET column1=$1,column2=$2,column3=$3,column5=$4 WHERE column1=$5`
	passArgs := func(args []interface{}) bool {
		if len(args) != 5 || args[0] != "test1" && args[1] != "test2" && args[2] != 69 && args[3] != "test3" && args[4] != "test1" {
			return false
		}
		return true
	}
	err := queryBuilderCommon(TestModel{
		Column1: "test1",
		Column2: "test2",
		Column3: 69,
		Column4: true,
		Column7: Dummy{
			Column5: "test3",
			Column6: "test4", // Shouldn't be inserted
		},
	}, passQuery, passArgs, queryBuilderUpdate)
	if err != nil {
		t.Error(err)
	}
}

func Test_QueryDelete(t *testing.T) {
	passQuery := `DELETE FROM public.test_model WHERE column1=$1`
	passArgs := func(args []interface{}) bool {
		if len(args) != 1 || args[0] != "test1" {
			return false
		}
		return true
	}
	err := queryBuilderCommon(TestModel{
		Column1: "test1",
	}, passQuery, passArgs, queryBuilderDelete)
	if err != nil {
		t.Error(err)
	}
}

func Test_QueryCount(t *testing.T) {
	passQuery := `SELECT COUNT(column1) FROM public.test_model`
	passArgs := func(args []interface{}) bool {
		if len(args) != 0 {
			return false
		}
		return true
	}
	err := queryBuilderCommon(TestModel{}, passQuery, passArgs, queryBuilderCount)
	if err != nil {
		t.Error(err)
	}
}
