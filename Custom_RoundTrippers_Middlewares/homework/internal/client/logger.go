package client

import (
	"homework/internal/format"
	"log"
	"net/http"
)

// loggingRoundTripper logs request and response data.
type loggingRoundTripper struct {
	next http.RoundTripper
}

func (l loggingRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	log.Println("request:" + format.RequestLog(request))
	resp, err := l.next.RoundTrip(request)
	if err != nil {
		log.Println("response error:", err)
	} else {
		log.Println("response: " + format.ResponseLog(resp.StatusCode, resp.Header, err))
	}
	return resp, err
}

func NewLoggingRoundTripper(next http.RoundTripper) http.RoundTripper {
	return loggingRoundTripper{next: next}
}
