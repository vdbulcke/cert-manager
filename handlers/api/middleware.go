package api

import (
	"net/http"
)

// API Object
type API struct{}

// NewAPI return an API object
func NewAPI() *API {
	return &API{}
}

// CommonAPIMiddleware common handler for API calls
func (api *API) CommonAPIMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
