package querybuilder

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

var connString string
var dbDriver string

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

func SetConnection(pgConn string, driver string) {
	connString = pgConn
	dbDriver = driver
}

func connectToDB() (*sql.DB, error) {
	if dbDriver != "" && connString != "" {
		conn, err := sql.Open(dbDriver, connString)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
	return nil, fmt.Errorf("Connection string or driver can't be empty")
}

func GetTransaction(ctx context.Context, options *sql.TxOptions) (*sql.Tx, error) {
	conn, err := connectToDB()
	if err != nil {
		return nil, err
	}
	return conn.BeginTx(ctx, nil)
}

func Initialize(t *sql.Tx) Conn {
	conn, err := connectToDB()
	if err != nil {
		// TODO: Return error
	}
	return Conn{
		trans: t,
		conn:  conn,
	}
}

func InitializeWithConn(c *sql.DB) Conn {
	return Conn{
		trans: nil,
		conn:  c,
	}
}

func (c Conn) Get(i interface{}) (interface{}, error) {
	if val, ok := i.(Model); ok {
		schema, table := val.GetTableName()
		isTableExist(c.conn, schema, table)

		query, args := queryBuilderJoin(i, schema, table)
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

		query, args := queryBuilderDelete(i, schema, table)
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

func (c Conn) Close() error {
	return c.conn.Close()
}
