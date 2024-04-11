package logger

import "net/http"

type logResponseWriter struct {
	http.ResponseWriter
	wroteHeader bool
	statusCode  int
}

func (rw *logResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.wroteHeader = true
}

func (rw *logResponseWriter) StatusCode() int {
	if rw.wroteHeader {
		return rw.statusCode
	}
	return http.StatusOK
}

func newLogResponseWriter(rw http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{
		ResponseWriter: rw,
		statusCode:     0,
		wroteHeader:    false,
	}
}
