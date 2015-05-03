package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var extraMimeTypes = map[string]string{
	".icon": "image-x-icon",
	".ttf":  "application/x-font-ttf",
	".woff": "application/x-font-woff",
	".eot":  "application/vnd.ms-fontobject",
	".svg":  "image/svg+xml",
}

// Error is error message sent to client as json response for backend errors
type Error struct {
	Message string `json:"error"`
}

// NewError creates a new Error
func NewError(err error) Error {
	return Error{err.Error()}
}

func getClientIP(r *http.Request) string {
	clientIP := r.Header.Get("X-Real-IP")
	if len(clientIP) > 0 {
		return clientIP
	}
	clientIP = r.Header.Get("X-Forwarded-For")
	clientIP = strings.Split(clientIP, ",")[0]
	if len(clientIP) > 0 {
		return strings.TrimSpace(clientIP)
	}
	return r.RemoteAddr
}

func assetContentType(name string) string {
	ext := filepath.Ext(name)
	result := mime.TypeByExtension(ext)

	if result == "" {
		result = extraMimeTypes[ext]
	}

	if result == "" {
		result = "text/plain; charset=utf-8"
	}

	return result
}

func asset(fileName string) ([]byte, error) {
	//fmt.Fprintf(os.Stderr, "asset: %s\n", fileName)
	return ioutil.ReadFile(fileName)
}

// TODO: not sure if it's worth to put GET, POST etc. filters
func get(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.NotFound(w, r)
			return
		}
		f(w, r)
	}
}

func post(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		f(w, r)
	}
}

func timeit(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		f(w, r)
		// TODO: save this for further analysis
		LogInfof("%s took %s\n", r.RequestURI, time.Since(timeStart))
	}
}

func hasUser(f http.HandlerFunc) http.HandlerFunc {
	// TODO: reject if no user
	return f
}

func hasConnection(f http.HandlerFunc) http.HandlerFunc {
	// TODO: reject if no connection
	return f
}
func serveStatic(w http.ResponseWriter, r *http.Request, path string) {
	data, err := asset(path)

	if err != nil {
		LogErrorf("asset('%s') failed with '%s'\n", path, err)
		serveString(w, r, 404, err.Error())
		return
	}

	if len(data) == 0 {
		serveString(w, r, 404, "Asset is empty")
		return
	}

	serveData(w, r, 200, assetContentType(path), data)
}

// GET /
func handleIndex(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	if uri != "/" {
		http.NotFound(w, r)
		return
	}
	serveStatic(w, r, "s/index.html")
}

// GET /s/:path
func handleStatic(w http.ResponseWriter, r *http.Request) {
	path := "s/" + r.URL.Path[len("/s/"):]
	//LogInfof("path='%s'\n", path)
	serveStatic(w, r, path)
}

func writeHeader(w http.ResponseWriter, code int, contentType string) {
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	w.WriteHeader(code)
}

func serveJSON(w http.ResponseWriter, r *http.Request, code int, data ...interface{}) error {
	writeHeader(w, code, "application/json")
	encoder := json.NewEncoder(w)
	return encoder.Encode(data[0])
}

func serveString(w http.ResponseWriter, r *http.Request, code int, format string, args ...interface{}) error {
	writeHeader(w, code, "text/plain")
	var err error
	if len(args) > 0 {
		_, err = w.Write([]byte(fmt.Sprintf(format, args...)))
	} else {
		_, err = w.Write([]byte(format))
	}
	return err
}

func serveData(w http.ResponseWriter, r *http.Request, code int, contentType string, data []byte) {
	if len(contentType) > 0 {
		w.Header().Set("Content-Type", contentType)
	}
	w.WriteHeader(code)
	w.Write(data)
}

// POST /api/connect
func handleConnect(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		serveJSON(w, r, 400, Error{"Url parameter is required"})
		return
	}

	opts := Options{URL: url}
	url, err := formatConnectionUrl(opts)

	if err != nil {
		serveJSON(w, r, 400, Error{err.Error()})
		return
	}

	client, err := NewClientFromUrl(url)
	if err != nil {
		serveJSON(w, r, 400, Error{err.Error()})
		return
	}

	err = client.Test()
	if err != nil {
		serveJSON(w, r, 400, Error{err.Error()})
		return
	}

	info, err := client.Info()

	if err == nil {
		if dbClient != nil {
			dbClient.db.Close()
		}

		dbClient = client
	}

	serveJSON(w, r, 200, info.Format()[0])
}

// GET /api/databases
func handleGetDatabases(w http.ResponseWriter, r *http.Request) {
	names, err := dbClient.Databases()

	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, names)
}

func handleQuery(w http.ResponseWriter, r *http.Request, query string) {
	result, err := dbClient.Query(query)

	if err != nil {
		LogErrorf("query: '%s', err: %s\n", query, err)
		serveJSON(w, r, 400, NewError(err))
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

	serveJSON(w, r, 200, result)
}

// GET | POST /api/query
func handleRunQuery(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.FormValue("query"))
	LogInfof("query: '%s'\n", query)

	if query == "" {
		serveJSON(w, r, 400, errors.New("Query parameter is missing"))
		return
	}

	handleQuery(w, r, query)
}

// GET | POST /api/explain
func handleExplainQuery(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.FormValue("query"))
	LogInfof("query: '%s'\n", query)

	if query == "" {
		serveJSON(w, r, 400, errors.New("Query parameter is missing"))
		return
	}

	handleQuery(w, r, fmt.Sprintf("EXPLAIN ANALYZE %s", query))
}

// GET /api/history
func handleHistory(w http.ResponseWriter, r *http.Request) {
	serveJSON(w, r, 200, dbClient.history)
}

// GET /api/bookmarks
func handleBookmarks(w http.ResponseWriter, r *http.Request) {
	bookmarks, err := readAllBookmarks(bookmarksPath())

	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, bookmarks)
}

// GET /api/connection
func handleConnectionInfo(w http.ResponseWriter, r *http.Request) {
	res, err := dbClient.Info()

	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, res.Format()[0])
}

// GET /api/activity
func handleActivity(w http.ResponseWriter, r *http.Request) {
	res, err := dbClient.Activity()
	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, res)
}

// GET /api/schemas
func handleGetSchemas(w http.ResponseWriter, r *http.Request) {
	names, err := dbClient.Schemas()

	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, names)
}

// GET /api/tables
func handleGetTables(w http.ResponseWriter, r *http.Request) {
	names, err := dbClient.Tables()

	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, names)
}

func handleGetTable(w http.ResponseWriter, r *http.Request, table string) {
	res, err := dbClient.Table(table)
	LogInfof("table: '%s'\n", table)

	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, res)
}

func apiGetTableRows(w http.ResponseWriter, r *http.Request, table string) {
	LogInfof("table='%s'\n", table)
	limit := 1000 // Number of rows to fetch
	limitVal := r.FormValue("limit")

	if limitVal != "" {
		num, err := strconv.Atoi(limitVal)

		if err != nil {
			serveJSON(w, r, 400, Error{"Invalid limit value"})
			return
		}

		if num <= 0 {
			serveJSON(w, r, 400, Error{"Limit should be greater than 0"})
			return
		}

		limit = num
	}

	opts := RowsOptions{
		Limit:      limit,
		SortColumn: r.FormValue("sort_column"),
		SortOrder:  r.FormValue("sort_order"),
	}

	res, err := dbClient.TableRows(table, opts)

	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, res)
}

func apiGetTableInfo(w http.ResponseWriter, r *http.Request, table string) {
	res, err := dbClient.TableInfo(table)

	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, res.Format()[0])
}

func apiGetTableIndexes(w http.ResponseWriter, r *http.Request, table string) {
	LogInfof("table='%s'\n", table)
	res, err := dbClient.TableIndexes(table)

	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, res)
}

// GET /api/tables/:table/:action
func handleTablesDispatch(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	uriPath := uri[len("/api/tables/"):]
	parts := strings.SplitN(uriPath, "/", 2)
	table := parts[0]
	LogInfof("table='%s'\n", table)
	if len(parts) == 1 {
		handleGetTable(w, r, table)
		return
	}
	cmd := parts[1]
	if cmd == "rows" {
		apiGetTableRows(w, r, table)
		return
	}
	if cmd == "info" {
		apiGetTableInfo(w, r, table)
		return
	}
	if cmd == "indexes" {
		apiGetTableIndexes(w, r, table)
		return
	}
	LogErrorf("unknown cmd: '%s'\n", cmd)
	http.NotFound(w, r)
}

func registerHTTPHandlers() {
	http.HandleFunc("/", get(handleIndex))
	http.HandleFunc("/s/", get(handleStatic))
	http.HandleFunc("/api/connect", timeit(post(handleConnect)))
	http.HandleFunc("/api/history", get(handleHistory))
	http.HandleFunc("/api/bookmarks", get(handleBookmarks))

	http.HandleFunc("/api/databases", timeit(hasConnection(handleGetDatabases)))
	http.HandleFunc("/api/connection", timeit(hasConnection(handleConnectionInfo)))
	http.HandleFunc("/api/activity", timeit(hasConnection(handleActivity)))
	http.HandleFunc("/api/schemas", timeit(hasConnection(handleGetSchemas)))
	http.HandleFunc("/api/tables", timeit(hasConnection(handleGetTables)))
	http.HandleFunc("/api/tables/", timeit(hasConnection(handleTablesDispatch)))
	http.HandleFunc("/api/query", timeit(hasConnection(handleRunQuery)))
	http.HandleFunc("/api/explain", timeit(hasConnection(handleExplainQuery)))
}

func startWebServer() {
	registerHTTPHandlers()
	httpAddr := fmt.Sprintf(":%v", options.HttpPort)
	fmt.Printf("Started running on %s\n", httpAddr)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		log.Fatalf("http.ListendAndServer() failed with %s\n", err)
	}
	fmt.Printf("Exited\n")
}
