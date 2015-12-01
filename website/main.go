package main

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	httpAddr = ":5555"
	// for auto-update
	latestMacVersion     = "0.1"
	latestMacDownloadURL = ""
	latestWinVersion     = "0.2"
	latestWinDownloadURL = "https://kjkpub.s3.amazonaws.com/software/databaseworkbench/rel/DatabaseWorkbench-setup-0.1.exe"
)

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

// url: /api/winupdatecheck?ver=${ver}
func handleWinUpdateCheck(w http.ResponseWriter, r *http.Request) {
	// TODO: log request for analytics
	// TODO: also this should recive usage data which should be saved
	// for analytics
	LogInfof("handleWinUpdateCheck\n")
	s := fmt.Sprintf(`ver: %s
url: %s`, latestWinVersion, latestWinDownloadURL)
	servePlainText(w, r, 200, s)
}

// url: /api/macupdatecheck?ver=${ver}
func handleMacUpdateCheck(w http.ResponseWriter, r *http.Request) {
	// TODO: log request for analytics
	// TODO: also this should recive usage data which should be saved
	// for analytics
	LogInfof("handleMacUpdateCheck\n")
	s := fmt.Sprintf(`ver: %s
url: %s`, latestMacVersion, latestMacDownloadURL)
	servePlainText(w, r, 200, s)
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
		http.Redirect(w, r, "/s/for-mac.html", http.StatusFound /* 302 */)
		return
	}

	// for windows and everything else
	http.Redirect(w, r, "/s/for-windows.html", http.StatusFound /* 302 */)
}

// url: /
func handleIndex(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	LogInfof("handleIndex: '%s'\n", uri)
	if uri == "/" {
		redirectIndex(w, r)
		return
	}
	http.NotFound(w, r)
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
	LogInfof("starting website on %s\n", httpAddr)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		LogErrorf("http.ListendAndServe() failed with '%s'\n", err)
	}
}
