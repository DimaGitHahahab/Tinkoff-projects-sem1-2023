package main

import (
	"homework/internal/handler"
	"homework/internal/router"
	"homework/internal/service"
	"log"
	"net"
	"net/http"
	"os"
)

func Address() string {
	host := os.Getenv("HTTP_HOST")
	port := os.Getenv("HTTP_PORT")
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "8080"
	}
	return net.JoinHostPort(host, port)
}

func main() {
	h := handler.NewHandler(service.NewService(service.NewStorage()))
	mux := router.NewRouter(h)

	if err := http.ListenAndServe(Address(), mux); err != nil {
		log.Fatal(err)
	}
}
