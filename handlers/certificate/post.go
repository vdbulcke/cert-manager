package certificate

import (
	"net/http"

	"github.com/vdbulcke/cert-manager/data"
	"github.com/vdbulcke/cert-manager/handlers/api"
)

// TODO:
// swagger:route GET /certificate/CreateCertificate Certificate CreateCertificate
// Return a list of tenants from the database
// responses:
//	200: certificateResponse
//	404: errorResponse
//  401: errorResponse

// CreateCertificate handles GET requests
func (h *APICertificateHandler) CreateCertificate(rw http.ResponseWriter, r *http.Request) *api.APIError {

	// look up cert input (set by the middleware) in request context
	certInput := r.Context().Value(APICertificateKey{}).(APICertificateInput)

	// checking certInput contains a pem
	if certInput.Pem == "" {
		return &api.APIError{
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: "Error parsing input",
		}
	}

	// setting default status code
	statusCode := http.StatusCreated
	var cert *data.Certificate
	var err error

	// checking if the input contains tags
	if len(certInput.Tags) != 0 {
		cert, err = h.certBackend.CreateCertificateWithTags(certInput.Pem, certInput.Tags)
		if err != nil {
			if _, ok := err.(*data.DBObjectAlreadyExist); ok {
				statusCode = http.StatusConflict
			} else if _, ok := err.(*data.DBObjectValidationError); ok {
				return &api.APIError{
					Err:     err,
					Code:    http.StatusBadRequest,
					Type:    &api.ValidationError{},
					Message: err.Error(),
				}

			} else {
				// default
				return &api.APIError{Err: err,
					Code:    http.StatusInternalServerError,
					Type:    &api.InternalServerError{},
					Message: "error creating cert",
				}
			}

		}

	} else {
		// creating a cert without tags
		cert, err = h.certBackend.CreateCertificate(certInput.Pem)
		if err != nil {
			if _, ok := err.(*data.DBObjectAlreadyExist); ok {
				statusCode = http.StatusConflict
			} else if _, ok := err.(*data.DBObjectValidationError); ok {
				return &api.APIError{
					Err:     err,
					Code:    http.StatusBadRequest,
					Type:    &api.ValidationError{},
					Message: err.Error(),
				}

			} else {
				// default
				return &api.APIError{Err: err,
					Code:    http.StatusInternalServerError,
					Type:    &api.InternalServerError{},
					Message: "error creating cert",
				}
			}

		}
	}

	h.logger.Debug("GetCertificateByID: Found cert", "cert", cert)

	// Write Status code
	rw.WriteHeader(statusCode)
	err = api.ToJSON(cert, rw)
	if err != nil {
		// we should never be here but log the error just incase
		h.logger.Error("GetCertificateByID: Error Serializing JSON", "cert", cert, "err", err)
		return &api.APIError{
			Err:     err,
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: "error formating json",
		}
	}

	return nil
}

// Helper Functions
