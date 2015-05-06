package main

import "net/http"

// CountingResponseWriter wraps http.ResponseWriter and counts how many
// bytes were written and which code was sent
type CountingResponseWriter struct {
	w            http.ResponseWriter
	bytesWritten int
	code         int
}

func NewCountingResponseWriter(w http.ResponseWriter) *CountingResponseWriter {
	return &CountingResponseWriter{
		w: w,
	}
}

func (w *CountingResponseWriter) Write(d []byte) (int, error) {
	n, err := w.w.Write(d)
	w.bytesWritten += n
	return n, err
}

func (w *CountingResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *CountingResponseWriter) WriteHeader(code int) {
	w.code = code
	w.w.WriteHeader(code)
}
