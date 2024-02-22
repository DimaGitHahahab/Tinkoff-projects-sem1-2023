package client

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

const (
	defaultTimeout  = 5
	defaultInterval = 3
	defaultRequests = 3
)

type BreakerConfig struct {
	MaxRequests int
	Interval    int
	Timeout     int
}

func newCircuitBreaker(cfg BreakerConfig) *gobreaker.CircuitBreaker {
	timeout := defaultTimeout
	interval := defaultInterval
	maxRequests := defaultRequests

	if cfg.Timeout != 0 {
		timeout = cfg.Timeout
	}
	if cfg.Interval != 0 {
		interval = cfg.Interval
	}
	if cfg.MaxRequests != 0 {
		maxRequests = cfg.MaxRequests
	}

	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "HTTPClient",
		MaxRequests: uint32(maxRequests),
		Interval:    time.Duration(interval) * time.Second,
		Timeout:     time.Duration(timeout) * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 2
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit breaker state change: %s -> %s\n", from, to)
		},
	})
}

// breakerRoundTripper wraps http.RoundTripper with gobreaker.CircuitBreaker.
type breakerRoundTripper struct {
	next    http.RoundTripper
	breaker *gobreaker.CircuitBreaker
}

func (b breakerRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	result, err := b.breaker.Execute(func() (interface{}, error) {
		return b.next.RoundTrip(request)
	})

	if err != nil {
		return nil, err

	}

	resp, ok := result.(*http.Response)
	if !ok {
		return nil, errors.New("unexpected response type from circuit breaker")
	}

	return resp, nil
}

func NewBreakerRoundTripper(next http.RoundTripper, breaker *gobreaker.CircuitBreaker) http.RoundTripper {
	return breakerRoundTripper{
		next:    next,
		breaker: breaker,
	}
}
