package main

import (
	"fmt"
	"homework/internal/client"
	"log"
)

func main() {
	// example of usage
	cbConfig := client.BreakerConfig{MaxRequests: 10, Timeout: 3}
	c, err := client.NewClient().WithLogging().WithCircuitBreaker(cbConfig).Build()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := c.Get("https://google.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.StatusCode)
}
