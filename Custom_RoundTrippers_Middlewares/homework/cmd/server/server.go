package main

import (
	"homework/internal/server"
	"net/http"
)

func main() {
	cfg := server.Config{
		Addr: ":8080",
		Mux:  http.NewServeMux(),
		Log:  true,
		Cred: server.BasicAuthCredentials{
			Username: "user",
			Password: "pass",
		},
	}

	s := server.NewServer(cfg)

	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}

}
