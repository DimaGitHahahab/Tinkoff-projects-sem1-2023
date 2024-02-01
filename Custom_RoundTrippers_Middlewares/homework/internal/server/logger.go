package server

import (
	"homework/internal/format"
	"log"
	"net/http"
)

// loggingMiddleware logs request and response data.
func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		log.Println("request: " + format.RequestLog(request))
		logRespWriter := &loggingResponseWriter{originalRW: responseWriter, header: http.Header{}, err: nil}
		h.ServeHTTP(logRespWriter, request)
		log.Println("response: " + format.ResponseLog(logRespWriter.code, logRespWriter.header, logRespWriter.err))
	})
}

// loggingResponseWriter stores response headers, http code and error.
type loggingResponseWriter struct {
	originalRW http.ResponseWriter
	header     http.Header
	code       int
	err        error
}

func (rw *loggingResponseWriter) Header() http.Header {
	rw.header = rw.originalRW.Header()
	return rw.originalRW.Header()
}

func (rw *loggingResponseWriter) Write(data []byte) (int, error) {
	size, err := rw.originalRW.Write(data)
	if err != nil {
		rw.err = err
	}
	return size, err
}

func (rw *loggingResponseWriter) WriteHeader(code int) {
	rw.code = code
	rw.originalRW.WriteHeader(code)
}

func NewLoggingMiddleware(h http.Handler) http.Handler {
	return loggingMiddleware(h)
}
