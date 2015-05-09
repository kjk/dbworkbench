package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
)

const (
	cookieAuthKeyHex = "4d1e71694154384ee43b20d06c710a2e6149126efadd140c66a4fa9bed9cb0bd"
	cookieEncrKeyHex = "514171b116075c359db13549dc4b15924a5c0d0615c8b9d81436ba4d712132a1"
	cookieName       = "dbwckie" // "database workbench cookie"
)

var (
	cookieAuthKey []byte
	cookieEncrKey []byte

	secureCookie *securecookie.SecureCookie
)

// HandlerWithCtxFunc is like http.HandlerFunc but with additional ReqContext argument
type HandlerWithCtxFunc func(*ReqContext, http.ResponseWriter, *http.Request)

// CookieValue contains data we set in browser cookie
type CookieValue struct {
	UserID            int
	IsLoggedIn        bool
	GoogleAnalyticsID string
}

// ReqOpts is a set of flags passed to withCtx
type ReqOpts uint

const (
	// OnlyGet tells to reject non-GET requests
	OnlyGet ReqOpts = 1 << iota
	// OnlyPost tells to reject non-POST requests
	OnlyPost
	// MustBeLoggedIn tells to reject requests if user is not logged in
	MustBeLoggedIn
	// MustHaveConnection tells to reject requests if conn_id request arg is not
	// present. Implies MustBeLoggedIn
	MustHaveConnection
	// IsJSON denotes a handler that is serving JSON requests and should send
	// errors as { "error": "error message" }
	IsJSON
)

// ReqContext contains data that is useful to access in every http handler
type ReqContext struct {
	Cookie       *CookieValue
	User         *User
	IsAdmin      bool
	TimeStart    time.Time
	ConnectionID int
	Client       *Client
}

func isAdminUser(u *DbUser) bool {
	if u != nil {
		switch u.Email {
		case "kkowalczyk@gmail.com":
			return true
		}
	}
	return false
}

func initCookieMust() {
	var err error
	cookieAuthKey, err = hex.DecodeString(cookieAuthKeyHex)
	fatalIfErr(err, "hex.DecodeString(cookieAuthKeyHex)")
	cookieEncrKey, err = hex.DecodeString(cookieEncrKeyHex)
	fatalIfErr(err, "hex.DecodeString(cookieEncrKeyHex)")
	secureCookie = securecookie.New(cookieAuthKey, cookieEncrKey)
	// verify auth/encr keys are correct
	val := map[string]string{
		"foo": "bar",
	}
	_, err = secureCookie.Encode(cookieName, val)
	fatalIfErr(err, "secureCookie.Encode")
}

func setCookie(w http.ResponseWriter, cookieVal *CookieValue) {
	if encoded, err := secureCookie.Encode(cookieName, cookieVal); err == nil {
		// TODO: set expiration (Expires    time.Time) long time in the future?
		cookie := &http.Cookie{
			Name:     cookieName,
			Value:    encoded,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	} else {
		LogErrorf("secureCookie.Encode() failed with '%s'\n", err)
	}
}

func genAndSetNewCookieValue(w http.ResponseWriter) *CookieValue {
	c := &CookieValue{
		GoogleAnalyticsID: generateUUID(),
	}
	setCookie(w, c)
	return c
}

func getOrCreateCookie(w http.ResponseWriter, r *http.Request) *CookieValue {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return genAndSetNewCookieValue(w)
	}
	var cv CookieValue
	if err = secureCookie.Decode(cookieName, cookie.Value, &cv); err != nil {
		// most likely expired cookie, so ignore and delete
		LogErrorf("secureCookie.Decode() failed with %s\n", err)
		return genAndSetNewCookieValue(w)
	}
	//LogVerbosef("Got cookie %#v\n", ret)
	if cv.GoogleAnalyticsID == "" {
		LogError("cv.GoogleAnalyticsID is empty string\n")
		cv.GoogleAnalyticsID = generateUUID()
		setCookie(w, &cv)
	}
	return &cv
}

func asset(fileName string) ([]byte, error) {
	//fmt.Fprintf(os.Stderr, "asset: %s\n", fileName)
	return ioutil.ReadFile(fileName)
}

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
			Cookie:    getOrCreateCookie(w, r),
			TimeStart: time.Now()}

		if opts&MustBeLoggedIn != 0 || opts&MustHaveConnection != 0 {
			if !ctx.Cookie.IsLoggedIn || ctx.Cookie.UserID == -1 {
				serveError(w, r, isJSON, "must be logged")
				return
			}
		}

		if ctx.Cookie.UserID != -1 {
			ctx.User, _ = dbGetUserCopyByIDCached(ctx.Cookie.UserID)
			if ctx.User == nil {
				// if we have valid UserID, we should be able to look up the user
				serveError(w, r, isJSON, fmt.Sprintf("dbGetUserByIDCached() returned nil for userId %d", ctx.Cookie.UserID))
				return
			}
			ctx.IsAdmin = isAdminUser(ctx.User.DbUser)
		}
		user := ctx.User

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
			if user.ConnInfo == nil || user.ConnInfo.ConnectionID != connID {
				var errMsg string
				if user.ConnInfo == nil {
					errMsg = "ctx.User.ConnInfo == nil"
				} else {
					errMsg = fmt.Sprintf("ctx.User.ConnInfo.ConnectionID != connID (%d != %d)", user.ConnInfo.ConnectionID, connID)
				}
				serveError(w, r, isJSON, errMsg)
				return
			}
			ctx.ConnectionID = connID
			ctx.Client = user.ConnInfo.Client
		}

		f(ctx, cw, r)
		if !strings.HasPrefix(r.RequestURI, "/s/") {
			// TODO: log this to a file for further analysis
			LogInfof("%s took %s, code: %d\n", r.RequestURI, time.Since(ctx.TimeStart), cw.Code)
		}

		go func(r *http.Request, gaID string) {
			err := gaLogPageView(r.UserAgent(), gaID, getClientIP(r), r.URL.Path, "", nil)

			if err != nil {
				log.Printf("Unable to log GA PageView: %v\n", err)
			}
		}(r, ctx.Cookie.GoogleAnalyticsID)
	}
}

func serveStatic(w http.ResponseWriter, r *http.Request, path string) {
	data, err := asset(path)

	if err != nil {
		LogErrorf("asset('%s') failed with '%s'\n", path, err)
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
	if uri != "/" {
		http.NotFound(w, r)
		return
	}
	serveStatic(w, r, "s/index.html")
}

// GET /s/:path
func handleStatic(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	path := "s/" + r.URL.Path[len("/s/"):]
	//LogInfof("path='%s'\n", path)
	serveStatic(w, r, path)
}

// POST /api/connect
func handleConnect(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		serveJSONError(w, r, "Url parameter is required")
		return
	}

	opts := Options{URL: url}
	url, err := formatConnectionUrl(opts)

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

	userID := ctx.User.DbUser.ID
	removeCurrentUserConnectionInfo(userID)
	connInfo := addConnectionInfo(url, userID, client)

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
	userID := ctx.User.DbUser.ID
	removeCurrentUserConnectionInfo(userID)
	v := struct {
		Message string
	}{
		Message: "ok",
	}
	serveJSON(w, r, v)
}

// GET /api/databases
func handleGetDatabases(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	names, err := ctx.User.ConnInfo.Client.Databases()
	if err != nil {
		serveJSONError(w, r, err)
		return
	}
	serveJSON(w, r, names)
}

func handleQuery(ctx *ReqContext, w http.ResponseWriter, r *http.Request, query string) {
	result, err := ctx.User.ConnInfo.Client.Query(query)

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
	serveJSON(w, r, ctx.User.ConnInfo.Client.history)
}

// GET /api/bookmarks
func handleBookmarks(w http.ResponseWriter, r *http.Request) {
	bookmarks, err := readAllBookmarks(bookmarksPath())
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	serveJSON(w, r, bookmarks)
}

// GET /api/connection
func handleConnectionInfo(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	res, err := ctx.User.ConnInfo.Client.Info()
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	serveJSON(w, r, res.Format()[0])
}

// GET /api/activity
func handleActivity(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	res, err := ctx.User.ConnInfo.Client.Activity()
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	serveJSON(w, r, res)
}

// GET /api/schemas
func handleGetSchemas(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	names, err := ctx.User.ConnInfo.Client.Schemas()
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	serveJSON(w, r, names)
}

// GET /api/tables
func handleGetTables(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	//LogInfof("connID: %d\n", ctx.ConnectionID)
	names, err := ctx.User.ConnInfo.Client.Tables()
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	serveJSON(w, r, names)
}

func handleGetTable(ctx *ReqContext, w http.ResponseWriter, r *http.Request, table string) {
	res, err := ctx.User.ConnInfo.Client.Table(table)
	//LogInfof("table: '%s'\n", table)
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	serveJSON(w, r, res)
}

func apiGetTableRows(ctx *ReqContext, w http.ResponseWriter, r *http.Request, table string) {
	LogInfof("table='%s'\n", table)
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

	res, err := ctx.User.ConnInfo.Client.TableRows(table, opts)
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	serveJSON(w, r, res)
}

func apiGetTableInfo(ctx *ReqContext, w http.ResponseWriter, r *http.Request, table string) {
	res, err := ctx.User.ConnInfo.Client.TableInfo(table)
	if err != nil {
		serveJSONError(w, r, err)
		return
	}

	serveJSON(w, r, res.Format()[0])
}

func apiGetTableIndexes(ctx *ReqContext, w http.ResponseWriter, r *http.Request, table string) {
	LogInfof("table='%s'\n", table)
	res, err := ctx.User.ConnInfo.Client.TableIndexes(table)
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
	LogInfof("User: %#v\n", ctx.User)

	v := struct {
		Email        string
		IsLoggedIn   bool
		ConnectionID int
	}{
		IsLoggedIn: ctx.Cookie.IsLoggedIn,
	}
	if ctx.User != nil {
		v.Email = ctx.User.DbUser.Email
		if ctx.User.ConnInfo != nil {
			v.ConnectionID = ctx.User.ConnInfo.ConnectionID
		}
	}
	LogInfof("v: %#v\n", v)
	serveJSONP(w, r, v, jsonp)
}

func registerHTTPHandlers() {
	http.HandleFunc("/", withCtx(handleIndex, OnlyGet))
	http.HandleFunc("/s/", withCtx(handleStatic, OnlyGet))

	http.HandleFunc("/api/activity", withCtx(handleActivity, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/bookmarks", handleBookmarks)
	http.HandleFunc("/api/connect", withCtx(handleConnect, OnlyPost|MustBeLoggedIn|IsJSON))
	http.HandleFunc("/api/connection", withCtx(handleConnectionInfo, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/databases", withCtx(handleGetDatabases, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/disconnect", withCtx(handleDisconnect, OnlyPost|MustHaveConnection|IsJSON))
	http.HandleFunc("/api/explain", withCtx(handleExplainQuery, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/history", withCtx(handleHistory, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/schemas", withCtx(handleGetSchemas, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/tables", withCtx(handleGetTables, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/tables/", withCtx(handleTablesDispatch, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/query", withCtx(handleRunQuery, MustHaveConnection|IsJSON))
	http.HandleFunc("/api/userinfo", withCtx(handleUserInfo, MustBeLoggedIn|IsJSON))

	http.HandleFunc("/logingoogle", handleLoginGoogle)
	http.HandleFunc("/logout", withCtx(handleLogout, 0))
	http.HandleFunc("/googleoauth2cb", withCtx(handleOauthGoogleCallback, 0))
	http.HandleFunc("/showmyhost", handleShowMyHost)
}

// GET /showmyhost, for testing only
func handleShowMyHost(w http.ResponseWriter, r *http.Request) {
	s := getMyHost(r)
	servePlainText(w, r, 200, "me: %s\n", s)
}

func startWebServer() {
	registerHTTPHandlers()
	httpAddr := fmt.Sprintf(":%v", options.HTTPPort)
	fmt.Printf("Started running on %s\n", httpAddr)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		log.Fatalf("http.ListendAndServer() failed with %s\n", err)
	}
	fmt.Printf("Exited\n")
}
