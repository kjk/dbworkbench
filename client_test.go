package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: needs to rename db name here and in data/booktown.sql to booktown_test
// so that I can have a test db for testing the program and be able to run tests

var testClient *Client
var testCommands map[string]string

func setupCommands() {
	testCommands = map[string]string{
		"createdb": "createdb",
		"psql":     "psql",
		"dropdb":   "dropdb",
	}

	if onWindows() {
		for k, v := range testCommands {
			testCommands[k] = v + ".exe"
		}
	}
}

func onWindows() bool {
	return runtime.GOOS == "windows"
}

func setup() {
	out, err := exec.Command(testCommands["createdb"], "-U", "postgres", "-h", "localhost", "booktown").CombinedOutput()

	if err != nil {
		fmt.Println("Database creation failed:", string(out))
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	out, err = exec.Command(testCommands["psql"], "-U", "postgres", "-h", "localhost", "-f", "./data/booktown.sql", "booktown").CombinedOutput()

	if err != nil {
		fmt.Println("Database import failed:", string(out))
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func setupClient() {
	testClient, _ = NewClientFromUrl("postgres://postgres@localhost/booktown?sslmode=disable")
}

func teardownClient() {
	if testClient != nil {
		testClient.db.Close()
	}
}

func teardown() {
	_, err := exec.Command(testCommands["dropdb"], "-U", "postgres", "-h", "localhost", "booktown").CombinedOutput()

	if err != nil {
		fmt.Println("Teardown error:", err)
	}
}

func testNewClientFromUrl(t *testing.T) {
	url := "postgres://postgres@localhost/booktown?sslmode=disable"
	client, err := NewClientFromUrl(url)

	if err != nil {
		defer client.db.Close()
	}

	assert.Equal(t, nil, err)
	assert.Equal(t, url, client.connectionString)
}

func testTest(t *testing.T) {
	assert.Equal(t, nil, testClient.Test())
}

func testInfo(t *testing.T) {
	res, err := testClient.Info()

	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, res)
}

func strInArray(s string, arr []string) bool {
	for _, s2 := range arr {
		if s2 == s {
			return true
		}
	}
	return false
}

func testDatabases(t *testing.T) {
	res, err := testClient.Databases()

	assert.Equal(t, nil, err)
	assert.True(t, strInArray("booktown", res))
	assert.True(t, strInArray("postgres", res))
}

func testTables(t *testing.T) {
	res, err := testClient.Tables()

	expected := []string{
		"alternate_stock",
		"authors",
		"book_backup",
		"book_queue",
		"books",
		"customers",
		"daily_inventory",
		"distinguished_authors",
		"editions",
		"employees",
		"favorite_authors",
		"favorite_books",
		"money_example",
		"my_list",
		"numeric_values",
		"publishers",
		"recent_shipments",
		"schedules",
		"shipments",
		"states",
		"stock",
		"stock_backup",
		"stock_view",
		"subjects",
		"text_sorting",
	}

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, res)
}

func testTable(t *testing.T) {
	res, err := testClient.Table("books")

	columns := []string{
		"column_name",
		"data_type",
		"is_nullable",
		"character_maximum_length",
		"character_set_catalog",
		"column_default",
	}

	assert.Equal(t, nil, err)
	assert.Equal(t, columns, res.Columns)
	assert.Equal(t, 4, len(res.Rows))
}

func testTableRows(t *testing.T) {
	res, err := testClient.TableRows("books", RowsOptions{})

	assert.Equal(t, nil, err)
	assert.Equal(t, 4, len(res.Columns))
	assert.Equal(t, 15, len(res.Rows))
}

func testTableInfo(t *testing.T) {
	res, err := testClient.TableInfo("books")

	assert.Equal(t, nil, err)
	assert.Equal(t, 4, len(res.Columns))
	assert.Equal(t, 1, len(res.Rows))
}

func testTableIndexes(t *testing.T) {
	res, err := testClient.TableIndexes("books")

	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(res.Columns))
	assert.Equal(t, 2, len(res.Rows))
}

func testQuery(t *testing.T) {
	res, err := testClient.Query("SELECT * FROM books")

	assert.Equal(t, nil, err)
	assert.Equal(t, 4, len(res.Columns))
	assert.Equal(t, 15, len(res.Rows))
}

func testQueryError(t *testing.T) {
	res, err := testClient.Query("SELCT * FROM books")

	assert.NotEqual(t, nil, err)
	assert.Equal(t, "pq: syntax error at or near \"SELCT\"", err.Error())
	assert.Equal(t, true, res == nil)
}

func testQueryInvalidTable(t *testing.T) {
	res, err := testClient.Query("SELECT * FROM books2")

	assert.NotEqual(t, nil, err)
	assert.Equal(t, "pq: relation \"books2\" does not exist", err.Error())
	assert.Equal(t, true, res == nil)
}

func testResultCsv(t *testing.T) {
	res, _ := testClient.Query("SELECT * FROM books ORDER BY id ASC LIMIT 1")
	csv := res.CSV()

	expected := "id,title,author_id,subject_id\n156,The Tell-Tale Heart,115,9\n"

	assert.Equal(t, expected, string(csv))
}

func testHistory(t *testing.T) {
	_, err := testClient.Query("SELECT * FROM books")
	query := testClient.history[len(testClient.history)-1].Query

	assert.Equal(t, nil, err)
	assert.Equal(t, "SELECT * FROM books", query)
}

func testHistoryError(t *testing.T) {
	_, err := testClient.Query("SELECT * FROM books123")
	query := testClient.history[len(testClient.history)-1].Query

	assert.NotEqual(t, nil, err)
	assert.NotEqual(t, "SELECT * FROM books123", query)
}

func TestAll(t *testing.T) {
	if onWindows() {
		// Dont have access to windows machines at the moment...
		return
	}

	setupCommands()
	teardown()
	setup()
	setupClient()

	testNewClientFromUrl(t)
	testTest(t)
	testInfo(t)
	testDatabases(t)
	testTables(t)
	testTable(t)
	testTableRows(t)
	testTableInfo(t)
	testTableIndexes(t)
	testQuery(t)
	testQueryError(t)
	testQueryInvalidTable(t)
	testResultCsv(t)
	testHistory(t)
	testHistoryError(t)

	teardownClient()
	teardown()
}
