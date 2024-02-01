package router

import (
	"homework/internal/handler"
	"net/http"
)

func NewRouter(h *handler.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/device", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.HandleCreate(w, r)
		case http.MethodGet:
			h.HandleGet(w, r)
		case http.MethodPut:
			h.HandleUpdate(w, r)
		case http.MethodDelete:
			h.HandleDelete(w, r)
		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
