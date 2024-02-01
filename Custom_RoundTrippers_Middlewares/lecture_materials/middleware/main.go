package main

import (
	"fmt"
	"net/http"
)

func getDevice(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get device")
	w.WriteHeader(http.StatusOK)
}

func getStorage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get storage")
	w.WriteHeader(http.StatusOK)
}

func globalMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("global mux middleware")
		h.ServeHTTP(w, r)
	})
}

func deviceMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("local get device middleware")
		h.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	deviceHandler := http.HandlerFunc(getDevice)
	storageHandler := http.HandlerFunc(getStorage)

	mux.Handle("/device", deviceMiddleware(deviceHandler))
	mux.Handle("/storage", storageHandler)

	muxWithMiddleware := globalMiddleware(mux)

	if err := http.ListenAndServe(":8080", muxWithMiddleware); err != nil {
		panic(err)
	}
}
