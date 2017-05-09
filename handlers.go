package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// HandlerWithCtxFunc is like http.HandlerFunc but with additional ReqContext argument
type HandlerWithCtxFunc func(*ReqContext, http.ResponseWriter, *http.Request)

// ReqOpts is a set of flags passed to withCtx
type ReqOpts uint

const (
	// OnlyGet tells to reject non-GET requests
	OnlyGet ReqOpts = 1 << iota
	// OnlyPost tells to reject non-POST requests
	OnlyPost
	// MustHaveConnection tells to reject requests if conn_id request arg is not
	// present
	MustHaveConnection
	// IsJSON denotes a handler that is serving JSON requests and should send
	// errors as { "error": "error message" }
	IsJSON

	defaultMaxDataSize = 100 * 1024 * 1024 // 100 MB
)

// ReqContext contains data that is useful to access in every http handler
type ReqContext struct {
	TimeStart    time.Time
	ConnectionID int
	ConnInfo     *ConnectionInfo
	Client       Client
}

var (
	// loaded only once at startup. maps a file path of the resource
	// to its data
	resourcesFromZip map[string][]byte
)

func withCtx(f HandlerWithCtxFunc, opts ReqOpts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isJSON := opts&IsJSON != 0

		cw := NewCountingResponseWriter(w)
		method := strings.ToUpper(r.Method)
		if opts&OnlyGet != 0 {
			if method != "GET" {
				serveError(w, r, isJSON, fmt.Sprintf("%s is not GET", r.Method))
				return
			}
		}
		if opts&OnlyPost != 0 {
			if method != "POST" {
				serveError(w, r, isJSON, fmt.Sprintf("%s is not POST", r.Method))
				return
			}
		}

		ctx := &ReqContext{
			TimeStart: time.Now(),
		}

		if opts&MustHaveConnection != 0 {
			connID, err := strconv.Atoi(r.FormValue("conn_id"))
			if err != nil || connID <= 0 {
				var errMsg string
				if err != nil {
					errMsg = err.Error()
				} else {
					errMsg = fmt.Sprintf("connID: %d (should be > 0)", connID)
				}
				serveError(w, r, isJSON, errMsg)
				return
			}
			connInfo := getConnectionInfoByID(connID)
			if nil == connInfo {
				errMsg := fmt.Sprintf("invalid conn_id %d", connID)
				serveError(w, r, isJSON, errMsg)
				return
			}
			ctx.ConnectionID = connID
			ctx.ConnInfo = connInfo
			ctx.Client = connInfo.Client
		}

		f(ctx, cw, r)
		if !strings.HasPrefix(r.RequestURI, "/s/") {
			LogInfof("%s took %s, code: %d\n", r.RequestURI, time.Since(ctx.TimeStart), cw.Code)
		}
	}
}

func normalizePath(s string) string {
	return strings.Replace(s, "\\", "/", -1)
}

func loadResourcesFromZipReader(zr *zip.Reader) error {
	for _, f := range zr.File {
		name := normalizePath(f.Name)
		rc, err := f.Open()
		if err != nil {
			return err
		}
		d, err := ioutil.ReadAll(rc)
		rc.Close()
		if err != nil {
			return err
		}
		// for simplicity of the build, the file that we embedded in zip
		// is bundle.min.js but the html refers to it as bundle.js
		if name == "s/dist/bundle.min.js" {
			name = "s/dist/bundle.js"
		}
		//LogInfof("Loaded '%s' of size %d bytes\n", name, len(d))
		resourcesFromZip[name] = d
	}
	return nil
}

// call this only once at startup
func loadResourcesFromZip(path string) error {
	resourcesFromZip = make(map[string][]byte)
	zrc, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer zrc.Close()
	return loadResourcesFromZipReader(&zrc.Reader)
}

func loadResourcesFromEmbeddedZip() error {
	//LogInfof("loadResourcesFromEmbeddedZip()\n")
	n := len(resourcesZipData)
	if n == 0 {
		return errors.New("len(resourcesZipData) == 0")
	}
	resourcesFromZip = make(map[string][]byte)
	r := bytes.NewReader(resourcesZipData)
	zrc, err := zip.NewReader(r, int64(n))
	if err != nil {
		return err
	}
	return loadResourcesFromZipReader(zrc)
}

func serveResourceFromZip(w http.ResponseWriter, r *http.Request, path string) {
	path = normalizePath(path)
	LogInfof("serving '%s' from zip\n", path)

	data := resourcesFromZip[path]

	if data == nil {
		LogErrorf("no data for file '%s'\n", path)
		servePlainText(w, r, 404, fmt.Sprintf("file '%s' not found", path))
		return
	}

	if len(data) == 0 {
		servePlainText(w, r, 404, "Asset is empty")
		return
	}

	serveData(w, r, 200, MimeTypeByExtensionExt(path), data)
}

func serveStatic(w http.ResponseWriter, r *http.Request, path string) {
	if options.ResourcesFromZip {
		serveResourceFromZip(w, r, path)
		return
	}

	data, err := ioutil.ReadFile(path)

	if err != nil {
		LogErrorf("ioutil.ReadFile('%s') failed with '%s'\n", path, err)
		servePlainText(w, r, 404, err.Error())
		return
	}

	if len(data) == 0 {
		servePlainText(w, r, 404, "Asset is empty")
		return
	}

	serveData(w, r, 200, MimeTypeByExtensionExt(path), data)
}

// GET /
func handleIndex(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	if uri == "/" {
		uri = "index.html"
	} else {
		uri = strings.ToLower(uri)
		if strings.HasSuffix(uri, ".html") {
			uri = uri[1:]
		}
	}
	path := filepath.Join("s", uri)
	serveStatic(w, r, path)
}

// GET /s/:path
func handleStatic(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:] // remove initial "/" i.e. "/s/*" => "s/*"
	//LogInfof("path='%s'\n", path)
	serveStatic(w, r, path)
}

// POST /api/connect
// args:
//	url     : database connection url formatted for Go driver
//  urlSafe : like url but with password replaced with ***
//  type    : database type ('postgres' or 'mysql')
//
func handleConnect(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	url := strings.TrimSpace(r.FormValue("url"))
	urlSafe := strings.TrimSpace(r.FormValue("urlSafe"))
	dbType := strings.TrimSpace(r.FormValue("type"))
	LogInfof("dbtype: '%s' urlSafe: '%s'\n", dbType, urlSafe)
	if url == "" || urlSafe == "" || dbType == "" {
		serveJSONError(w, r, fmt.Errorf("url or urlSafe ('%s') or type ('%s') argument is missing", urlSafe, dbType))
		return
	}
	var client Client
	var err error
	switch dbType {
	case dbTypePostgres:
		client, err = connectPostgres(url)
	case dbTypeMysql:
		client, err = connectMysql(url)
	default:
		err = fmt.Errorf("invalid 'type' argument ('%s')", dbType)
	}

	if err != nil {
		msg := strings.Replace(err.Error(), url, urlSafe, -1)
		LogErrorf("failed to connect with '%s'\n", msg)
		serveJSONError(w, r, msg)
		return
	}

	recordDatabaseOpened()

	info, err := client.Info()
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	connInfo := addConnectionInfo(url, client)

	i := info.Format()[0]
	currDb, ok := i["current_database"]
	if !ok {
		serveJSONError(w, r, "must provide database")
		return
	}
	currDbStr, ok := currDb.(string)
	if !ok {
		serveJSONError(w, r, "must provide a database")
		return
	}
	v := struct {
		ConnectionID    int
		CurrentDatabase string
		Capabilities    ClientCapabilities
	}{
		ConnectionID:    connInfo.ConnectionID,
		CurrentDatabase: currDbStr,
		Capabilities:    client.GetCapabilities(),
	}
	serveJSON(w, r, v)
}

// POST /api/disconnect
func handleDisconnect(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	LogInfof("conn id='%d'\n", ctx.ConnectionID)

	err := connectionDisconnect(ctx.ConnectionID)
	if err != nil {
		serveJSONError(w, r, err)
		return
	}
	v := struct {
		Message string
	}{
		Message: "ok",
	}
	serveJSON(w, r, v)
}

// GET /api/databases
func handleGetDatabases(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	updateConnectionLastAccess(ctx.ConnInfo.ConnectionID)
	names, err := ctx.ConnInfo.Client.Databases()
	if err != nil {
		serveJSONError(w, r, err)
		return
	}
	serveJSON(w, r, names)
}

func runQuery(ctx *ReqContext, w http.ResponseWriter, r *http.Request, query string) {
	updateConnectionLastAccess(ctx.ConnInfo.ConnectionID)
	result, err := ctx.ConnInfo.Client.Query(query)

	if err != nil {
		LogErrorf("query '%s' failed with '%s'\n", query, err)
		serveJSONError(w, r, err)
		return
	}

	q := r.URL.Query()

	if len(q["format"]) > 0 && q["format"][0] == "csv" {
		// TODO: add database name
		filename := fmt.Sprintf("db-%v.csv", time.Now().Unix())
		w.Header().Set("Content-disposition", "attachment;filename="+filename)
		serveData(w, r, 200, "text/csv", result.CSV())
		return
	}

	serveJSON(w, r, result)
}

/*
GET | POST /api/query
args:
    query : string, query to execute
returns; json in the format
{
    "columns" : ["id", "name", ...], // names of database columns
    "rows" : [
        [ , , ... ], // data for row one
        [ , , ... ] // data for row two etc.
    ]
}
*/
func handleQuery(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.FormValue("query"))
	//LogInfof("query: '%s'\n", query)

	if query == "" {
		serveJSONError(w, r, "Query parameter is missing")
		return
	}

	recordQueryExecuted()
	runQuery(ctx, w, r, query)
}

// QueryAsyncStatus describes that status of the (possibly stil executing) query
type QueryAsyncStatus struct {
	QueryID   string
	RowsCount int   `json:"rows_count"`
	DataSize  int64 `json:"data_size"`
	Finished  bool  `json:"finished"`
	// if NotComplete is true, we didn't get all data because
	// we reached max rows or max data limit
	NotComplete          bool     `json:"not_complete"`
	TimeToFirstResultMs  float64  `json:"time_to_first_result_ms"`
	TotalQueryTimeMs     float64  `json:"total_query_time_ms"`
	Columns              []string `json:"columns"`
	ErrorString          string   `json:"error"`
	err                  error
	rows                 []interface{}
	clientLastAccessTime time.Time // so that we can delete stale data after a while
}

var (
	// ID is an int but it's more convenient to treat it as a string
	currQueryAsyncID     int
	queryAsyncIDToStatus map[string]*QueryAsyncStatus
	muQueryAsync         sync.Mutex
)

func init() {
	queryAsyncIDToStatus = make(map[string]*QueryAsyncStatus)
}

// TODO: call this every once in a while
func gcStaleAsyncQueries() {
	muQueryAsync.Lock()
	defer muQueryAsync.Unlock()
	// I'm not sure what happens when we delete key while iterating
	// a map, so first collect keys to delete and then delete
	var toDelete []string
	for queryID, status := range queryAsyncIDToStatus {
		if time.Now().Sub(status.clientLastAccessTime) > time.Hour {
			toDelete = append(toDelete, queryID)
		}
	}
	for _, queryID := range toDelete {
		delete(queryAsyncIDToStatus, queryID)
	}
}

func getNextQueryAsyncID() string {
	muQueryAsync.Lock()
	currQueryAsyncID++
	id := strconv.Itoa(currQueryAsyncID)
	queryAsyncIDToStatus[id] = &QueryAsyncStatus{
		QueryID: id,
	}
	muQueryAsync.Unlock()
	return id
}

func withQueryStatus(queryID string, f func(s *QueryAsyncStatus)) error {
	muQueryAsync.Lock()
	defer muQueryAsync.Unlock()
	status := queryAsyncIDToStatus[queryID]
	if status == nil {
		LogErrorf("non-existent queryID '%s'\n", queryID)
		return fmt.Errorf("non-existent queryID '%s'\n", queryID)
	}
	f(status)
	return nil
}

func setQueryStatusError(queryID string, err error) {
	withQueryStatus(queryID, func(s *QueryAsyncStatus) {
		s.err = err
		if err != nil {
			s.ErrorString = err.Error()
		}
		s.Finished = true
	})
}

func durToMs(d time.Duration) float64 {
	return float64(d / time.Millisecond)
}

func durMsSince(t time.Time) float64 {
	return durToMs(time.Now().Sub(t))
}

func doQueryAsync(client Client, query string, queryID string, maxRows int, maxDataSize int64) {
	db := client.Connection()

	rows, err := db.Queryx(query)

	if err != nil {
		setQueryStatusError(queryID, err)
		return
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		setQueryStatusError(queryID, err)
		return
	}

	withQueryStatus(queryID, func(s *QueryAsyncStatus) {
		s.Columns = cols
	})

	client.AddToHistory(query)
	firstRow := false
	timeStart := time.Now()
	nRows := 0
	var dataSize int64
	for rows.Next() {
		row, err := rows.SliceScan()
		if err != nil {
			setQueryStatusError(queryID, err)
			return
		}
		if firstRow {
			withQueryStatus(queryID, func(s *QueryAsyncStatus) {
				s.TimeToFirstResultMs = durMsSince(timeStart)
			})
		}

		dataSize += int64(updateRow(row))
		if maxDataSize != -1 && dataSize > maxDataSize {
			withQueryStatus(queryID, func(s *QueryAsyncStatus) {
				s.NotComplete = true
			})
			break
		}

		withQueryStatus(queryID, func(s *QueryAsyncStatus) {
			s.rows = append(s.rows, row)
			s.TotalQueryTimeMs = durMsSince(timeStart)
			s.DataSize = dataSize
		})
		nRows++
		if maxRows != -1 && nRows >= maxRows {
			withQueryStatus(queryID, func(s *QueryAsyncStatus) {
				s.NotComplete = true
			})
			break
		}
	}

	withQueryStatus(queryID, func(s *QueryAsyncStatus) {
		s.Finished = true
		s.RowsCount = len(s.rows)
	})

}

/*
GET | POST /api/queryasync
args:
   conn_id       : string, connection
   query         : string, query to execute
   max_rows      : int, optional, max number of rows to fetch from the database
   max_data_size : int64, optional, max amount of data to fetch from the database

Both max_rows and max_data_size can be given, both are respected.

If max_rows and max_data_size are not given, we default to max_data_size of 100 MB.

returns: json in the format
{
  "query_id": 15
}
*/
func handleQueryAsync(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	var err error
	query := strings.TrimSpace(r.FormValue("query"))
	if query == "" {
		LogErrorf("Query parameter is missing, uri: '%s'\n", r.URL.RequestURI())
		serveJSONError(w, r, "Query parameter is missing")
		return
	}

	maxRowsStr := strings.TrimSpace(r.FormValue("max_rows"))
	maxRows := -1
	if maxRowsStr != "" {
		maxRows, err = strconv.Atoi(maxRowsStr)
		if err != nil {
			err = fmt.Errorf("invalid 'max_rows': '%s', err: '%s'", maxRowsStr, err)
			LogErrorf("%s\n", err)
			serveJSONError(w, r, err.Error())
			return
		}
	}
	maxDataSizeStr := strings.TrimSpace(r.FormValue("max_data_size"))
	var maxDataSize int64 = -1
	if maxDataSizeStr != "" {
		maxDataSize, err = strconv.ParseInt(maxDataSizeStr, 10, 64)
		if err != nil {
			err = fmt.Errorf("invalid 'max_data_size': '%s', err: '%s'", maxDataSizeStr, err)
			LogErrorf("%s\n", err)
			serveJSONError(w, r, err.Error())
			return
		}
	}
	if maxDataSize == -1 && maxRows == -1 {
		maxDataSize = defaultMaxDataSize
	}

	LogInfof("query: '%s', max_rows: '%s', max_data_size: '%s'\n", query, maxRowsStr, maxDataSizeStr)

	connInfo := ctx.ConnInfo
	queryID := getNextQueryAsyncID()

	go doQueryAsync(connInfo.Client, query, queryID, maxRows, maxDataSize)
	updateConnectionLastAccess(connInfo.ConnectionID)

	res := struct {
		QueryID string `json:"query_id"`
	}{
		QueryID: queryID,
	}
	serveJSON(w, r, res)
}

func getQueryStatus(queryID string) (*QueryAsyncStatus, error) {
	var statusCopy QueryAsyncStatus
	err := withQueryStatus(queryID, func(s *QueryAsyncStatus) {
		statusCopy = *s
	})
	return &statusCopy, err
}

/*
GET | POST /api/queryasyncstatus
args:
  conn_id  : string, connection
  query_id : string, id of the query

returns: json in the format
{
  "rows_count": 123,
  "data_size": 8234, // in bytes, approximate, we calc by adding size of all values
  "finished": false,
  "time_to_first_result_ms": 34, // time it took to get the first result from the server
  "total_query_time_ms": 1234, // time it took to get all results
  "columns": [ "id", "name" ],
  "error": null, // string if there was an error
}
*/
func handleQueryAsyncStatus(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	queryID := strings.TrimSpace(r.FormValue("query_id"))
	LogInfof("query_id: '%s'\n", queryID)
	status, err := getQueryStatus(queryID)
	if err != nil {
		serveJSONError(w, r, err)
	} else {
		serveJSON(w, r, status)
	}
}

/*
GET | POST /api/queryasyncdata
args:
  conn_id  : string, connection
  query_id : string
  start    : int, first row to return
  count    : int, number of rows, start + count should be < total rows count
result: json in the format
{
    "rows": [
         [ ... ],
         [ ... ],
     ]
}
*/
func handleQueryAsyncData(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	queryID := strings.TrimSpace(r.FormValue("query_id"))
	startStr := strings.TrimSpace(r.FormValue("start"))
	start, err := strconv.Atoi(startStr)
	if err != nil {
		LogErrorf("invalid 'start' argument '%s'", startStr)
		serveJSONError(w, r, err)
		return
	}
	countStr := strings.TrimSpace(r.FormValue("count"))
	count, err := strconv.Atoi(countStr)
	if err != nil {
		LogErrorf("invalid 'start' argument '%s'", countStr)
		serveJSONError(w, r, err)
		return
	}

	LogInfof("query_id: '%s', start: %d, count: %d \n", queryID, start, count)

	// hold the lock for the shortest possible time
	muQueryAsync.Lock()
	s := queryAsyncIDToStatus[queryID]
	if s == nil {
		muQueryAsync.Unlock()
		err = fmt.Errorf("invalid query_id '%s'", queryID)
		LogErrorf("%s\n", err.Error())
		serveJSONError(w, r, err)
		return
	}

	// start is 0-based
	end := start + count
	rowsCount := s.RowsCount
	if end > rowsCount {
		muQueryAsync.Unlock()
		err = fmt.Errorf("start+count (%d) too high (max is %d)'", end, rowsCount)
		LogErrorf("%s\n", err.Error())
		serveJSONError(w, r, err)
		return
	}
	// make a copy of the results
	rows := make([]interface{}, count, count)
	for i := 0; i < count; i++ {
		rows[i] = s.rows[start+i]
	}
	muQueryAsync.Unlock()
	res := struct {
		Rows []interface{} `json:"rows"`
	}{
		Rows: rows,
	}
	serveJSON(w, r, res)
}

/*
GET | POST /api/explain
args:
  conn_id : string, connection
  query   : string, query to explain
*/
func handleExplain(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.FormValue("query"))
	//LogInfof("query: '%s'\n", query)

	if query == "" {
		serveJSONError(w, r, "Query parameter is missing")
		return
	}

	runQuery(ctx, w, r, fmt.Sprintf("EXPLAIN ANALYZE %s", query))
}

// GET /api/history
func handleHistory(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	serveJSON(w, r, ctx.ConnInfo.Client.GetHistory())
}

// GET /api/getbookmarks
func handleGetBookmarks(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	jsonp := strings.TrimSpace(r.FormValue("jsonp"))

	bookmarks, err := readBookmarks()
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	sortedBookmarks := sortBookmarks(bookmarks)

	serveJSONP(w, r, sortedBookmarks, jsonp)
}

func getFormInt(r *http.Request, name string) (int, error) {
	s := strings.TrimSpace(r.FormValue(name))
	if s == "" {
		return 0, fmt.Errorf("missing form value '%s'", name)
	}
	return strconv.Atoi(s)
}

func getFormBool(r *http.Request, name string) (bool, error) {
	s := strings.ToLower(strings.TrimSpace(r.FormValue(name)))
	switch s {
	case "1", "true":
		return true, nil
	case "0", "false":
		return false, nil
	}
	return false, fmt.Errorf("invalid bool value '%s' for form key '%s'", r.FormValue(name), name)
}

// POST /api/addbookmark
func handleAddBookmark(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	id, err := getFormInt(r, "id")
	newBookmark := Bookmark{
		ID:       id,
		Nick:     r.FormValue("nick"),
		Type:     r.FormValue("type"),
		Host:     r.FormValue("host"),
		Port:     r.FormValue("port"),
		User:     r.FormValue("user"),
		Password: r.FormValue("password"),
		Database: r.FormValue("database"),
	}

	// TODO: validate fields make sense (type is pg or mysql, nick is not empty)
	bookmarks, err := addBookmark(newBookmark)
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	sortedBookmarks := sortBookmarks(bookmarks)

	serveJSON(w, r, sortedBookmarks)
}

// POST /api/removebookmark
func handleRemoveBookmark(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	id, err := getFormInt(r, "id")
	if err != nil {
		serveJSONError(w, r, err)
		return
	}
	bookmarks, err := removeBookmark(id)
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	sortedBookmarks := sortBookmarks(bookmarks)
	serveJSON(w, r, sortedBookmarks)
}

// GET /api/connection
func handleConnectionInfo(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	updateConnectionLastAccess(ctx.ConnInfo.ConnectionID)
	res, err := ctx.ConnInfo.Client.Info()
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	serveJSON(w, r, res.Format()[0])
}

// GET /api/activity
func handleActivity(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	updateConnectionLastAccess(ctx.ConnInfo.ConnectionID)
	res, err := ctx.ConnInfo.Client.Activity()
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	serveJSON(w, r, res)
}

/*
GET /api/schemas
Note: not used by frontend
*/
func handleGetSchemas(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	updateConnectionLastAccess(ctx.ConnInfo.ConnectionID)
	names, err := ctx.ConnInfo.Client.Schemas()
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	serveJSON(w, r, names)
}

// GET /api/tables
func handleGetTables(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	//LogInfof("connID: %d\n", ctx.ConnectionID)
	updateConnectionLastAccess(ctx.ConnInfo.ConnectionID)
	names, err := ctx.ConnInfo.Client.Tables()
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	serveJSON(w, r, names)
}

func handleGetTable(ctx *ReqContext, w http.ResponseWriter, r *http.Request, table string) {
	updateConnectionLastAccess(ctx.ConnInfo.ConnectionID)
	res, err := ctx.ConnInfo.Client.Table(table)
	//LogInfof("table: '%s'\n", table)
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	serveJSON(w, r, res)
}

func apiGetTableInfo(ctx *ReqContext, w http.ResponseWriter, r *http.Request, table string) {
	res, err := ctx.ConnInfo.Client.TableInfo(table)
	if err != nil {
		serveJSONError(w, r, err)
		return
	}
	serveJSON(w, r, res.Format()[0])
}

func apiGetTableIndexes(ctx *ReqContext, w http.ResponseWriter, r *http.Request, table string) {
	LogInfof("table='%s'\n", table)
	updateConnectionLastAccess(ctx.ConnInfo.ConnectionID)
	res, err := ctx.ConnInfo.Client.TableIndexes(table)
	if err != nil {
		serveJSONError(w, r, err)
		return
	}
	serveJSON(w, r, res)
}

/*
GET /api/tables/:table/:action
args:
  action : "info", "indexes"
*/

func handleTablesDispatch(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	uriPath := uri[len("/api/tables/"):]
	parts := strings.SplitN(uriPath, "/", 2)
	table := parts[0]
	LogInfof("table='%s'\n", table)
	if len(parts) == 1 {
		handleGetTable(ctx, w, r, table)
		return
	}
	cmd := parts[1]
	switch cmd {
	case "info":
		apiGetTableInfo(ctx, w, r, table)
	case "indexes":
		apiGetTableIndexes(ctx, w, r, table)
	default:
		LogErrorf("unknown cmd: '%s'\n", cmd)
		http.NotFound(w, r)
	}
}

/*
GET /api/userinfo
args:
  jsonp : jsonp wrapper, optional
Returns information about the user
*/
func handleUserInfo(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	jsonp := strings.TrimSpace(r.FormValue("jsonp"))
	//LogInfof("User: %#v\n", ctx.User)

	v := struct {
		ConnectionID int
	}{
		ConnectionID: 0,
	}
	connID := getFirstConnectionID()
	if -1 != connID {
		v.ConnectionID = connID
	}
	LogVerbosef("v: %#v\n", v)
	serveJSONP(w, r, v, jsonp)
}

/*
GET /api/launchbrowser
args:
  url : url to open in the default browser
*/
func handleLaunchBrowser(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	url := strings.TrimSpace(r.FormValue("url"))
	LogInfof("url: '%s'\n", url)

	if url == "" {
		serveJSONError(w, r, "'url' parameter is missing")
		return
	}

	err := openDefaultBrowser(url)
	if err != nil {
		LogError(err.Error())
	}

	v := struct {
	}{}
	serveJSON(w, r, v)
}

func registerHTTPHandlers() {
	http.HandleFunc("/", withCtx(handleIndex, OnlyGet))
	http.HandleFunc("/s/", withCtx(handleStatic, OnlyGet))

	http.HandleFunc("/api/activity", withCtx(handleActivity, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/getbookmarks", withCtx(handleGetBookmarks, IsJSON))
	http.HandleFunc("/api/addbookmark", withCtx(handleAddBookmark, OnlyPost|IsJSON))
	http.HandleFunc("/api/removebookmark", withCtx(handleRemoveBookmark, OnlyPost|IsJSON))
	http.HandleFunc("/api/connect", withCtx(handleConnect, OnlyPost|IsJSON))
	http.HandleFunc("/api/connection", withCtx(handleConnectionInfo, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/databases", withCtx(handleGetDatabases, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/disconnect", withCtx(handleDisconnect, OnlyPost|MustHaveConnection|IsJSON))
	http.HandleFunc("/api/explain", withCtx(handleExplain, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/history", withCtx(handleHistory, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/schemas", withCtx(handleGetSchemas, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/tables", withCtx(handleGetTables, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/tables/", withCtx(handleTablesDispatch, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/query", withCtx(handleQuery, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/queryasync", withCtx(handleQueryAsync, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/queryasyncstatus", withCtx(handleQueryAsyncStatus, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/queryasyncdata", withCtx(handleQueryAsyncData, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/userinfo", withCtx(handleUserInfo, IsJSON))
	http.HandleFunc("/api/launchbrowser", withCtx(handleLaunchBrowser, IsJSON))
	http.HandleFunc("/showmyhost", handleShowMyHost)
}

// GET /showmyhost, for testing only
func handleShowMyHost(w http.ResponseWriter, r *http.Request) {
	s := getMyHost(r)
	servePlainText(w, r, 200, "me: %s\n", s)
}

func startWebServer() {
	registerHTTPHandlers()
	httpAddr := fmt.Sprintf("%s:%v", options.HTTPHost, options.HTTPPort)
	fmt.Printf("Started running on %s, dev mode: %v\n", httpAddr, options.IsDev)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		log.Fatalf("http.ListenAndServe() failed with %s\n", err)
	}
	fmt.Printf("Exited\n")
}
