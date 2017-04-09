package main

import (
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strings"
)

// Request.RemoteAddress contains port, which we want to remove i.e.:
// "[::1]:58292" => "[::1]"
func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

func getIPFromRequest(r *http.Request) string {
	hdr := r.Header
	hdrReal := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrReal == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: should return first non-local address
		return parts[0]
	}
	return hdrReal
}

// PathExists returns true if path exists
// treats any error (e.g. lack of access due to permissions) as non-existence
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsLinux returns true if we're running on linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// IsMac returns true if we're running on mac
func IsMac() bool {
	return runtime.GOOS == "darwin"
}

// UserHomeDir returns $HOME dir
func UserHomeDir() string {
	// user.Current() returns nil if cross-compiled e.g. on mac for linux
	if usr, _ := user.Current(); usr != nil {
		return usr.HomeDir
	}
	return os.Getenv("HOME")
}

// ExpandTildeInPath expands ~ in path
func ExpandTildeInPath(s string) string {
	if strings.HasPrefix(s, "~") {
		return UserHomeDir() + s[1:]
	}
	return s
}
