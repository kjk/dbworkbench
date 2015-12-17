package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	pgDatabasesStmt = `SELECT datname FROM pg_database WHERE NOT datistemplate ORDER BY datname ASC`

	// Note: probably not used
	pgSchemasStmt = `SELECT schema_name FROM information_schema.schemata ORDER BY schema_name ASC`

	pgInfoStmt = `SELECT
  session_user
, current_user
, current_database()
, current_schemas(false)
, inet_client_addr()
, inet_client_port()
, inet_server_addr()
, inet_server_port()
, version()`

	pgTableIndexesStmt = `SELECT indexname, indexdef FROM pg_indexes WHERE tablename = $1`

	pgTableInfoStmt = `SELECT
  pg_size_pretty(pg_table_size($1)) AS data_size
, pg_size_pretty(pg_indexes_size($1)) AS index_size
, pg_size_pretty(pg_total_relation_size($1)) AS total_size
, (SELECT reltuples FROM pg_class WHERE oid = $1::regclass) AS rows_count`

	pgTableSchemaStmt = `SELECT
column_name, data_type, is_nullable, character_maximum_length, character_set_catalog, column_default
FROM information_schema.columns
WHERE table_name = $1`

	pgTablesStmt = `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_schema,table_name`

	pgActivityStmt = `SELECT
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
)

// ClientPg describes Postgres db client
type ClientPg struct {
	db               *sqlx.DB
	history          []HistoryRecord
	connectionString string
}

// NewClientPgFromURL opens a Postgres db connection
func NewClientPgFromURL(uri string) (Client, error) {
	if options.Debug {
		fmt.Println("Creating a new client for:", uri)
	}

	db, err := sqlx.Open("postgres", uri)

	if err != nil {
		return nil, err
	}

	client := ClientPg{
		db:               db,
		connectionString: uri,
		history:          NewHistory(),
	}

	return &client, nil
}

// Connection returns underlying db connection
func (c *ClientPg) Connection() *sqlx.DB {
	return c.db
}

// Info returns information about a postgres db connection
func (c *ClientPg) Info() (*Result, error) {
	return dbQuery(c.db, pgInfoStmt)
}

// Databases returns list of databases in a given postgres connection
func (c *ClientPg) Databases() ([]string, error) {
	return dbFetchRows(c.db, pgDatabasesStmt)
}

// Schemas returns list of schemas
func (c *ClientPg) Schemas() ([]string, error) {
	return dbFetchRows(c.db, pgSchemasStmt)
}

// Tables returns list of tables
func (c *ClientPg) Tables() ([]string, error) {
	return dbFetchRows(c.db, pgTablesStmt)
}

// Table returns schema for a given table
func (c *ClientPg) Table(table string) (*Result, error) {
	return dbQuery(c.db, pgTableSchemaStmt, table)
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
	return dbQuery(c.db, pgTableInfoStmt, table)
}

// TableIndexes returns info about indexes for a given table
func (c *ClientPg) TableIndexes(table string) (*Result, error) {
	res, err := dbQuery(c.db, pgTableIndexesStmt, table)

	if err != nil {
		return nil, err
	}

	return res, err
}

// Activity returns all active queriers on the server
func (c *ClientPg) Activity() (*Result, error) {
	return dbQuery(c.db, pgActivityStmt)
}

// Query executes a given query and returns the results
func (c *ClientPg) Query(query string) (*Result, error) {
	res, err := dbQuery(c.db, query)

	// Save history records only if query did not fail
	if err == nil {
		c.history = append(c.history, NewHistoryRecord(query))
	}

	return res, err
}

// History returns history of queries
func (c *ClientPg) History() []HistoryRecord {
	return c.history
}

func connectPostgres(uri string) (Client, error) {
	sslModes := []string{"require", "disable", "verify-full"}
	var firstError error
	for _, sslMode := range sslModes {
		fullURI := uri + "?sslmode=" + sslMode
		client, err := NewClientPgFromURL(fullURI)
		if err != nil {
			if firstError == nil {
				firstError = err
			}
			LogVerbosef("NewClientPgFromURL('%s') failed with '%s'\n", fullURI, err)
			continue
		}
		db := client.Connection()
		err = db.Ping()
		if err == nil {
			return client, nil
		}
		LogVerbosef("client.Test() failed with '%s', uri: '%s'\n", err, fullURI)
		if firstError == nil {
			firstError = err
		}
	}
	return nil, firstError
}
