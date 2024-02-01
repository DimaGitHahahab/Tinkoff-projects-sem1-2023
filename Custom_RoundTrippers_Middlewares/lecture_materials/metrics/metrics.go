package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var reqSumMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "http",
	Subsystem: "server",
	Name:      "requests_seconds",
}, []string{"uri", "method", "status"})

var respBytesMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "http",
	Subsystem: "server",
	Name:      "response_bytes",
	Buckets:   []float64{64, 128, 256, 512, 1024, 2048, 4096, 8192},
}, []string{"uri", "method", "status"})

func ServerMetricsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		now := time.Now()

		w := &statusWriter{
			status: 200,
			w:      writer,
		}

		h.ServeHTTP(w, request)

		labels := prometheus.Labels{
			"uri":    request.URL.Path,
			"method": request.Method,
			"status": fmt.Sprintf("%d", w.status),
		}
		reqSumMetric.With(labels).Observe(time.Since(now).Seconds())
		respBytesMetric.With(labels).Observe(float64(w.count.Load()))
	})
}

type statusWriter struct {
	w      http.ResponseWriter
	status int
	count  atomic.Int64
}

func (s *statusWriter) Header() http.Header {
	return s.w.Header()
}

func (s *statusWriter) Write(bytes []byte) (int, error) {
	n, err := s.w.Write(bytes)
	s.count.Add(int64(n))
	return n, err
}

func (s *statusWriter) WriteHeader(statusCode int) {
	s.status = statusCode
	s.w.WriteHeader(statusCode)
}
