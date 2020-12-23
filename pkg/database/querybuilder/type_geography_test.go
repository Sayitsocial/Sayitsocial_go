package querybuilder

import (
	"testing"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/database/querybuilder/types"
)

type GeographyTest struct {
	Column1 string                `row:"column1" pk:"auto" type:"exact"`
	Column2 types.GeographyPoints `row:"column2" type:"exact"`
}

func (GeographyTest) GetTableName() (string, string) {
	return "public", "test_model"
}

func Test_GeographyGet(t *testing.T) {
	passQuery := `SELECT public.test_model.column1,ST_X(public.test_model.column2::geometry),ST_Y(public.test_model.column2::geometry) FROM public.test_model`
	err := queryBuilderCommon(GeographyTest{}, passQuery, Config{}, func(args []interface{}) bool {
		return len(args) == 0
	}, queryBuilderJoin)
	if err != nil {
		t.Error(err)
	}
}

func Test_GeographyGetWithArgs(t *testing.T) {
	passQuery := `SELECT public.test_model.column1,ST_X(public.test_model.column2::geometry),ST_Y(public.test_model.column2::geometry) FROM public.test_model  WHERE ST_DWithin(public.test_model.column2,ST_MakePoint(1.1,2.2),3000)`
	err := queryBuilderCommon(GeographyTest{
		Column2: types.GeographyPoints{
			Longitude: "1.1",  // Required
			Latitude:  "2.2",  // Required
			Radius:    "3000", // Required
		},
	}, passQuery, Config{}, func(args []interface{}) bool {
		return len(args) == 0
	}, queryBuilderJoin)
	if err != nil {
		t.Error(err)
	}
}

func Test_GeographyCreate(t *testing.T) {
	passQuery := `INSERT INTO public.test_model (column1,column2) VALUES ($1,ST_SetSRID(ST_MakePoint(1.1,2.2),4326))`
	err := queryBuilderCommon(GeographyTest{
		Column1: "test",
		Column2: types.GeographyPoints{
			Longitude: "1.1", // Required
			Latitude:  "2.2", // Required
		},
	}, passQuery, Config{}, func(args []interface{}) bool {
		return !(len(args) != 1 || args[0] != "test")
	}, queryBuilderCreate)
	if err != nil {
		t.Error(err)
	}
}

func Test_GeographyUpdate(t *testing.T) {
	passQuery := `UPDATE public.test_model SET column1=$1,column2=ST_SetSRID(ST_MakePoint(1.1,2.2),4326) WHERE column1=$2`
	err := queryBuilderCommon(GeographyTest{
		Column1: "test",
		Column2: types.GeographyPoints{
			Longitude: "1.1", // Required
			Latitude:  "2.2", // Required
		},
	}, passQuery, Config{}, func(args []interface{}) bool {
		return !(len(args) != 2 || (args[0] != "test" && args[1] != "ST_SetSRID(ST_MakePoint(%v,%v),4326)"))
	}, queryBuilderUpdate)
	if err != nil {
		t.Error(err)
	}
}

// func Test_GeographyDelete(t *testing.T) {
// 	passQuery := `DELETE FROM public.test_model WHERE column2=$1`
// 	err := queryBuilderCommon(GeographyTest{
// 		Column2: GeographyPoints{
// 			Longitude: "1.1", // Required
// 			Latitude:  "2.2", // Required
// 		},
// 	}, passQuery, func(args []interface{}) bool {
// 		if val, ok := args[0].(GeographyPoints); ok {
// 			return (val.Longitude == "1.1" && val.Latitude == "2.2")
// 		}
// 		return false
// 	}, queryBuilderDelete)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func Test_GeographyCount(t *testing.T) {
// 	passQuery := `SELECT COUNT(column1) FROM public.test_model WHERE column2=$2`
// 	err := queryBuilderCommon(GeographyTest{
// 		Column2: GeographyPoints{
// 			Longitude: "1.1", // Required
// 			Latitude:  "2.2", // Required
// 		},
// 	}, passQuery, func(args []interface{}) bool {
// 		if val, ok := args[0].(GeographyPoints); ok {
// 			return (val.Longitude == "1.1" && val.Latitude == "2.2")
// 		}
// 		return false
// 	}, queryBuilderCount)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
