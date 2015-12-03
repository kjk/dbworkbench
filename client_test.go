package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: needs to rename db name here and in data/world.sql to world_test
// so that I can have a test db for testing the program and be able to run tests

var testClient *Client

func onWindows() bool {
	return runtime.GOOS == "windows"
}

func exeName(cmd string) string {
	if onWindows() {
		return cmd + ".exe"
	}
	return cmd
}

func buildArgs() []string {
	pgHost := "localhost"
	pgPort := "5432"
	pgUser := "postgres"

	s := os.Getenv("WERCKER_POSTGRESQL_HOST")
	if s != "" {
		pgHost = s
	}

	s = os.Getenv("WERCKER_POSTGRESQL_PORT")
	if s != "" {
		pgPort = s
	}

	s = os.Getenv("WERCKER_POSTGRESQL_USERNAME")
	if s != "" {
		pgUser = s
	}

	res := []string{
		"-U", pgUser,
		"-h", pgHost,
		"-p", pgPort,
	}
	return res
}

func dbURL() string {
	host := os.Getenv("WERCKER_POSTGRESQL_HOST")
	if host == "" {
		return "postgres://postgres@localhost/world?sslmode=disable"
	}
	port := os.Getenv("WERCKER_POSTGRESQL_PORT")
	user := os.Getenv("WERCKER_POSTGRESQL_USERNAME")
	pwd := os.Getenv("WERCKER_POSTGRESQL_PASSWORD")
	s := fmt.Sprintf("postgres://%s:%s@%s:%s/world?sslmode=disable", user, pwd, host, port)
	return s
}

func setupCmdPgPassword(cmd *exec.Cmd) {
	pwd := os.Getenv("WERCKER_POSTGRESQL_PASSWORD")
	if pwd != "" {
		cmd.Env = append(cmd.Env, "PGPASSWORD="+pwd)
	}
}

func setup() {
	args := buildArgs()
	args = append(args, "world")
	cmd := exec.Command(exeName("createdb"), args...)
	setupCmdPgPassword(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Database creation failed:", string(out))
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	args = buildArgs()
	args = append(args, "-f", "./data/world.sql", "world")
	cmd = exec.Command(exeName("psql"), args...)
	setupCmdPgPassword(cmd)
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Database import failed:", string(out))
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func setupClient() {
	testClient, _ = NewClientFromUrl(dbURL())
}

func teardownClient() {
	if testClient != nil {
		testClient.db.Close()
	}
}

func teardown() {
	args := buildArgs()
	args = append(args, "world")
	cmd := exec.Command(exeName("dropdb"), args...)
	setupCmdPgPassword(cmd)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Teardown error:", err)
	}
}

func testNewClientFromURL(t *testing.T) {
	url := dbURL()
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
	assert.True(t, strInArray("world", res))
	assert.True(t, strInArray("postgres", res))
}

func testTables(t *testing.T) {
	res, err := testClient.Tables()

	expected := []string{
		"city",
		"country",
		"countrylanguage",
	}

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, res)
}

/*
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
*/

func TestAll(t *testing.T) {
	if onWindows() {
		// Dont have access to windows machines at the moment...
		return
	}

	teardown()
	setup()
	setupClient()

	testNewClientFromURL(t)
	testTest(t)
	testInfo(t)
	testDatabases(t)
	testTables(t)
	// TODO: update for world database
	/*
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
	*/
	teardownClient()
	teardown()
}
