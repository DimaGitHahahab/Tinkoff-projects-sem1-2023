package format

import (
	"fmt"
	"net/http"
	"strings"
)

// RequestLog formats request data to string for logging.
func RequestLog(r *http.Request) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("method: %v, endpoint: %v\n", r.Method, r.URL.Path))
	if len(r.Header) != 0 {
		sb.WriteString("headers:\n")
		sb.WriteString(mapStrToStrSlice(r.Header))
	}
	queryParams := r.URL.Query()
	if len(queryParams) != 0 {
		sb.WriteString("query params:\n")
		sb.WriteString(mapStrToStrSlice(queryParams))
	}
	return sb.String()
}

// ResponseLog formats response data to string for logging.
func ResponseLog(code int, header http.Header, err error) string {
	var sb strings.Builder
	if len(header) != 0 {
		sb.WriteString("headers:\n")
		sb.WriteString(mapStrToStrSlice(header))
	}
	if err != nil {
		sb.WriteString(err.Error())
	}
	sb.WriteString(fmt.Sprintf("status code: %v\n", code))
	return sb.String()
}

// mapStrToStrSlice formats http.Header or url.Values (query params) to string.
func mapStrToStrSlice(data map[string][]string) string {
	var sb strings.Builder
	for name, values := range data {
		for _, value := range values {
			sb.WriteString(fmt.Sprintf("%v : %v\n", name, value))
		}
	}
	return sb.String()
}
