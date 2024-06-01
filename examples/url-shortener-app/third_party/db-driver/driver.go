package db_driver

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"strings"
	"sync"
)

// TODO: this code needs to be improved.

// The framework does not depend on third party packages, but we need some kind of database driver for example.
// This is a mock driver. However, you can replace this driver with a real database driver.

func init() {
	sql.Register("db-driver", &dbDriver{})
	sql.Register("db-driver-mock", &dbDriver{})
}

type dbDriver struct {
}

func (d *dbDriver) Open(_ string) (driver.Conn, error) {
	return &connection{
		storage: sync.Map{},
	}, nil
}

type connection struct {
	storage sync.Map
}

func (c *connection) Begin() (driver.Tx, error) {
	panic("unsupported method")
}

func (c *connection) Prepare(query string) (driver.Stmt, error) {
	return &stmt{
		connection: c,
		query:      query,
	}, nil
}

func (c *connection) Close() error {
	return nil
}

func (c *connection) Exec(query string, args []driver.Value) (driver.Result, error) {
	_, err := c.Query(query, args)
	return nil, err
}

func (c *connection) Query(query string, args []driver.Value) (driver.Rows, error) {
	if strings.Contains(query, "CREATE TABLE IF NOT EXISTS redirects") {
		return nil, nil
	}
	switch query {
	case `INSERT INTO redirects(url_key, url) VALUES(?, ?)`:
		c.storage.Store(args[0], args[1])
		return nil, nil
	case `SELECT url FROM redirects WHERE url_key = ?`:
		url, ok := c.storage.Load(args[0])
		data := make([][]driver.Value, 0, 1)
		if ok {
			data = append(data, []driver.Value{url})
		}
		return &rows{
			columns: []string{"url"},
			rows:    data,
		}, nil
	case `SELECT VERSION()`:
		return &rows{
			columns: []string{"version"},
			rows: [][]driver.Value{
				{"0.0.1"},
			},
		}, nil
	}
	errMessage := "We use a database mock so that the current framework package does not depend on other packages. "
	errMessage += "Please connect any other database driver to the application to use a real database. "
	return nil, fmt.Errorf("unsupported query: %s. %s", query, errMessage)
}

type stmt struct {
	connection *connection
	query      string
}

func (s *stmt) Exec(args []driver.Value) (driver.Result, error) {
	return s.connection.Exec(s.query, args)
}

func (s *stmt) Query(args []driver.Value) (driver.Rows, error) {
	return s.connection.Query(s.query, args)
}

func (s *stmt) NumInput() int {
	return -1
}

func (s *stmt) Close() error {
	return nil
}

type rows struct {
	columns  []string
	rows     [][]driver.Value
	position int
}

func (r *rows) Columns() []string {
	return r.columns
}

func (r *rows) Close() error {
	return nil
}

func (r *rows) Next(dest []driver.Value) error {
	if r.position >= len(r.rows) {
		return io.EOF
	}
	row := r.rows[r.position]
	for index, value := range row {
		dest[index] = value
	}
	r.position++
	return nil
}
