package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// http://dev.mysql.com/doc/refman/5.0/en/information-schema.html
// Useful statements:
// select version()
// `SELECT VARIABLE_NAME, VARIABLE_VALUE FROM INFORMATION_SCHEMA.GLOBAL_VARIABLES WHERE VARIABLE_NAME = 'VERSION';`

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

// Connection returns underlying db connection
func (c *ClientMysql) Connection() *sqlx.DB {
	return c.db
}

// Info returns information about a postgres db connection
func (c *ClientMysql) Info() (*Result, error) {
	// note: doesb't have as many fields as in postgres
	q := `SELECT user() AS session_user,
current_user,
database() as current_database,
version() AS version`
	return dbQuery(c.db, q)
}

// Databases returns list of databases in a given postgres connection
func (c *ClientMysql) Databases() ([]string, error) {
	// http://dev.mysql.com/doc/refman/5.0/en/show-databases.html
	q := `SHOW DATABASES`
	return dbFetchRows(c.db, q)
}

// Schemas returns list of schemas
func (c *ClientMysql) Schemas() ([]string, error) {
	// Note: probably not used
	q := `select schema_name from information_schema.schemata`
	return dbFetchRows(c.db, q)
}

// Tables returns list of tables
func (c *ClientMysql) Tables() ([]string, error) {
	// http://dev.mysql.com/doc/refman/5.0/en/show-tables.html
	// TODO: possibliy rewrite as a query since it differs depending on mysql version
	// https://dev.mysql.com/doc/refman/5.0/en/tables-table.html
	q := `SHOW TABLES`
	// TODO: add equivalent of table_schema = 'public'
	/*
			q := `SELECT
		table_name FROM information_schema.tables
		WHERE table_type = 'BASE TABLE'
		ORDER BY table_schema, table_name`
	*/
	return dbFetchRows(c.db, q)
}

// Table returns schema for a given table
func (c *ClientMysql) Table(table string) (*Result, error) {
	// https://dev.mysql.com/doc/refman/5.0/en/columns-table.html
	// TODO: don't know if CHARACTER_SET_NAME is the same as character_set_catalog
	q := `SELECT 
column_name, data_type, is_nullable, character_maximum_length, character_set_name, column_default
FROM information_schema.columns
WHERE table_name = ?`
	return dbQuery(c.db, q, table)
}

// TableRows returns all rows from a query
func (c *ClientMysql) TableRows(table string, opts RowsOptions) (*Result, error) {
	sql := fmt.Sprintf(`SELECT * FROM %s`, table)

	if opts.SortColumn != "" {
		if opts.SortOrder == "" {
			opts.SortOrder = "ASC"
		}

		sql += fmt.Sprintf(" ORDER BY %s %s", opts.SortColumn, opts.SortOrder)
	}

	if opts.Limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", opts.Limit)
	}
	LogInfof("sql: '%s'\n", sql)
	return dbQuery(c.db, sql)
}

// TableInfo returns information about a given table
func (c *ClientMysql) TableInfo(table string) (*Result, error) {
	// TODO: clearly wrong
	q := `SELECT
  pg_size_pretty(pg_table_size($1)) AS data_size
, pg_size_pretty(pg_indexes_size($1)) AS index_size
, pg_size_pretty(pg_total_relation_size($1)) AS total_size
, (SELECT reltuples FROM pg_class WHERE oid = $1::regclass) AS rows_count`

	return dbQuery(c.db, q, table)
}

// TableIndexes returns info about indexes for a given table
func (c *ClientMysql) TableIndexes(table string) (*Result, error) {
	// http://stackoverflow.com/questions/5213339/how-to-see-indexes-for-a-database-or-table
	q := fmt.Sprintf(`SHOW INDEX FROM %s`, table)
	res, err := dbQuery(c.db, q)

	if err != nil {
		return nil, err
	}

	return res, err
}

// Activity returns all active queriers on the server
func (c *ClientMysql) Activity() (*Result, error) {
	q := `SHOW FULL PROCESSLIST;`
	return dbQuery(c.db, q)
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

func connectMysql(uri string) (Client, error) {
	client, err := NewClientMysqlFromURL(uri)
	if err != nil {
		LogVerbosef("NewClientMysqlFromURL('%s') failed with '%s'\n", uri, err)
		return nil, err
	}
	db := client.Connection()
	err = db.Ping()
	if err != nil {
		LogVerbosef("client.Test() failed with '%s', uri: '%s'\n", err, uri)
		return nil, err
	}
	return client, nil
}
