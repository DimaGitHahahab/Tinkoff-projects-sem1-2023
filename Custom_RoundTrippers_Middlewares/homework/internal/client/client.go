package client

import (
	"net/http"
	"os"
	"strconv"
)

type Client struct {
	*http.Client
}
type Config struct {
	Log     bool
	Breaker BreakerConfig
}

func NewClient(cfg Config) (*Client, error) {
	transport := http.DefaultTransport

	cb, err := configureCircuitBreaker(cfg.Breaker)
	if err != nil {
		return nil, err
	}
	if cb != nil {
		transport = NewBreakerRoundTripper(transport, cb)
	}

	if cfg.Log {
		transport = NewLoggingRoundTripper(transport)
	}

	return &Client{
		Client: &http.Client{
			Transport: transport,
		},
	}, nil

}

func ReadConf() (Config, error) {
	cfg := Config{}
	if os.Getenv("LOG") == "true" {
		cfg.Log = true
	}

	if os.Getenv("MAX_REQUESTS") != "" {
		maxReq, err := strconv.Atoi(os.Getenv("MAX_REQUESTS"))
		if err != nil {
			return cfg, err
		}
		cfg.Breaker.MaxRequests = maxReq
	}

	if os.Getenv("INTERVAL") != "" {
		interval, err := strconv.Atoi(os.Getenv("INTERVAL"))
		if err != nil {
			return cfg, err
		}
		cfg.Breaker.Interval = interval
	}

	if os.Getenv("TIMEOUT") != "" {
		timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
		if err != nil {
			return cfg, err
		}
		cfg.Breaker.Timeout = timeout
	}

	return cfg, nil
}
