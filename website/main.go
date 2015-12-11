package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/kjk/u"
)

var (
	httpAddr = ":5555"

	// for auto-update
	latestMacVersion = "0.1.1"
	latestWinVersion = "0.1.1"

	dataDir string
)

func latestMacDownloadURL() string {
	return fmt.Sprintf("https://kjkpub.s3.amazonaws.com/software/dbhero/rel/DBHero-%s.zip", latestMacVersion)
}

func latestWinDownloadURL() string {
	return fmt.Sprintf("https://kjkpub.s3.amazonaws.com/software/dbhero/rel/DBHero-setup-%s.exe", latestWinVersion)
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
	if u.PathExists(dataDir) {
		return dataDir
	}
	dataDir = u.ExpandTildeInPath("~/data/dbhero-website")
	if u.PathExists(dataDir) {
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
		return nil
	}
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		LogErrorf("ioutil.ReadAll() failed with '%s'\n", err)
		return nil
	}
	if len(d) == 0 {
		return nil
	}
	return d
}

func prependVerIPTime(d []byte, ip, ver string) []byte {
	now := time.Now()
	t := now.Unix()
	day := now.Format("2006-01-02")
	s := fmt.Sprintf("time: %d\nday: %s\nip: %s\nver_in_url: %s\n", t, day, ip, ver)
	return append([]byte(s), d...)
}

// url: /api/winupdatecheck?ver=${ver}
func handleWinUpdateCheck(w http.ResponseWriter, r *http.Request) {
	LogInfof("handleWinUpdateCheck\n")

	d := readRawPostData(r)
	ver := r.FormValue("ver")
	ip := getIPFromRequest(r)

	s := fmt.Sprintf(`ver: %s
url: %s`, latestWinVersion, latestWinDownloadURL())
	servePlainText(w, r, 200, s)

	d = prependVerIPTime(d, ip, ver)
	recordUsage(d)
}

// url: /api/macupdatecheck?ver=${ver}
func handleMacUpdateCheck(w http.ResponseWriter, r *http.Request) {
	LogInfof("handleMacUpdateCheck\n")

	d := readRawPostData(r)
	ver := r.FormValue("ver")
	ip := getIPFromRequest(r)

	s := fmt.Sprintf(`ver: %s
url: %s`, latestMacVersion, latestMacDownloadURL())
	servePlainText(w, r, 200, s)

	d = prependVerIPTime(d, ip, ver)
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

func redirectIndex2(w http.ResponseWriter, r *http.Request) {
	ua := r.UserAgent()
	if isMacUserAgent(ua) {
		http.Redirect(w, r, "for-mac", http.StatusFound /* 302 */)
		return
	}

	// for windows and everything else
	http.Redirect(w, r, "for-windows", http.StatusFound /* 302 */)
}

func redirectIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/s/index.html", http.StatusFound /* 302 */)
}

// url: /
func handleIndex(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	LogInfof("handleIndex: '%s'\n", uri)
	if uri == "/" {
		redirectIndex(w, r)
		return
	}
	// map /foo and /foo.html to /s/foo.html such file exists
	name := uri[1:]
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

// url: /s/:path
func handleStatic(w http.ResponseWriter, r *http.Request) {
	LogInfof("handleStatic: '%s'\n", r.URL.Path)
	name := r.URL.Path[len("/s/"):]
	resourcePath := filepath.Join("www", name)
	serveStatic(w, r, resourcePath)
}

func initHandlers() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/s/", handleStatic)
	http.HandleFunc("/api/winupdatecheck", handleWinUpdateCheck)
	http.HandleFunc("/api/macupdatecheck", handleMacUpdateCheck)
}

func main() {
	initHandlers()
	// TODO: open log files
	openUsageFileMust()
	LogInfof("starting website on %s\n", httpAddr)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		LogErrorf("http.ListendAndServe() failed with '%s'\n", err)
	}
}
