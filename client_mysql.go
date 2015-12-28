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

var (
	mysqlCapabilities = ClientCapabilities{
		HasAnalyze: false,
	}
)

// ClientMysql describes MySQL (and derivatives) client
type ClientMysql struct {
	db               *sqlx.DB
	history          []HistoryRecord
	connectionString string
}

// NewClientMysqlFromURL opens a mysql db connection
func NewClientMysqlFromURL(uri string) (Client, error) {
	db, err := sqlx.Open("mysql", uri)
	if err != nil {
		return nil, err
	}
	return &ClientMysql{
		db:               db,
		connectionString: uri,
		history:          NewHistory(),
	}, nil
}

// GetCapabilities returns mysql capabilities
func (c *ClientMysql) GetCapabilities() ClientCapabilities {
	return mysqlCapabilities
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
	// Note: returning character_set_name as character_set_catalog
	// for pg compat. TODO: don't know if they mean the same thing
	q := `SELECT 
column_name, data_type, is_nullable, character_maximum_length, character_set_name as character_set_catalog, column_default
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

//http://stackoverflow.com/questions/14569940/mysql-list-tables-and-sizes-order-by-size
//http://stackoverflow.com/questions/5060366/mysql-fastest-way-to-count-number-of-rows
//http://stackoverflow.com/questions/9620198/how-to-get-the-sizes-of-the-tables-of-a-mysql-database
// TableInfo returns information about a given table
func (c *ClientMysql) TableInfo(table string) (*Result, error) {
	// TODO: filter by TABLE_SCHEMA i.e. database name
	q := `SELECT
  DATA_LENGTH AS data_size
, INDEX_LENGTH AS index_size
, TABLE_ROWS AS rows_count
FROM information_schema.tables
WHERE table_name = ?
`
	return dbQuery(c.db, q, table)
}

/*
TABLE_CATALOG,
TABLE_SCHEMA,TABLE_NAME,TABLE_TYPE,ENGINE,VERSION,ROW_FORMAT,TABLE_ROWS,AVG_ROW_LENGTH,
DATA_LENGTH,
MAX_DATA_LENGTH,
NDEX_LENGTH,
DATA_FREE,AUTO_INCREMENT,CREATE_TIME,UPDATE_TIME,CHECK_TIME,TABLE_COLLATION,CHECKSUM,CREATE_OPTIONS,TABLE_COMMENT
*/

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
