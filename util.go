package main

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func fatalIfErr(err error, what string) {
	if err != nil {
		log.Fatalf("%s failed with %s\n", what, err)
	}
}

func isMac() bool {
	return runtime.GOOS == "darwin"
}

func isLinux() bool {
	return runtime.GOOS == "linux"
}

func isWindows() bool {
	return runtime.GOOS == "windows"
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

// IntInArray returns true if int is in array of ints
func IntInArray(arr []int, n int) bool {
	for _, n2 := range arr {
		if n2 == n {
			return true
		}
	}
	return false
}

// IntAppendIfNotExists adds int to array if it's not in the array yet
func IntAppendIfNotExists(arr []int, n int) []int {
	if IntInArray(arr, n) {
		return arr
	}
	return append(arr, n)
}

func getMyHost(r *http.Request) string {
	return "http://" + r.Host
}

func openDefaultBrowserMac(uri string) error {
	cmd := exec.Command("open", uri)
	return cmd.Run()
}

// https://github.com/skratchdot/open-golang/blob/master/open/exec_windows.go

func opneDefaultBrowserWin2(uri string) error {
	return exec.Command("cmd", "/c", "start", uri).Run()
}

func openDefaultBrowserWin(uri string) error {
	runDll32 := filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")
	return exec.Command(runDll32, "url.dll,FileProtocolHandler", uri).Run()
}

func openDefaultBrowser(uri string) error {
	var err error
	if isMac() {
		err = openDefaultBrowserMac(uri)
	} else if isWindows() {
		err = openDefaultBrowserWin(uri)
	} else {
		err = fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	return err
}
