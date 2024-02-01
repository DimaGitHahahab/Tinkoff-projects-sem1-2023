package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println(http.Client{
		Transport: newRateLimitRoundTripper(config{}, http.DefaultTransport),
	})
}
