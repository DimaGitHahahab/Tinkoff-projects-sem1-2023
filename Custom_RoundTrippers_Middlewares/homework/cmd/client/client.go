package main

import (
	"fmt"
	"homework/internal/client"
	"log"
)

func main() {
	cfg, err := client.ReadConf()
	if err != nil {
		log.Fatal(err)
	}

	c, err := client.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// example of usage
	resp, err := c.Get("https://google.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.StatusCode)

}
