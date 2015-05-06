package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
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

// CookieValue contains data we set in browser cookie
type CookieValue struct {
	UserID            int
	IsLoggedIn        bool
	GoogleAnalyticsID string
}

type ReqOpts uint

const (
	OnlyGet ReqOpts = 1 << iota
	OnlyPost
)

// ReqContext contains data that is useful to access in every http handler
type ReqContext struct {
	Cookie    *CookieValue
	DbUser    *DbUser
	IsAdmin   bool
	TimeStart time.Time
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

// Error is error message sent to client as json response for backend errors
type Error struct {
	Message string `json:"error"`
}

// NewError creates a new Error
func NewError(err error) Error {
	return Error{err.Error()}
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

func asset(fileName string) ([]byte, error) {
	//fmt.Fprintf(os.Stderr, "asset: %s\n", fileName)
	return ioutil.ReadFile(fileName)
}

// HandlerWithCtxFunc is like http.HandlerFunc but with additional ReqContext argument
type HandlerWithCtxFunc func(*ReqContext, http.ResponseWriter, *http.Request)

// TODO: wrap w within CountingResponseWriter
func withctx(f HandlerWithCtxFunc, opts ReqOpts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := strings.ToUpper(r.Method)
		if opts&OnlyGet != 0 {
			if method != "GET" {
				http.NotFound(w, r)
				return
			}
		}
		if opts&OnlyPost != 0 {
			if method != "POST" {
				http.NotFound(w, r)
				return
			}
		}

		ctx := &ReqContext{
			Cookie:    getOrCreateCookie(w, r),
			TimeStart: time.Now(),
		}

		if ctx.Cookie.UserID != -1 {
			ctx.DbUser, _ = dbGetUserByIDCached(ctx.Cookie.UserID)
			if ctx.DbUser == nil {
				// if we have valid UserID, we should be able to look up the user
				LogErrorf("dbGetUserByIDCached() returned nil for userId %d, url: %s\n", ctx.Cookie.UserID, r.RequestURI)
				http.NotFound(w, r)
				return
			}
			ctx.IsAdmin = isAdminUser(ctx.DbUser)
		}

		f(ctx, w, r)
		// TODO: log this to a file for further analysis
		LogInfof("%s took %s\n", r.RequestURI, time.Since(ctx.TimeStart))

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
		serveString(w, r, 404, err.Error())
		return
	}

	if len(data) == 0 {
		serveString(w, r, 404, "Asset is empty")
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

func writeHeader(w http.ResponseWriter, code int, contentType string) {
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	w.WriteHeader(code)
}

func serveJSON(w http.ResponseWriter, r *http.Request, code int, v interface{}) error {
	writeHeader(w, code, "application/json")
	encoder := json.NewEncoder(w)
	return encoder.Encode(v)
}

func serveJSONP(w http.ResponseWriter, r *http.Request, code int, v interface{}, jsonp string) error {
	writeHeader(w, code, "application/json")

	if jsonp == "" {
		return serveJSON(w, r, code, v)
	}

	b, err := json.Marshal(v)
	if err != nil {
		// should never happen
		LogErrorf("json.MarshalIndent() failed with %q\n", err)
		return err
	}
	res := []byte(jsonp)
	res = append(res, '(')
	res = append(res, b...)
	res = append(res, ')')
	_, err = w.Write(res)
	return err
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
func handleConnect(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
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
func handleGetDatabases(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	names, err := dbClient.Databases()

	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, names)
}

func handleQuery(ctx *ReqContext, w http.ResponseWriter, r *http.Request, query string) {
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
func handleRunQuery(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.FormValue("query"))
	LogInfof("query: '%s'\n", query)

	if query == "" {
		serveJSON(w, r, 400, errors.New("Query parameter is missing"))
		return
	}

	handleQuery(ctx, w, r, query)
}

// GET | POST /api/explain
func handleExplainQuery(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.FormValue("query"))
	LogInfof("query: '%s'\n", query)

	if query == "" {
		serveJSON(w, r, 400, errors.New("Query parameter is missing"))
		return
	}

	handleQuery(ctx, w, r, fmt.Sprintf("EXPLAIN ANALYZE %s", query))
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
func handleConnectionInfo(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	res, err := dbClient.Info()

	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, res.Format()[0])
}

// GET /api/activity
func handleActivity(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	res, err := dbClient.Activity()
	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, res)
}

// GET /api/schemas
func handleGetSchemas(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	names, err := dbClient.Schemas()

	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, names)
}

// GET /api/tables
func handleGetTables(ctx *ReqContext, w http.ResponseWriter, r *http.Request) {
	names, err := dbClient.Tables()

	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, names)
}

func handleGetTable(ctx *ReqContext, w http.ResponseWriter, r *http.Request, table string) {
	res, err := dbClient.Table(table)
	LogInfof("table: '%s'\n", table)

	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, res)
}

func apiGetTableRows(ctx *ReqContext, w http.ResponseWriter, r *http.Request, table string) {
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

func apiGetTableInfo(ctx *ReqContext, w http.ResponseWriter, r *http.Request, table string) {
	res, err := dbClient.TableInfo(table)

	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, res.Format()[0])
}

func apiGetTableIndexes(ctx *ReqContext, w http.ResponseWriter, r *http.Request, table string) {
	LogInfof("table='%s'\n", table)
	res, err := dbClient.TableIndexes(table)

	if err != nil {
		serveJSON(w, r, 400, NewError(err))
		return
	}

	serveJSON(w, r, 200, res)
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
	LogInfof("jsonp: '%s'\n", jsonp)
	if ctx.DbUser != nil {
		LogInfof("dbUser: %#v\n", ctx.DbUser)
	}

	v := struct {
		Email      string
		IsLoggedIn bool
	}{
		IsLoggedIn: ctx.Cookie.IsLoggedIn,
	}
	if ctx.DbUser != nil {
		v.Email = ctx.DbUser.Email
	}
	LogInfof("v: %#v\n", v)
	serveJSONP(w, r, 200, v, jsonp)
}

func registerHTTPHandlers() {
	http.HandleFunc("/", withctx(handleIndex, OnlyGet))
	http.HandleFunc("/s/", withctx(handleStatic, OnlyGet))
	http.HandleFunc("/api/connect", withctx(handleConnect, OnlyPost))
	http.HandleFunc("/api/history", handleHistory)
	http.HandleFunc("/api/bookmarks", handleBookmarks)

	http.HandleFunc("/api/databases", withctx(handleGetDatabases, 0))
	http.HandleFunc("/api/connection", withctx(handleConnectionInfo, 0))
	http.HandleFunc("/api/activity", withctx(handleActivity, 0))
	http.HandleFunc("/api/schemas", withctx(handleGetSchemas, 0))
	http.HandleFunc("/api/tables", withctx(handleGetTables, 0))
	http.HandleFunc("/api/tables/", withctx(handleTablesDispatch, 0))
	http.HandleFunc("/api/query", withctx(handleRunQuery, 0))
	http.HandleFunc("/api/explain", withctx(handleExplainQuery, 0))
	http.HandleFunc("/api/userinfo", withctx(handleUserInfo, 0))

	http.HandleFunc("/logingoogle", handleLoginGoogle)
	http.HandleFunc("/logout", withctx(handleLogout, 0))
	http.HandleFunc("/googleoauth2cb", withctx(handleOauthGoogleCallback, 0))
	http.HandleFunc("/showmyhost", handleShowMyHost)
}

// GET /showmyhost, for testing only
func handleShowMyHost(w http.ResponseWriter, r *http.Request) {
	s := getMyHost(r)
	serveString(w, r, 200, "me: %s\n", s)
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
