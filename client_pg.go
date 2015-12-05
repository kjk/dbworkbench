package main

import (
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	pgDatabasesStmt = `SELECT datname FROM pg_database WHERE NOT datistemplate ORDER BY datname ASC`

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
func NewClientPgFromURL(url string) (*ClientPg, error) {
	if options.Debug {
		fmt.Println("Creating a new client for:", url)
	}

	db, err := sqlx.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	client := ClientPg{
		db:               db,
		connectionString: url,
		history:          NewHistory(),
	}

	return &client, nil
}

// Test checks if a db connection is valid
func (c *ClientPg) Test() error {
	return c.db.Ping()
}

// Info returns information about a postgres db connection
func (c *ClientPg) Info() (*Result, error) {
	return c.query(pgInfoStmt)
}

// Databases returns list of databases in a given postgres connection
func (c *ClientPg) Databases() ([]string, error) {
	return c.fetchRows(pgDatabasesStmt)
}

// Schemas returns list of schemas
func (c *ClientPg) Schemas() ([]string, error) {
	return c.fetchRows(pgSchemasStmt)
}

// Tables returns list of tables
func (c *ClientPg) Tables() ([]string, error) {
	return c.fetchRows(pgTablesStmt)
}

// Table returns schema for a given table
func (c *ClientPg) Table(table string) (*Result, error) {
	return c.query(pgTableSchemaStmt, table)
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

	return c.query(sql)
}

// TableInfo returns information about a given table
func (c *ClientPg) TableInfo(table string) (*Result, error) {
	return c.query(pgTableInfoStmt, table)
}

// TableIndexes returns info about indexes for a given table
func (c *ClientPg) TableIndexes(table string) (*Result, error) {
	res, err := c.query(pgTableIndexesStmt, table)

	if err != nil {
		return nil, err
	}

	return res, err
}

// Activity returns all active queriers on the server
func (c *ClientPg) Activity() (*Result, error) {
	return c.query(pgActivityStmt)
}

// Query executes a given query and returns the results
func (c *ClientPg) Query(query string) (*Result, error) {
	res, err := c.query(query)

	// Save history records only if query did not fail
	if err == nil {
		c.history = append(c.history, NewHistoryRecord(query))
	}

	return res, err
}

func (c *ClientPg) query(query string, args ...interface{}) (*Result, error) {
	rows, err := c.db.Queryx(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := Result{Columns: cols}

	for rows.Next() {
		obj, err := rows.SliceScan()

		for i, item := range obj {
			if item == nil {
				obj[i] = nil
			} else {
				t := reflect.TypeOf(item).Kind().String()

				if t == "slice" {
					obj[i] = string(item.([]byte))
				}
			}
		}

		if err == nil {
			result.Rows = append(result.Rows, obj)
		}
	}

	return &result, nil
}

// Fetch all rows as strings for a single column
func (c *ClientPg) fetchRows(q string) ([]string, error) {
	res, err := c.query(q)

	if err != nil {
		return nil, err
	}

	// Init empty slice so json.Marshal will encode it to "[]" instead of "null"
	var results []string

	for _, row := range res.Rows {
		results = append(results, row[0].(string))
	}

	return results, nil
}
