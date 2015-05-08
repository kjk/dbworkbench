package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

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

func serveError(w http.ResponseWriter, r *http.Request, isJSON bool, errMsg string) {
	if isJSON {
		serveJSONError(w, r, errMsg)
		return
	}
	LogErrorf("uri: '%s', err: '%s'\n", r.RequestURI, errMsg)
	http.NotFound(w, r)
}

func writeHeader(w http.ResponseWriter, code int, contentType string) {
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	w.WriteHeader(code)
}

// err can be an error, a string or anything that can be converted to string
func serveJSONError(w http.ResponseWriter, r *http.Request, errMsg interface{}) {
	writeHeader(w, 400, "application/json") // Note: maybe different code, like 500?
	msg := fmt.Sprintf("%s", errMsg)
	LogErrorf("url: '%s', err: '%s'\n", r.RequestURI, msg)
	v := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(v); err != nil {
		LogErrorf("err: %s\n", err)
	}
}

func serveJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	writeHeader(w, 200, "application/json")
	encoder := json.NewEncoder(w)
	return encoder.Encode(v)
}

func serveJSONP(w http.ResponseWriter, r *http.Request, v interface{}, jsonp string) error {
	if jsonp == "" {
		return serveJSON(w, r, v)
	}

	writeHeader(w, 200, "application/json")
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

func servePlainText(w http.ResponseWriter, r *http.Request, code int, format string, args ...interface{}) error {
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
