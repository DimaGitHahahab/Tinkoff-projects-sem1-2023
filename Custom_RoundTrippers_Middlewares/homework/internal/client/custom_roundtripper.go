package client

import "net/http"

// customRoundTripper allows chaining multiple http.RoundTripper.
type customRoundTripper struct {
	next   http.RoundTripper
	custom http.RoundTripper
}

func (c customRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	resp, err := c.custom.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	if c.next != nil {
		return c.next.RoundTrip(request)
	} else {
		return resp, nil
	}
}
