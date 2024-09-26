package db_driver

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"sync"
)

// The framework does not depend on third party packages, but we need some kind of database driver for example.
// This is a mock driver. However, you can replace this driver with a real database driver.

var (
	storages = map[string]*sync.Map{}
	mutex    = sync.Mutex{}
)

func init() {
	sql.Register("db-driver", &dbDriver{})
	sql.Register("db-driver-mock", &dbDriver{})
}

type dbDriver struct{}

func (d *dbDriver) Open(name string) (driver.Conn, error) {
	return &connection{
		storage: getStorage(name),
	}, nil
}

type connection struct {
	storage *sync.Map
	closed  bool
}

func (c *connection) Begin() (driver.Tx, error) {
	return nil, fmt.Errorf("transaction is not supported by this database driver")
}

func (c *connection) Prepare(query string) (driver.Stmt, error) {
	return c.PrepareContext(context.Background(), query)
}

func (c *connection) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return &stmt{
		ctx:        ctx,
		connection: c,
		query:      query,
	}, nil
}

func (c *connection) Close() error {
	mutex.Lock()
	c.closed = true
	mutex.Unlock()
	return nil
}

func (c *connection) isClosed() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return c.closed
}

type stmt struct {
	ctx        context.Context
	connection *connection
	query      string
}

func (s *stmt) Exec(args []driver.Value) (driver.Result, error) {
	return s.ExecContext(s.ctx, convertValuesToNamedValues(args))
}

func (s *stmt) Query(args []driver.Value) (driver.Rows, error) {
	return s.QueryContext(s.ctx, convertValuesToNamedValues(args))
}

func (s *stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	_, err := s.QueryContext(ctx, args)
	return nil, err
}

func (s *stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (rows driver.Rows, err error) {
	if s.connection.isClosed() {
		return nil, sql.ErrConnDone
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("an error occurred while executing the query '%s' in the database: %v", s.query, r)
		}
	}()
	for _, item := range queries {
		if (item.matcher == nil && s.query == item.query) || (item.matcher != nil && item.matcher(s.query, item.query)) {
			return item.handler(ctx, s.connection, args)
		}
	}
	errMessage := "We use a database mock so that the framework package does not depend on other packages. "
	errMessage += "Please connect any other database driver to the application to use a real database. "
	return nil, fmt.Errorf("unsupported query: %s. %s", s.query, errMessage)
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

func getStorage(name string) *sync.Map {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := storages[name]; !ok {
		storages[name] = &sync.Map{}
	}
	return storages[name]
}

func convertValuesToNamedValues(values []driver.Value) []driver.NamedValue {
	namedValues := make([]driver.NamedValue, len(values))
	for index, value := range values {
		namedValues[index] = driver.NamedValue{
			Name:    "",
			Ordinal: index + 1,
			Value:   value,
		}
	}
	return namedValues
}

var (
	_ driver.Driver             = (*dbDriver)(nil)
	_ driver.Conn               = (*connection)(nil)
	_ driver.ConnPrepareContext = (*connection)(nil)
	_ driver.Stmt               = (*stmt)(nil)
	_ driver.StmtExecContext    = (*stmt)(nil)
	_ driver.StmtQueryContext   = (*stmt)(nil)
	_ driver.Rows               = (*rows)(nil)
)
