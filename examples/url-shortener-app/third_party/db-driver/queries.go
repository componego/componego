package db_driver

import (
	"context"
	"database/sql/driver"
	"strings"
)

// List of hardcoded database queries because the driver in this example does not provide a real connection to the database.
var queries = []struct {
	query   string
	matcher func(actualQuery string, expectedQuery string) bool
	handler func(ctx context.Context, conn *connection, args []driver.NamedValue) (driver.Rows, error)
}{
	{
		query: `CREATE TABLE IF NOT EXISTS redirects`,
		matcher: func(actualQuery string, expectedQuery string) bool {
			return strings.Contains(actualQuery, expectedQuery)
		},
		handler: func(_ context.Context, _ *connection, _ []driver.NamedValue) (driver.Rows, error) {
			return nil, nil
		},
	},
	{
		query: `INSERT INTO redirects(url_key, url) VALUES(?, ?)`,
		handler: func(_ context.Context, conn *connection, args []driver.NamedValue) (driver.Rows, error) {
			conn.storage.Store(args[0].Value, args[1].Value)
			return nil, nil
		},
	},
	{
		query: `SELECT url FROM redirects WHERE url_key = ?`,
		handler: func(_ context.Context, conn *connection, args []driver.NamedValue) (driver.Rows, error) {
			return &rows{
				columns: []string{"url"},
				rows: func() [][]driver.Value {
					if url, ok := conn.storage.Load(args[0].Value); ok {
						return [][]driver.Value{
							{url.(string)},
						}
					}
					return make([][]driver.Value, 0)
				}(),
			}, nil
		},
	},
	{
		query: `SELECT VERSION()`,
		handler: func(_ context.Context, _ *connection, _ []driver.NamedValue) (driver.Rows, error) {
			return &rows{
				columns: []string{"version"},
				rows: [][]driver.Value{
					{"0.0.1"},
				},
			}, nil
		},
	},
}
