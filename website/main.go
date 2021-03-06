package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

const (
	// for auto-update
	latestMacVersion = "0.2.3"
	latestWinVersion = "0.2.3"

	indexMac = "/gui-database-client-for-mysql-postgresql-mac-osx"
	indexWin = "/gui-database-client-for-mysql-postgresql-windows"
)

var (
	flgUsageTest bool

	httpAddr = ":5555"

	dataDir string

	urlToFileMap = map[string]string{
		indexMac: "for-mac.html",
		indexWin: "for-windows.html",
	}

	redirects = map[string]string{
		"/for-mac":     indexMac,
		"/for-windows": indexWin,
		// those were main urls before mysql support
		"/gui-database-client-for-postgresql-mac-osx": indexMac,
		"/gui-database-client-for-postgresql-windows": indexWin,
	}
)

func latestMacDownloadURL() string {
	return fmt.Sprintf("https://kjkpub.s3.amazonaws.com/software/dbhero/rel/dbHero-%s.zip", latestMacVersion)
}

func latestWinDownloadURL() string {
	return fmt.Sprintf("https://kjkpub.s3.amazonaws.com/software/dbhero/rel/dbHero-setup-%s.exe", latestWinVersion)
}

// LogVerbosef logs verbose info
func LogVerbosef(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	fmt.Print(s)
}

// LogInfof logs additional info
func LogInfof(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	fmt.Print(s)
}

// LogErrorf logs errors
func LogErrorf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	fmt.Print(s)
}

// LogFatalf logs an error and terminates the app
func LogFatalf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	fmt.Print(s)
	log.Fatal(s)
}

var extraMimeTypes = map[string]string{
	".icon": "image-x-icon",
	".ttf":  "application/x-font-ttf",
	".woff": "application/x-font-woff",
	".eot":  "application/vnd.ms-fontobject",
	".svg":  "image/svg+xml",
}

// MimeTypeByExtensionExt is like mime.TypeByExtension but supports more types
// and defaults to text/plain
func MimeTypeByExtensionExt(name string) string {
	ext := strings.ToLower(filepath.Ext(name))
	result := mime.TypeByExtension(ext)

	if result == "" {
		result = extraMimeTypes[ext]
	}

	if result == "" {
		result = "text/plain; charset=utf-8"
	}

	return result
}

// data dir is ../../data on the server or ~/data/dbhero-website locally
// the important part is that it's outside of directory with the code
func getDataDir() string {
	if dataDir != "" {
		return dataDir
	}
	// on the server, must be done first because ExpandTildeInPath()
	// doesn't work when cross-compiled on mac for linux
	dataDir = filepath.Join("..", "..", "data")
	if PathExists(dataDir) {
		return dataDir
	}
	dataDir = ExpandTildeInPath("~/data/dbhero-website")
	if PathExists(dataDir) {
		return dataDir
	}
	log.Fatal("data directory (../../data or ~/data/dbhero-website) doesn't exist")
	return ""
}

func writeHeader(w http.ResponseWriter, code int, contentType string) {
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	w.WriteHeader(code)
}

func servePlainText(w http.ResponseWriter, r *http.Request, code int, format string, args ...interface{}) {
	writeHeader(w, code, "text/plain")
	var err error
	if len(args) > 0 {
		_, err = w.Write([]byte(fmt.Sprintf(format, args...)))
	} else {
		_, err = w.Write([]byte(format))
	}
	if err != nil {
		LogErrorf("err: '%s'\n", err)
	}
}

func serveData(w http.ResponseWriter, r *http.Request, code int, contentType string, data []byte) {
	if len(contentType) > 0 {
		w.Header().Set("Content-Type", contentType)
	}
	w.WriteHeader(code)
	w.Write(data)
}

// returns nil if not a POST request or error reading data
func readRawPostData(r *http.Request) []byte {
	if r.Method != "POST" {
		LogErrorf("readRawPostData: r.Method is %s, not POST\n", r.Method)
		return nil
	}
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		LogErrorf("ioutil.ReadAll() failed with '%s'\n", err)
		return nil
	}
	//LogInfof("readRawPostData: len(d) = %d\n'%s'\n", len(d), string(d))
	if len(d) == 0 {
		return nil
	}
	return d
}

func prependInfoFromRequest(d []byte, ip, ver, os string) []byte {
	now := time.Now()
	t := now.Unix()
	day := now.Format("2006-01-02")
	s := fmt.Sprintf("time: %d\nday: %s\nip: %s\nver_in_url: %s\nos_in_url: %s\n", t, day, ip, ver, os)
	d1 := []byte(s)
	return append(d1, d...)
}

// url: /api/winupdatecheck?ver=${ver}
func handleWinUpdateCheck(w http.ResponseWriter, r *http.Request) {
	//LogInfof("handleWinUpdateCheck\n")

	d := readRawPostData(r)
	ver := r.FormValue("ver")
	ip := getIPFromRequest(r)

	s := fmt.Sprintf(`ver: %s
url: %s`, latestWinVersion, latestWinDownloadURL())
	servePlainText(w, r, 200, s)

	d = prependInfoFromRequest(d, ip, ver, "windows")
	//LogInfof("handleWinUpdateCheck: '%s'\n", string(d))
	recordUsage(d)
}

// url: /api/macupdatecheck?ver=${ver}
func handleMacUpdateCheck(w http.ResponseWriter, r *http.Request) {
	//LogInfof("handleMacUpdateCheck\n")

	d := readRawPostData(r)
	ver := r.FormValue("ver")
	ip := getIPFromRequest(r)

	s := fmt.Sprintf(`ver: %s
url: %s`, latestMacVersion, latestMacDownloadURL())
	servePlainText(w, r, 200, s)

	d = prependInfoFromRequest(d, ip, ver, "windows")
	recordUsage(d)
}

// heuristic to determine if request is coming from Windows
func isWindowsUserAgent(ua string) bool {
	return strings.Contains(ua, "Windows")
}

// heuristic to determine if request is coming from Mac
func isMacUserAgent(ua string) bool {
	return strings.Contains(ua, "Macintosh")
}

func redirectIndex(w http.ResponseWriter, r *http.Request) {
	ua := r.UserAgent()
	if isMacUserAgent(ua) {
		http.Redirect(w, r, "gui-database-client-for-postgresql-mac-osx", http.StatusFound /* 302 */)
		return
	}

	// for windows and everything else
	http.Redirect(w, r, "gui-database-client-for-postgresql-windows", http.StatusFound /* 302 */)
}

// url: /
func handleIndex(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	LogVerbosef("handleIndex: '%s'\n", uri)
	if uri == "/" {
		redirectIndex(w, r)
		return
	}

	if redirect := redirects[uri]; redirect != "" {
		http.Redirect(w, r, redirect, http.StatusFound /* 302 */)
		return
	}

	var name string
	if name = urlToFileMap[uri]; name == "" {
		// map /foo and /foo.html to /s/foo.html such file exists
		name = uri[1:]
	}

	path := filepath.Join("www", name)
	d, err := ioutil.ReadFile(path)
	if err != nil {
		path += ".html"
		d, err = ioutil.ReadFile(path)
	}
	if err != nil {
		LogErrorf("ioutil.ReadFile('%s') failed with '%s'\n", path, err)
		http.NotFound(w, r)
		return
	}
	serveData(w, r, 200, MimeTypeByExtensionExt(path), d)
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

func serveTemplate(w http.ResponseWriter, r *http.Request, path string, v interface{}) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		LogErrorf("tempalte.ParseFiles('%s') failed with '%s'\n", path, err)
		servePlainText(w, r, 500, fmt.Sprintf("tempalte.ParseFiles('%s') failed with '%s'\n", path, err))
		return
	}
	name := filepath.Base(path)
	err = tmpl.ExecuteTemplate(w, name, v)
	if err != nil {
		LogErrorf("tmpl.ExecuteTemplate('%s') failed with '%s'\n", name, err)
	}
}

// url: /s/:path
func handleStatic(w http.ResponseWriter, r *http.Request) {
	LogInfof("handleStatic: '%s'\n", r.URL.Path)
	name := r.URL.Path[len("/s/"):]
	resourcePath := filepath.Join("www", name)
	serveStatic(w, r, resourcePath)
}

// url: /admin/usage?s=acergytvw
// s is the simplest for of authentication - via unguessable url
func handleUsage(w http.ResponseWriter, r *http.Request) {
	LogInfof("handleUsage: '%s'\n", r.URL.Path)
	secret := r.FormValue("s")
	if secret != "acergytvw" {
		servePlainText(w, r, 500, "can't see me")
		return
	}
	//path := "usage-for-testing.txt"
	path := usageFilePath()
	res, err := parseUsage(path)
	if err != nil {
		LogErrorf("parseUsage() failed with '%s'\n", err)
		servePlainText(w, r, 500, fmt.Sprintf("error: '%s'", err))
		return
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		LogErrorf("json.Marshal() failed with '%s'\n", err)
		servePlainText(w, r, 500, fmt.Sprintf("error: '%s'", err))
		return
	}
	path = filepath.Join("www", "usage_stats.html")
	v := struct {
		UsersJSON string
	}{
		UsersJSON: string(resJSON),
	}
	serveTemplate(w, r, path, v)
}

// https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/
func makeHTTPServer() *http.Server {
	mux := &http.ServeMux{}

	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/s/", handleStatic)
	mux.HandleFunc("/admin/usage", handleUsage)
	mux.HandleFunc("/api/winupdatecheck", handleWinUpdateCheck)
	mux.HandleFunc("/api/macupdatecheck", handleMacUpdateCheck)

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		// TODO: 1.8 only
		// IdleTimeout:  120 * time.Second,
		Handler: mux,
	}
	// TODO: track connections and their state
	return srv
}

func hostPolicy(ctx context.Context, host string) error {
	if strings.HasSuffix(host, "dbheroapp.com") {
		return nil
	}
	return errors.New("acme/autocert: only *.dbheroapp.com hosts are allowed")
}

func parseCmdLine() {
	flag.BoolVar(&flgUsageTest, "usage-test", false, "test parsing of usage.txt")
	flag.Parse()
}

func main() {
	parseCmdLine()
	if flgUsageTest {
		testUsageParse()
		os.Exit(0)
	}
	openUsageFileMust()

	if IsLinux() {
		srv := makeHTTPServer()
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
		}
		srv.Addr = ":443"
		srv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
		LogInfof("Started runing HTTPS on %s\n", srv.Addr)
		go func() {
			srv.ListenAndServeTLS("", "")
		}()
	}

	srv := makeHTTPServer()
	if IsLinux() {
		httpAddr = ":80"
	}
	srv.Addr = httpAddr
	LogInfof("starting website on %s\n", httpAddr)
	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("http.ListendAndServer() failed with %s\n", err)
	}
}
