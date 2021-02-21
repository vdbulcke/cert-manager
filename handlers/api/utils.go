package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GetIDFromRequest return uuid from request
func GetIDFromRequest(r *http.Request) (uuid.UUID, error) {
	// parse the id from the url
	vars := mux.Vars(r)

	// convert the id into an uuid and return
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		// should never happen
		return id, err
	}

	return id, nil
}

// ToJSON serializes the given interface into a string based JSON format
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)

	return e.Encode(i)
}

// FromJSON deserializes the object from JSON string
// in an io.Reader to the given interface
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}
