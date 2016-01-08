package main

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	pgCapabilities = ClientCapabilities{
		HasAnalyze: true,
	}
)

// ClientPg describes Postgres db client
type ClientPg struct {
	History
	db               *sqlx.DB
	connectionString string
}

// NewClientPgFromURL opens a Postgres db connection
func NewClientPgFromURL(uri string) (Client, error) {
	db, err := sqlx.Open("postgres", uri)

	if err != nil {
		return nil, err
	}

	return &ClientPg{
		db:               db,
		connectionString: uri,
	}, nil
}

// GetCapabilities returns mysql capabilities
func (c *ClientPg) GetCapabilities() ClientCapabilities {
	return pgCapabilities
}

// Connection returns underlying db connection
func (c *ClientPg) Connection() *sqlx.DB {
	return c.db
}

// Info returns information about a postgres db connection
func (c *ClientPg) Info() (*Result, error) {
	q := `SELECT
  session_user
, current_user
, current_database()
, current_schemas(false)
, inet_client_addr()
, inet_client_port()
, inet_server_addr()
, inet_server_port()
, version()`
	return dbQuery(c.db, q)
}

// Databases returns list of databases in a given postgres connection
func (c *ClientPg) Databases() ([]string, error) {
	q := `SELECT datname FROM pg_database WHERE NOT datistemplate ORDER BY datname ASC`
	return dbFetchRows(c.db, q)
}

// Schemas returns list of schemas
// Note: probably not used
func (c *ClientPg) Schemas() ([]string, error) {
	q := `SELECT schema_name FROM information_schema.schemata ORDER BY schema_name ASC`
	return dbFetchRows(c.db, q)
}

// Tables returns list of tables
func (c *ClientPg) Tables() ([]*TableInfo, error) {
	q := `SELECT
			column_name,
			data_type,
			is_nullable,
			character_maximum_length,
			character_set_catalog,
			column_default,
			table_name,
			table_schema
		FROM information_schema.columns
		WHERE table_name IN (SELECT table_name FROM information_schema.tables WHERE table_schema = 'public')`

	return dbQueryTableInfo(c.db, q)
}

// Table returns schema for a given table
func (c *ClientPg) Table(table string) (*Result, error) {
	q := `SELECT
			column_name,
			data_type, is_nullable,
			character_maximum_length,
			character_set_catalog,
			column_default,
			table_name,
			table_schema
		FROM information_schema.columns
		WHERE table_name = $1`

	return dbQuery(c.db, q, table)
}

// TableRows returns all rows from a query
func (c *ClientPg) TableRows(table string, opts RowsOptions) (*Result, error) {
	sql := fmt.Sprintf(`SELECT * FROM "%s"`, table)

	if opts.SortColumn != "" {
		if opts.SortOrder == "" {
			opts.SortOrder = "ASC"
		}

		sql += fmt.Sprintf(" ORDER BY %s %s", opts.SortColumn, opts.SortOrder)
	}

	if opts.Limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}

	return dbQuery(c.db, sql)
}

// TableInfo returns information about a given table
func (c *ClientPg) TableInfo(table string) (*Result, error) {
	q := `SELECT
  pg_table_size($1) AS data_size
, pg_indexes_size($1) AS index_size
, (SELECT reltuples FROM pg_class WHERE oid = $1::regclass) AS rows_count`
	return dbQuery(c.db, q, table)
}

// TableIndexes returns info about indexes for a given table
func (c *ClientPg) TableIndexes(table string) (*Result, error) {
	q := `SELECT indexname, indexdef FROM pg_indexes WHERE tablename = $1`
	res, err := dbQuery(c.db, q, table)

	if err != nil {
		return nil, err
	}

	return res, err
}

// Activity returns all active queriers on the server
func (c *ClientPg) Activity() (*Result, error) {
	q := `SELECT
  datname,
  query,
  state,
  waiting,
  query_start,
  state_change,
  pid,
  datid,
  application_name,
  client_addr
  FROM pg_stat_activity
  WHERE state IS NOT NULL`
	return dbQuery(c.db, q)
}

// Query executes a given query and returns the results
func (c *ClientPg) Query(query string) (*Result, error) {
	res, err := dbQuery(c.db, query)

	// Save history records only if query did not fail
	if err == nil {
		c.AddToHistory(query)
	}

	return res, err
}

// a heuristic that detects if a given error is a tcp connection timeout
// The error messages look like:
// 'dial tcp 173.194.251.111:5432: getsockopt: operation timed out'
func isTimeoutError(err error) bool {
	return strings.Contains(err.Error(), "operation timed out")
}

func connectPostgres(uri string) (Client, error) {
	// TODO: doing 'verify-full' probably doesn't make sense as it's a superset
	// of 'require'
	sslModes := []string{"require", "disable", "verify-full"}
	var firstError error
	for _, sslMode := range sslModes {
		fullURI := uri + "?sslmode=" + sslMode
		client, err := NewClientPgFromURL(fullURI)
		if err != nil {
			if firstError == nil {
				firstError = err
			}
			LogErrorf("NewClientPgFromURL('%s') failed with '%s'\n", fullURI, err)
			continue
		}
		db := client.Connection()
		err = db.Ping()
		if err == nil {
			return client, nil
		}
		LogErrorf("db.Ping() failed with '%s', uri: '%s'\n", err, fullURI)
		if firstError == nil {
			firstError = err
		}
		// don't retry connections if the issue was a timeout. That would triple
		// how long does it take to timeout
		if isTimeoutError(err) {
			return nil, err
		}
	}
	return nil, firstError
}
