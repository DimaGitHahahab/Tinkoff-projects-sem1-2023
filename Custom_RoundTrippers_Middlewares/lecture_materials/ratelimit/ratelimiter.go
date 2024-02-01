package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/time/rate"
)

type config struct {
	Max      int
	Burst    int
	Duration time.Duration
}

var metric = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "http",
		Subsystem: "client",
		Name:      "requests_count",
	},
	[]string{"method"},
)

var errToManyReqeusts = errors.New("too many requests")

type rateLimitRoundTripper struct {
	limiter *rate.Limiter

	next http.RoundTripper
}

func newRateLimitRoundTripper(cfg config, next http.RoundTripper) http.RoundTripper {
	r := rate.Limit(cfg.Max) / rate.Limit(cfg.Duration.Seconds())
	limiter := rate.NewLimiter(r, cfg.Burst)

	return rateLimitRoundTripper{
		limiter: limiter,
		next:    next,
	}
}

func (rl rateLimitRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	metric.With(prometheus.Labels{"method": r.Method}).Inc()

	if !rl.limiter.Allow() {
		return nil, errToManyReqeusts
	}

	return rl.next.RoundTrip(r)
}
