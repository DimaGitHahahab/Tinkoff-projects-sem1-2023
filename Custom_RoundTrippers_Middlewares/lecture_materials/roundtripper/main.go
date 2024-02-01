package main

import (
	"fmt"
	"net/http"
)

type logRoundTripper struct {
	next http.RoundTripper
}

func (l logRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	fmt.Println("prepare to send http request")
	return l.next.RoundTrip(r)
}

func main() {
	client := http.Client{
		Transport: logRoundTripper{
			next: http.DefaultTransport,
		},
	}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/device", http.NoBody)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("status code: %d\n", resp.StatusCode)
}
