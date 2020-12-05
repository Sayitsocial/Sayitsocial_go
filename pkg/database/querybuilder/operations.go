package querybuilder

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/Sayitsocial/Sayitsocial_go/pkg/helpers"
)

type Model interface {
	GetTableName() (string, string)
}

type Conn struct {
	trans *sql.Tx
	conn  *sql.DB
}

func getSlicePtr(i interface{}) interface{} {
	return reflect.New(reflect.SliceOf(reflect.TypeOf(i))).Interface()
}

func Initialize(t *sql.Tx) Conn {
	return Conn{
		trans: t,
		conn:  GetConn(),
	}
}

func (c Conn) Get(i interface{}) (interface{}, error) {
	if val, ok := i.(Model); ok {
		schema, table := val.GetTableName()
		isTableExist(c.conn, schema, table)

		query, args := queryBuilderJoin(i, schema, table)
		helpers.LogInfo(query)
		row, err := c.conn.Query(query, args...)
		if err != nil {
			return nil, err
		}

		s := getSlicePtr(i)
		err = GetIntoStruct(row, s)
		return s, err
	}
	return nil, fmt.Errorf("Provided interface is not of type 'Model'")
}

func (c Conn) Create(i interface{}) error {
	if val, ok := i.(Model); ok {
		schema, table := val.GetTableName()
		isTableExist(c.conn, schema, table)

		query, args := queryBuilderCreate(i, schema, table)
		_, err := c.conn.Query(query, args...)
		if err != nil {
			return err
		}
	}
	return fmt.Errorf("Provided interface is not of type 'Model'")
}

func (c Conn) Delete(i interface{}) error {
	if val, ok := i.(Model); ok {
		schema, table := val.GetTableName()
		isTableExist(c.conn, schema, table)

		query, args := queryBuilderDelete(i, table)
		_, err := c.conn.Query(query, args...)
		if err != nil {
			return err
		}
	}
	return fmt.Errorf("Provided interface is not of type 'Model'")
}

func (c Conn) Update(i interface{}) error {
	if val, ok := i.(Model); ok {
		schema, table := val.GetTableName()
		isTableExist(c.conn, schema, table)

		query, args := queryBuilderUpdate(i, schema, table)
		_, err := c.conn.Query(query, args...)
		if err != nil {
			return err
		}
	}
	return fmt.Errorf("Provided interface is not of type 'Model'")
}

func (c Conn) Count(i interface{}) ([]int, error) {
	if val, ok := i.(Model); ok {
		schema, table := val.GetTableName()
		isTableExist(c.conn, schema, table)

		query, args := queryBuilderCount(i, schema, table)
		row, err := c.conn.Query(query, args...)
		if err != nil {
			return nil, err
		}

		s := make([]int, 0)
		err = GetIntoVar(row, s)
		return s, err
	}
	return nil, fmt.Errorf("Provided interface is not of type 'Model'")
}

func (c Conn) Close() {
	err := c.conn.Close()
	if err != nil {
		helpers.LogError(err.Error())
	}
}
