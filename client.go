package main

import (
	"reflect"

	"github.com/jmoiron/sqlx"
)

// Client describes a database connection
type Client interface {
	Connection() *sqlx.DB
	Info() (*Result, error)
	Databases() ([]string, error)
	Schemas() ([]string, error)
	Tables() ([]string, error)
	Table(table string) (*Result, error)
	TableRows(table string, opts RowsOptions) (*Result, error)
	TableInfo(table string) (*Result, error)
	TableIndexes(table string) (*Result, error)
	Activity() (*Result, error)
	Query(query string) (*Result, error)
	History() []HistoryRecord
	GetCapabilities() ClientCapabilities
}

// ClientCapabilities describes capabilities of the client.
// Front-end might customize the view depending on capabilities
type ClientCapabilities struct {
	// does it support query analyze
	HasAnalyze bool
}

// utility functions

func updateRow(row []interface{}) {
	for i, item := range row {
		if item == nil {
			row[i] = nil
			continue
		}
		t := reflect.TypeOf(item).Kind().String()

		if t == "slice" {
			row[i] = string(item.([]byte))
		}
	}
}

func dbQuery(db *sqlx.DB, query string, args ...interface{}) (*Result, error) {
	rows, err := db.Queryx(query, args...)

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
		row, err := rows.SliceScan()
		updateRow(row)

		if err == nil {
			result.Rows = append(result.Rows, row)
		}
	}

	return &result, nil
}

func dbFetchRows(db *sqlx.DB, q string) ([]string, error) {
	res, err := dbQuery(db, q)

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
