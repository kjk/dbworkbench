package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
)

// Row describes a database row
type Row []interface{}

// Result describes results of database query
type Result struct {
	Columns []string `json:"columns"`
	Rows    []Row    `json:"rows"`
}

// RowsOptions holds browsing options for table rows
type RowsOptions struct {
	Limit      int    // Number of rows to fetch
	SortColumn string // Column to sort by
	SortOrder  string // Sort direction (ASC, DESC)
}

// Format converts rows to a list of maps whose key is column
// name and value is value for that column in a given row.
// TODO: that is wasteful, it'll be more efficient to send this as
// an array that describes column names and array of arrays for rows
func (res *Result) Format() []map[string]interface{} {
	var items []map[string]interface{}

	for _, row := range res.Rows {
		item := make(map[string]interface{})

		for i, c := range res.Columns {
			item[c] = row[i]
		}

		items = append(items, item)
	}

	return items
}

// CSV creates csv representation of rows
func (res *Result) CSV() []byte {
	buff := &bytes.Buffer{}
	writer := csv.NewWriter(buff)

	writer.Write(res.Columns)

	for _, row := range res.Rows {
		record := make([]string, len(res.Columns))

		for i, item := range row {
			if item != nil {
				record[i] = fmt.Sprintf("%v", item)
			} else {
				record[i] = ""
			}
		}

		err := writer.Write(record)

		if err != nil {
			fmt.Println(err)
			break
		}
	}

	writer.Flush()
	return buff.Bytes()
}
