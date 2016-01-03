package main

import (
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"github.com/jmoiron/sqlx"
)

// Client describes a database connection
type Client interface {
	Connection() *sqlx.DB
	Info() (*Result, error)
	Databases() ([]string, error)
	Schemas() ([]string, error)
	Tables() ([]*TableInfo, error)
	Table(table string) (*Result, error)
	TableRows(table string, opts RowsOptions) (*Result, error)
	TableInfo(table string) (*Result, error)
	TableIndexes(table string) (*Result, error)
	Activity() (*Result, error)
	Query(query string) (*Result, error)
	History() []HistoryRecord
	GetCapabilities() ClientCapabilities
}

type ColumnInfo struct {
	Name                   string `json:"column_name"`
	DataType               string `json:"data_type"`
	IsNullable             string `json:"is_nullable"`
	CharacterMaximumLength string `json:"character_maximum_length"`
	CharacterSetCatalog    string `json:"character_set_catalog"`
	Default                string `json:"column_default"`
	// Extra                  map[string]interface{} `json:"extra"`	// TODO: implement general case also for frontend
}

type TableInfo struct {
	SchemaName string       `json:"table_schema"`
	TableName  string       `json:"table_name"`
	Columns    []ColumnInfo `json:"columns"`
}

// ClientCapabilities describes capabilities of the client.
// Front-end might customize the view depending on capabilities
type ClientCapabilities struct {
	// does it support query analyze
	HasAnalyze bool
}

// utility functions

func updateRow(row []interface{}) int {
	var size int
	for i, item := range row {
		if item == nil {
			row[i] = nil
			continue
		}
		switch v := item.(type) {
		default:
			LogInfof("unhandled type %T\n", item)
		case bool:
			size++
		case int, float32:
			size += 4
		case string:
			size += len(v)
		case []byte:
			size += len(v)
		case int64, float64:
			size += 8
		case time.Time:
			size += int(unsafe.Sizeof(v))
		}
		t := reflect.TypeOf(item).Kind().String()

		if t == "slice" {
			row[i] = string(item.([]byte))
		}
	}
	return size
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

func getTableInfo(tableInfos []*TableInfo, tableName string) *TableInfo {
	for _, tableInfo := range tableInfos {
		if tableInfo.TableName == tableName {
			return tableInfo
		}
	}
	return nil
}

func dbQueryTableInfo(db *sqlx.DB, query string) ([]*TableInfo, error) {
	var tableInfos []*TableInfo

	results, err := dbQuery(db, query)
	if err != nil {
		return tableInfos, err
	}

	reformatedData := make([]map[interface{}]interface{}, 0)

	for _, row := range results.Rows {
		m := make(map[interface{}]interface{})
		for j, col := range results.Columns {
			m[col] = row[j]
		}
		reformatedData = append(reformatedData, m)
	}

	for _, row := range reformatedData {
		LogInfof("row %v\n", row)

		var tableInfo = getTableInfo(tableInfos, row["table_name"].(string))

		LogInfof("current tableInfo %v\n", tableInfo)

		if tableInfo == nil {
			tableInfo = &TableInfo{}

			tableInfo.SchemaName = row["table_schema"].(string)
			tableInfo.TableName = row["table_name"].(string)
			tableInfos = append(tableInfos, tableInfo)
		}

		var column = ColumnInfo{}

		for columnType, value := range row {
			switch columnType {
			case "column_name":
				if row["column_name"] != nil {
					column.Name = value.(string)
				}
			case "data_type":
				if row["data_type"] != nil {
					column.DataType = row["data_type"].(string)
				}
			case "is_nullable":
				if row["is_nullable"] != nil {
					column.IsNullable = row["is_nullable"].(string)
				}
			case "character_maximum_length":
				if row["character_maximum_length"] != nil {
					// MySQL is string
					// PostgreSQL is int64
					vType := reflect.TypeOf(row["character_maximum_length"]).Kind()
					if vType == reflect.Int64 {
						column.CharacterMaximumLength = strconv.FormatInt(row["character_maximum_length"].(int64), 10)
					} else if vType == reflect.String {
						column.CharacterMaximumLength = row["character_maximum_length"].(string)
					} else {
						LogErrorf("This case shouldn't happen")
					}
				}
			case "character_set_catalog":
				if row["character_set_catalog"] != nil {
					column.CharacterSetCatalog = row["character_set_catalog"].(string)
				}
			case "column_default":
				if row["column_default"] != nil {
					column.Default = row["column_default"].(string)
				}
			default:
				// TODO: Implement General case
				// LogInfof("DEFAULT TYPE %v\n", columnType)
				// LogInfof("DEFAULT VALUE %v\n", value)
				// if columnType != nil {
				// 	if row[columnType] != nil {
				// 		LogInfof("DEFAULT extraColumnType %v\n", row[columnType])
				// 		column.Extra[columnType] = row[columnType]
				// 	}
				// }
			}
		}

		tableInfo.Columns = append(tableInfo.Columns, column)
	}

	LogInfof("tableInfos %v\n", tableInfos)

	return tableInfos, nil
}
