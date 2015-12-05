package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// http://dev.mysql.com/doc/refman/5.0/en/information-schema.html
// Useful statements:
// select version()

const (
	// TODO: write me
	mysqlDatabasesStmt = `SELECT version()`

	mysqlSchemasStmt = `select schema_name from information_schema.schemata`

	// returns version of mysql database e.g. 5.5.46
	mysqlVersionStmt = `SELECT VARIABLE_NAME, VARIABLE_VALUE FROM INFORMATION_SCHEMA.GLOBAL_VARIABLES WHERE VARIABLE_NAME = 'VERSION';`

	// TODO: write me
	mysqlInfoStmt        = `SELECT version()`
	mysqlTablesStmt      = `SELECT version()`
	mysqlTableSchemaStmt = `SELECT version()`
)

// ClientMysql describes MySQL (and derivatives) client
type ClientMysql struct {
	db               *sqlx.DB
	history          []HistoryRecord
	connectionString string
}

// NewClientMysqlFromURL opens a Postgres db connection
func NewClientMysqlFromURL(uri string) (Client, error) {
	db, err := sqlx.Open("mysql", uri)
	if err != nil {
		return nil, err
	}
	client := ClientMysql{
		db:               db,
		connectionString: uri,
		history:          NewHistory(),
	}

	return &client, nil
}

// Test checks if a db connection is valid
func (c *ClientMysql) Test() error {
	return c.db.Ping()
}

// Info returns information about a postgres db connection
func (c *ClientMysql) Info() (*Result, error) {
	return dbQuery(c.db, mysqlInfoStmt)
}

// Databases returns list of databases in a given postgres connection
func (c *ClientMysql) Databases() ([]string, error) {
	return dbFetchRows(c.db, mysqlDatabasesStmt)
}

// Schemas returns list of schemas
func (c *ClientMysql) Schemas() ([]string, error) {
	return dbFetchRows(c.db, mysqlSchemasStmt)
}

// Tables returns list of tables
func (c *ClientMysql) Tables() ([]string, error) {
	return dbFetchRows(c.db, mysqlTablesStmt)
}

// Table returns schema for a given table
func (c *ClientMysql) Table(table string) (*Result, error) {
	return dbQuery(c.db, mysqlTableSchemaStmt, table)
}

// TableRows returns all rows from a query
func (c *ClientMysql) TableRows(table string, opts RowsOptions) (*Result, error) {
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
func (c *ClientMysql) TableInfo(table string) (*Result, error) {
	return dbQuery(c.db, pgTableInfoStmt, table)
}

// TableIndexes returns info about indexes for a given table
func (c *ClientMysql) TableIndexes(table string) (*Result, error) {
	res, err := dbQuery(c.db, pgTableIndexesStmt, table)

	if err != nil {
		return nil, err
	}

	return res, err
}

// Activity returns all active queriers on the server
func (c *ClientMysql) Activity() (*Result, error) {
	return dbQuery(c.db, pgActivityStmt)
}

// Query executes a given query and returns the results
func (c *ClientMysql) Query(query string) (*Result, error) {
	res, err := dbQuery(c.db, query)

	// Save history records only if query did not fail
	if err == nil {
		c.history = append(c.history, NewHistoryRecord(query))
	}

	return res, err
}

// History returns history of queries
func (c *ClientMysql) History() []HistoryRecord {
	return c.history
}

// Close closes a database connection
func (c *ClientMysql) Close() error {
	return c.db.Close()
}
