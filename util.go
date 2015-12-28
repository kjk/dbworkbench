package main

import (
	"log"
	"mime"
	"net/http"
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
