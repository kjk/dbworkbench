package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testClient Client
	isPg       bool
	isMysql    bool
)

func setupClient(t *testing.T) bool {
	var err error
	connURL := os.Getenv("DBHERO_TEST_CONN")
	fmt.Printf("connURL: '%s'\n", connURL)
	assert.NotEqual(t, "", connURL)
	if strings.HasPrefix(connURL, "postgres") {
		isPg = true
		testClient, err = connectPostgres(connURL)
	} else {
		isMysql = true
		testClient, err = connectMysql(connURL)
	}
	assert.NoError(t, err)
	return err == nil
}

func teardownClient() {
	if testClient != nil {
		testClient.Connection().Close()
	}
}

func testTest(t *testing.T) {
	//assert.Equal(t, nil, testClient.Test())
}

func testInfo(t *testing.T) {
	res, err := testClient.Info()
	assert.NoError(t, err)
	assert.NotNil(t, res)
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

	assert.NoError(t, err)
	assert.True(t, strInArray("world", res))
	if isPg {
		assert.True(t, strInArray("postgres", res))
	}
}

func testTables(t *testing.T) {
	res, err := testClient.Tables()

	expected := []string{
		"City",
		"Country",
		"CountryLanguage",
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}

func testTable(t *testing.T) {
	res, err := testClient.Table("City")

	columns := []string{
		"column_name",
		"data_type",
		"is_nullable",
		"character_maximum_length",
		"character_set_catalog",
		"column_default",
	}

	assert.NoError(t, err)
	assert.Equal(t, columns, res.Columns)
	assert.Equal(t, 5, len(res.Rows))
}

func testTableRows(t *testing.T) {
	res, err := testClient.TableRows("City", RowsOptions{})

	assert.NoError(t, err)
	assert.Equal(t, 5, len(res.Columns))
	assert.Equal(t, 4079, len(res.Rows))
}

func testTableInfo(t *testing.T) {
	res, err := testClient.TableInfo("City")

	assert.NoError(t, err)
	assert.Equal(t, 3, len(res.Columns))
	assert.Equal(t, 1, len(res.Rows))
}

func testTableIndexes(t *testing.T) {
	res, err := testClient.TableIndexes("City")
	//res.DumpFull()
	assert.NoError(t, err)
	assert.Equal(t, 13, len(res.Columns))
	assert.Equal(t, 2, len(res.Rows))
}

func testQuery(t *testing.T) {
	res, err := testClient.Query("SELECT * FROM City")
	assert.NoError(t, err)
	assert.Equal(t, 5, len(res.Columns))
	assert.Equal(t, 4079, len(res.Rows))
}

func strEq(s string, args ...string) bool {
	for _, arg := range args {
		if s == arg {
			return true
		}
	}
	return false
}

func testQueryError(t *testing.T) {
	res, err := testClient.Query("SELCT * FROM City")
	assert.Error(t, err)
	assert.Nil(t, res)
	cond := func() bool {
		s := err.Error()
		if s == "" {
			return true
		}
		return strings.HasPrefix(s, "Error 1064: You have an error in your SQL syntax")
	}
	assert.Condition(t, cond, err.Error())
}

func testQueryInvalidTable(t *testing.T) {
	res, err := testClient.Query("SELECT * FROM books2")
	assert.Error(t, err)
	assert.Nil(t, res)
	cond := func() bool {
		return strEq(err.Error(), `pq: relation "books2" does not exist`, `Error 1146: Table 'world.books2' doesn't exist`)
	}
	assert.Condition(t, cond, err.Error())
}

/*

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
*/

func TestAll(t *testing.T) {
	fmt.Printf("TestAll: started\n")
	if isWindows() {
		// Dont have access to windows machines at the moment...
		return
	}

	if !setupClient(t) {
		return
	}

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
	// TODO: update for world database
	/*
		testResultCsv(t)
		testHistory(t)
		testHistoryError(t)
	*/
	teardownClient()
}
