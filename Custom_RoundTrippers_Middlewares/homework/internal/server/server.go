package server

import (
	"net/http"
)

type Config struct {
	Addr string
	Mux  *http.ServeMux
	Log  bool
	Cred BasicAuthCredentials
}

// NewServer returns configured http.Server
func NewServer(cfg Config, middlewares ...func(handler http.Handler) http.Handler) *http.Server {
	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}
	if cfg.Mux == nil {
		cfg.Mux = http.NewServeMux()
	}

	var handler http.Handler = cfg.Mux

	for _, customMiddleware := range middlewares {
		handler = customMiddleware(handler)
	}

	if cfg.Cred.Password != "" && cfg.Cred.Username != "" {
		handler = basicAuthMiddleware(handler, cfg.Cred)
	}

	if cfg.Log {
		handler = loggingMiddleware(handler)
	}

	return &http.Server{
		Addr:    cfg.Addr,
		Handler: handler,
	}
}
