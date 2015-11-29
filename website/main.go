package main

import (
	"fmt"
	"net/http"
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

func initHandlers() {
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
