package querybuilder

import "testing"

type SortTest struct {
	Column1 string `row:"column1" pk:"auto" type:"exact"`
	Column2 SortBy `row:"column2" type:"exact"`
}

func (SortTest) GetTableName() (string, string) {
	return "public", "test_model"
}

func Test_SortGet(t *testing.T) {
	passQuery := `SELECT public.test_model.column1 FROM public.test_model`
	err := queryBuilderCommon(SortTest{}, passQuery, func(args []interface{}) bool {
		return len(args) == 0
	}, queryBuilderJoin)
	if err != nil {
		t.Error(err)
	}
}

func Test_SortGetWithArgs(t *testing.T) {
	passQuery := `SELECT public.test_model.column1 FROM public.test_model    ORDER BY column1 ASC`
	err := queryBuilderCommon(SortTest{
		Column2: SortBy{
			Column: "column1",
			Mode:   "ASC",
		},
	}, passQuery, func(args []interface{}) bool {
		return len(args) == 0
	}, queryBuilderJoin)
	if err != nil {
		t.Error(err)
	}
}

func Test_SortCreate(t *testing.T) {
	passQuery := `INSERT INTO public.test_model (column1) VALUES ($1)`
	err := queryBuilderCommon(SortTest{
		Column1: "test",
		Column2: SortBy{
			Column: "column1",
			Mode:   "ASC",
		},
	}, passQuery, func(args []interface{}) bool {
		return !(len(args) != 1 || args[0] != "test")
	}, queryBuilderCreate)
	if err != nil {
		t.Error(err)
	}
}

func Test_SortUpdate(t *testing.T) {
	passQuery := `UPDATE public.test_model SET column1=$1 WHERE column1=$2`
	err := queryBuilderCommon(SortTest{
		Column1: "test",
		Column2: SortBy{
			Column: "column1",
			Mode:   "ASC",
		},
	}, passQuery, func(args []interface{}) bool {
		return !(len(args) != 2 || (args[0] != "test" && args[1] != "test"))
	}, queryBuilderUpdate)
	if err != nil {
		t.Error(err)
	}
}
