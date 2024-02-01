package client

import (
	"errors"
	"github.com/sony/gobreaker"
	"log"
	"net/http"
	"time"
)

type BreakerConfig struct {
	MaxRequests int
	Interval    int
	Timeout     int
}

func configureCircuitBreaker(cfg BreakerConfig) (*gobreaker.CircuitBreaker, error) {
	if cfg.Timeout != 0 || cfg.Interval != 0 || cfg.MaxRequests != 0 {
		if cfg.Timeout == 0 {
			cfg.Timeout = 5
		}
		if cfg.Interval == 0 {
			cfg.Interval = 3
		}
		if cfg.MaxRequests == 0 {
			cfg.MaxRequests = 3
		}
		cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        "HTTPClient",
			MaxRequests: uint32(cfg.MaxRequests),
			Interval:    time.Duration(cfg.Interval) * time.Second,
			Timeout:     time.Duration(cfg.Timeout) * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures > 2
			},
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				log.Printf("Circuit breaker state change: %s -> %s\n", from, to)
			},
		})
		return cb, nil
	}
	return nil, nil
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
