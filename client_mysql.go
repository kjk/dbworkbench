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
	// http://dev.mysql.com/doc/refman/5.0/en/show-databases.html
	mysqlDatabasesStmt = `SHOW DATABASES`

	// Note: probably not used
	mysqlSchemasStmt = `select schema_name from information_schema.schemata`

	// note: doesb't have as many fields as in postgres
	mysqlInfoStmt = `SELECT user() AS session_user,
current_user,
database() as current_database,
version() AS version`

	// returns version of mysql database e.g. 5.5.46
	mysqlVersionStmt = `SELECT VARIABLE_NAME, VARIABLE_VALUE FROM INFORMATION_SCHEMA.GLOBAL_VARIABLES WHERE VARIABLE_NAME = 'VERSION';`

	// http://dev.mysql.com/doc/refman/5.0/en/show-tables.html
	// TODO: possibliy rewrite as a query since it differs depending on mysql version
	// https://dev.mysql.com/doc/refman/5.0/en/tables-table.html
	mysqlTablesStmt = `SHOW TABLES`
	// TODO: add equivalent of table_schema = 'public'
	mysqlTablesStmt2 = `SELECT
table_name FROM information_schema.tables
WHERE table_type = 'BASE TABLE'
ORDER BY table_schema, table_name`

	// https://dev.mysql.com/doc/refman/5.0/en/columns-table.html
	// TODO: don't know if CHARACTER_SET_NAME is the same as character_set_catalog
	mysqlTableSchemaStmt = `SELECT 
column_name, data_type, is_nullable, character_maximum_length, character_set_name, column_default
FROM information_schema.columns
WHERE table_name = ?`

	mysqlActivityStmt = `SHOW FULL PROCESSLIST;`

	//mysqlTableIndexesStmt = `SELECT indexname, indexdef FROM pg_indexes WHERE tablename = $1`
	// http://stackoverflow.com/questions/5213339/how-to-see-indexes-for-a-database-or-table
	mysqlTableIndexesStmt = `SHOW INDEX FROM %s`
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

// Connection returns underlying db connection
func (c *ClientMysql) Connection() *sqlx.DB {
	return c.db
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
	return dbQuery(c.db, pgTableInfoStmt, table)
}

// TableIndexes returns info about indexes for a given table
func (c *ClientMysql) TableIndexes(table string) (*Result, error) {
	q := fmt.Sprintf(mysqlTableIndexesStmt, table)
	res, err := dbQuery(c.db, q)

	if err != nil {
		return nil, err
	}

	return res, err
}

// Activity returns all active queriers on the server
func (c *ClientMysql) Activity() (*Result, error) {
	return dbQuery(c.db, mysqlActivityStmt)
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
