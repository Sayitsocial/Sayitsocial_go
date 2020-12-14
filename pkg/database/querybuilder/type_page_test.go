package querybuilder

import "testing"

type PageTest struct {
	Column1 string `row:"column1" pk:"auto" type:"exact"`
	Column2 Page   `row:"column2" type:"exact"`
}

func (PageTest) GetTableName() (string, string) {
	return "public", "test_model"
}

func Test_PageGet(t *testing.T) {
	passQuery := `SELECT public.test_model.column1 FROM public.test_model`
	err := queryBuilderCommon(PageTest{}, passQuery, func(args []interface{}) bool {
		return len(args) == 0
	}, queryBuilderJoin)
	if err != nil {
		t.Error(err)
	}
}

func Test_PageGetWithArgs(t *testing.T) {
	passQuery := `SELECT public.test_model.column1 FROM public.test_model     LIMIT 15 OFFSET 30`
	err := queryBuilderCommon(PageTest{
		Column2: Page{
			Limit:  15,
			Offset: 30,
		},
	}, passQuery, func(args []interface{}) bool {
		return len(args) == 0
	}, queryBuilderJoin)
	if err != nil {
		t.Error(err)
	}
}

func Test_PageCreate(t *testing.T) {
	passQuery := `INSERT INTO public.test_model (column1) VALUES ($1)`
	err := queryBuilderCommon(PageTest{
		Column1: "test",
		Column2: Page{
			Limit:  15,
			Offset: 30,
		},
	}, passQuery, func(args []interface{}) bool {
		return !(len(args) != 1 || args[0] != "test")
	}, queryBuilderCreate)
	if err != nil {
		t.Error(err)
	}
}

func Test_PageUpdate(t *testing.T) {
	passQuery := `UPDATE public.test_model SET column1=$1 WHERE column1=$2`
	err := queryBuilderCommon(PageTest{
		Column1: "test",
		Column2: Page{
			Limit:  15,
			Offset: 30,
		},
	}, passQuery, func(args []interface{}) bool {
		return !(len(args) != 2 || (args[0] != "test" && args[1] != "test"))
	}, queryBuilderUpdate)
	if err != nil {
		t.Error(err)
	}
}
