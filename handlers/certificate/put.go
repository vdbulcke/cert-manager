package certificate

import (
	"fmt"
	"net/http"

	"github.com/vdbulcke/cert-manager/data"
	"github.com/vdbulcke/cert-manager/handlers/api"
)

// swagger:route PUT /certificate/UpdateCertificateTag/{id} Certificate UpdateCertificateTag
// Return updated certificate
// responses:
//	202: certificateResponse
//	404: errorResponse
//  401: errorResponse
//  400: errorResponse
//  500: errorResponse

// UpdateCertificateTag handles PUT requests
func (h *APICertificateHandler) UpdateCertificateTag(rw http.ResponseWriter, r *http.Request) *api.APIError {

	// getting uuid
	uuid, err := api.GetIDFromRequest(r)
	if err != nil {
		return &api.APIError{
			Code:    http.StatusBadRequest,
			Type:    &api.ValidationError{},
			Message: fmt.Sprintf("Cannot find id from request %s", r.URL.String()),
		}
	}

	// look up cert input (set by the middleware) in request context
	certTagInput, ok := r.Context().Value(APICertificateTagKey{}).(APICertificateTagInput)
	if !ok {
		return &api.APIError{
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: "Error parsing input",
		}
	}

	// checking certInput contains a pem
	if len(certTagInput.Tags) == 0 {

		return &api.APIError{
			Code:    http.StatusBadRequest,
			Type:    &api.ValidationError{},
			Message: "Empty list",
		}
	}

	// update certificate with new tags
	cert, err := h.certBackend.SetCertTagsNameByID(uuid, certTagInput.Tags)
	if err != nil {
		if _, ok := err.(*data.DBObjectNotFound); ok {
			return &api.APIError{
				Err:     err,
				Code:    http.StatusNotFound,
				Type:    &api.ObjectNotFoundError{},
				Message: err.Error(),
			}

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

	h.logger.Debug("GetCertificateByID: Found cert", "cert", cert)

	// Write Status code
	rw.WriteHeader(http.StatusAccepted)
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
