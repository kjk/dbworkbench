package main

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	httpAddr = ":5555"
)

func logInfof(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	fmt.Print(s)
}

func logErrorf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	fmt.Print(s)
}

// url: /api/macupdatecheck
func handleWinUpdateCheck(w http.ResponseWriter, r *http.Request) {
}

// url: /api/macupdatecheck
func handleMacUpdateCheck(w http.ResponseWriter, r *http.Request) {
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
		http.Redirect(w, r, "for-mac.html", http.StatusFound /* 302 */)
		return
	}

	// for windows and everything else
	http.Redirect(w, r, "for-windows.html", http.StatusFound /* 302 */)
}

// url: /*
func handleIndex(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	if uri == "/" {
		redirectIndex(w, r)
		return
	}

}

func initHandlers() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/winupdatecheck", handleWinUpdateCheck)
	http.HandleFunc("/api/macupdatecheck", handleMacUpdateCheck)
}

func main() {
	initHandlers()
	logInfof("starting website on %s\n", httpAddr)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		logErrorf("http.ListendAndServe() failed with '%s'\n", err)
	}

}
