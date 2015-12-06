package main

import (
	"reflect"

	"github.com/jmoiron/sqlx"
)

// Client describes a database connection
type Client interface {
	Test() error
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
	Close() error
}

// utility functions

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
