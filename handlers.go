package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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
)

// ReqContext contains data that is useful to access in every http handler
type ReqContext struct {
	TimeStart    time.Time
	ConnectionID int
	ConnInfo     *ConnectionInfo
	Client       *Client
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
			// TODO: log this to a file for further analysis
			LogInfof("%s took %s, code: %d\n", r.RequestURI, time.Since(ctx.TimeStart), cw.Code)
		}
	}
}

func serveStatic(w http.ResponseWriter, r *http.Request, path string) {
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

func loadResourcesFromZipReader(zr *zip.Reader) error {
	for _, f := range zr.File {
		name := f.Name
		rc, err := f.Open()
		if err != nil {
			return err
		}
		d, err := ioutil.ReadAll(rc)
		rc.Close()
		if err != nil {
			return err
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

func serveResourceFromZip(w http.ResponseWriter, r *http.Request, path string) {
	//LogInfof("serving '%s' from zip\n", path)
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

// GET /s/:path
func handleStatic(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	resourcePath := "s/" + r.URL.Path[len("/s/"):]
	//LogInfof("path='%s'\n", path)
	if options.ResourcesFromZip {
		serveResourceFromZip(w, r, resourcePath)
	} else {
		serveStatic(w, r, resourcePath)
	}
}

// POST /api/connect
func handleConnect(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		serveJSONError(w, r, "Url parameter is required")
		return
	}

	opts := Options{URL: url}
	url, err := formatConnectionURL(opts)

	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	client, err := NewClientFromUrl(url)
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	err = client.Test()
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	info, err := client.Info()
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	connInfo := addConnectionInfo(url, client)

	i := info.Format()[0]
	currDb, ok := i["current_database"]
	if !ok {
		serveJSONError(w, r, "no current_database")
		return
	}
	currDbStr, ok := currDb.(string)
	if !ok {
		serveJSONError(w, r, "invalid type")
		return
	}
	v := struct {
		ConnectionID    int
		CurrentDatabase string
	}{
		ConnectionID:    connInfo.ConnectionID,
		CurrentDatabase: currDbStr,
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

func handleQuery(ctx *ReqContext, w http.ResponseWriter, r *http.Request, query string) {
	updateConnectionLastAccess(ctx.ConnInfo.ConnectionID)
	result, err := ctx.ConnInfo.Client.Query(query)

	if err != nil {
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

// GET | POST /api/query
func handleRunQuery(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.FormValue("query"))
	//LogInfof("query: '%s'\n", query)

	if query == "" {
		serveJSONError(w, r, "Query parameter is missing")
		return
	}

	handleQuery(ctx, w, r, query)
}

// GET | POST /api/explain
func handleExplainQuery(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.FormValue("query"))
	//LogInfof("query: '%s'\n", query)

	if query == "" {
		serveJSONError(w, r, "Query parameter is missing")
		return
	}

	handleQuery(ctx, w, r, fmt.Sprintf("EXPLAIN ANALYZE %s", query))
}

// GET /api/history
func handleHistory(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	serveJSON(w, r, ctx.ConnInfo.Client.history)
}

// GET /api/bookmarks
func handleBookmarks(w http.ResponseWriter, r *http.Request) {
	bookmarks, err := readAllBookmarks()
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	serveJSON(w, r, bookmarks)
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

// GET /api/schemas
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

func apiGetTableRows(ctx *ReqContext, w http.ResponseWriter, r *http.Request, table string) {
	LogInfof("table='%s'\n", table)
	updateConnectionLastAccess(ctx.ConnInfo.ConnectionID)
	limit := 1000 // Number of rows to fetch
	limitVal := r.FormValue("limit")

	if limitVal != "" {
		num, err := strconv.Atoi(limitVal)
		if err != nil {
			serveJSONError(w, r, "Invalid limit value")
			return
		}

		if num <= 0 {
			serveJSONError(w, r, "Limit should be greater than 0")
			return
		}

		limit = num
	}

	opts := RowsOptions{
		Limit:      limit,
		SortColumn: r.FormValue("sort_column"),
		SortOrder:  r.FormValue("sort_order"),
	}

	res, err := ctx.ConnInfo.Client.TableRows(table, opts)
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

// GET /api/tables/:table/:action
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
	case "rows":
		apiGetTableRows(ctx, w, r, table)
	case "info":
		apiGetTableInfo(ctx, w, r, table)
	case "indexes":
		apiGetTableIndexes(ctx, w, r, table)
	default:
		LogErrorf("unknown cmd: '%s'\n", cmd)
		http.NotFound(w, r)
	}
}

// GET /api/userinfo
// Returns information about the user
// Arguments:
//  - jsonp : jsonp wrapper, optional
func handleUserInfo(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	jsonp := strings.TrimSpace(r.FormValue("jsonp"))
	//LogInfof("User: %#v\n", ctx.User)

	v := struct {
		Email        string // TODO: remove
		ConnectionID int
	}{
		Email: "foo@bar.com",
	}
	connID := getFirstConnectionID()
	if -1 != connID {
		v.ConnectionID = connID
	}
	LogInfof("v: %#v\n", v)
	serveJSONP(w, r, v, jsonp)
}

func registerHTTPHandlers() {
	http.HandleFunc("/", withCtx(handleIndex, OnlyGet))
	http.HandleFunc("/s/", withCtx(handleStatic, OnlyGet))

	http.HandleFunc("/api/activity", withCtx(handleActivity, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/bookmarks", handleBookmarks)
	http.HandleFunc("/api/connect", withCtx(handleConnect, OnlyPost|IsJSON))
	http.HandleFunc("/api/connection", withCtx(handleConnectionInfo, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/databases", withCtx(handleGetDatabases, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/disconnect", withCtx(handleDisconnect, OnlyPost|MustHaveConnection|IsJSON))
	http.HandleFunc("/api/explain", withCtx(handleExplainQuery, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/history", withCtx(handleHistory, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/schemas", withCtx(handleGetSchemas, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/tables", withCtx(handleGetTables, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/tables/", withCtx(handleTablesDispatch, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/query", withCtx(handleRunQuery, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/userinfo", withCtx(handleUserInfo, IsJSON))

	http.HandleFunc("/showmyhost", handleShowMyHost)
}

// GET /showmyhost, for testing only
func handleShowMyHost(w http.ResponseWriter, r *http.Request) {
	s := getMyHost(r)
	servePlainText(w, r, 200, "me: %s\n", s)
}

func startWebServer() {
	registerHTTPHandlers()
	httpAddr := fmt.Sprintf("127.0.0.1:%v", options.HTTPPort)
	fmt.Printf("Started running on %s\n", httpAddr)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		log.Fatalf("http.ListendAndServer() failed with %s\n", err)
	}
	fmt.Printf("Exited\n")
}
