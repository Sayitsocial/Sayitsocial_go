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
	trans  *sql.Tx
	conn   *sql.DB
	Config Config
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
	return conn.BeginTx(ctx, options)
}

func Initialize(c *sql.DB, t *sql.Tx) (*Conn, error) {
	var conn *sql.DB
	var err error
	if c != nil {
		conn = c
	} else {
		conn, err = connectToDB()
		if err != nil {
			return nil, err
		}
	}
	return &Conn{
		trans: t,
		conn:  conn,
	}, nil
}

func getSlicePtr(i interface{}) interface{} {
	return reflect.New(reflect.SliceOf(reflect.TypeOf(i))).Interface()
}

func (c *Conn) queryMethod(i interface{}, method func([]colHolder, []foreignHolder, Config, string, string) (string, []interface{})) (*sql.Rows, error) {
	if val, ok := i.(Model); ok {
		schema, table := val.GetTableName()
		cols, foreign := generateColHolder(i, schema+"."+table, false)

		query, args := method(cols, foreign, c.Config, schema, table)
		row, err := c.conn.Query(query, args...)
		if err != nil {
			return nil, err
		}
		return row, nil
	}
	return nil, fmt.Errorf("Provided interface is not of type 'Model'")
}

func (c *Conn) Get(i interface{}) (interface{}, error) {
	s := getSlicePtr(i)
	row, err := c.queryMethod(i, queryBuilderJoin)
	if err != nil {
		return s, err
	}
	return s, getIntoStruct(row, s)
}

func (c *Conn) Create(i interface{}) error {
	_, err := c.queryMethod(i, queryBuilderCreate)
	if err != nil {
		return err
	}

	return nil
}

func (c *Conn) Delete(i interface{}) error {
	_, err := c.queryMethod(i, queryBuilderDelete)
	if err != nil {
		return err
	}
	return nil
}

func (c *Conn) Update(i interface{}) error {
	_, err := c.queryMethod(i, queryBuilderUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (c *Conn) Count(i interface{}) ([]int, error) {
	s := make([]int, 0)
	row, err := c.queryMethod(i, queryBuilderCount)
	if err != nil {
		return s, err
	}
	return s, getIntoVar(row, &s)
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
