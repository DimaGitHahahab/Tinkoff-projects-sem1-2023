package server

import "net/http"

type BasicAuthCredentials struct {
	Username string
	Password string
}

// basicAuthMiddleware checks if request has basic auth credentials.
func basicAuthMiddleware(h http.Handler, auth BasicAuthCredentials) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		user, pass, ok := request.BasicAuth()
		if !ok || user != auth.Username || pass != auth.Password {
			http.Error(responseWriter, "Unauthorized", http.StatusUnauthorized)
			return

		}
		h.ServeHTTP(responseWriter, request)
	})
}
