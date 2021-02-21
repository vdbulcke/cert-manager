package certificate

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// getFingerprintFromRequest return uuid from request
func getFingerprintFromRequest(r *http.Request) (string, error) {
	// parse the id from the url
	vars := mux.Vars(r)

	// convert the id
	id := vars["id"]

	// convert sha256 to lower case without ':' in hex
	fingerprint := convertSHA256(id)
	// check
	err := validateSHA256(fingerprint)
	if err != nil {
		return "", nil
	}

	return fingerprint, nil
}

func convertSHA256(sha256 string) string {
	fingerprint := strings.ReplaceAll(sha256, ":", "")
	return strings.ToLower(fingerprint)
}

func validateSHA256(sha256 string) error {

	// check length
	if len(sha256) != 64 {
		return fmt.Errorf(fmt.Sprintf("length must be 64 found %d", len(sha256)))
	}

	// check if valid Hex string
	_, err := hex.DecodeString(sha256)
	if err != nil {
		return err
	}

	return nil

}
