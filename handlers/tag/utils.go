package tag

import (
	"net/http"

	"github.com/gorilla/mux"
)

// getTagNameFromRequest return uuid from request
func getTagNameFromRequest(r *http.Request) (string, error) {
	// parse the id from the url
	vars := mux.Vars(r)

	return vars["name"], nil
}
