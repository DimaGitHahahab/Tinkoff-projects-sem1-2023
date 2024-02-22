package client

import (
	"net/http"

	"github.com/sony/gobreaker"
)

type Client struct {
	log     bool
	breaker *gobreaker.CircuitBreaker
	*http.Client
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) WithLogging() *Client {
	c.log = true
	return c
}

func (c *Client) WithCircuitBreaker(cfg BreakerConfig) *Client {
	c.breaker = newCircuitBreaker(cfg)
	return c
}

func (c *Client) Build() (*Client, error) {
	transport := http.DefaultTransport

	if c.breaker != nil {
		transport = NewBreakerRoundTripper(transport, c.breaker)
	}

	if c.log {
		transport = NewLoggingRoundTripper(transport)
	}

	c.Client = &http.Client{
		Transport: transport,
	}

	return c, nil
}
