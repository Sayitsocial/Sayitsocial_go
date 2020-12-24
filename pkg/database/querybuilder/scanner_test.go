package querybuilder

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type Test struct {
	Column1 string `sorm:"column1,pk_autoinc"`
	Column2 string `sorm:"column2"`
}

func (Test) GetTableName() (string, string) {
	return "public", "test_model"
}

func Test_GetIntoStruct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"column1", "column2"}).AddRow("test", "test1")
	mock.ExpectQuery("SELECT public.test_model.column1,public.test_model.column2 FROM public.test_model").WillReturnRows(rows)

	model, err := Initialize(db, nil)
	if err != nil {
		t.Error(err)
		return
	}
	row, err := model.queryMethod(Test{}, queryBuilderJoin)
	if err != nil {
		t.Error(err)
		return
	}

	s := getSlicePtr(Test{})
	getIntoStruct(row, s)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if val, ok := s.(Test); ok {
		if val.Column1 != "test" && val.Column2 != "test1" {
			t.Errorf("Expected &[{test test1}], got %v", s)
			return
		}
	}
}

func Test_GetIntoVar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow("3")
	mock.ExpectQuery("^SELECT (.+) FROM public.test_model$").WillReturnRows(rows)

	model, err := Initialize(db, nil)
	if err != nil {
		t.Error(err)
	}

	row, err := model.queryMethod(Test{}, queryBuilderCount)
	if err != nil {
		t.Error(err)
	}

	s := make([]int, 0)
	getIntoVar(row, &s)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if !(len(s) > 0 && s[0] == 3) {
		t.Errorf("Expected [3], got %v", s)
	}
}
