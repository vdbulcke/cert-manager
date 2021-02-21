package certificate

import (
	"fmt"
	"net/http"

	"github.com/vdbulcke/cert-manager/data"
	"github.com/vdbulcke/cert-manager/handlers/api"
)

// swagger:route DELETE /certificate/DeleteCertificateTagsByID/{id} Certificate DeleteCertificateTagsByID
// Return Delete Tags from Certificate
// responses:
//	202: certificateResponse
//	404: errorResponse
//  401: errorResponse
//  400: errorResponse
//  500: errorResponse

// DeleteCertificateTagsByID handles DELETE requests
func (h *APICertificateHandler) DeleteCertificateTagsByID(rw http.ResponseWriter, r *http.Request) *api.APIError {

	// Get uuid from request
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

	// lookup this certificate
	cert, err := h.certBackend.DeleteCertificateTagsByID(uuid, certTagInput.Tags)
	if err != nil {
		if _, ok := err.(*data.DBObjectNotFound); ok {
			h.logger.Debug("DeleteCertificateByID: object not found", "uuid", uuid)

			return &api.APIError{
				Err:     err,
				Code:    http.StatusNotFound,
				Type:    &api.ObjectNotFoundError{},
				Message: err.Error(),
			}

		}

		// default
		h.logger.Debug("DeleteCertificateByID: unexpected error searching for certificate", "uuid", uuid, "err", err)

		return &api.APIError{
			Err:     err,
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: fmt.Sprintf("Error Deleting certificate id=%s", uuid.String()),
		}

	}

	// Write Status code
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

// swagger:route DELETE /certificate/DeleteCertificateByID/{id} Certificate DeleteCertificateByID
// Deletes certificate from DB
// responses:
//	204: noContentResponse
//	404: errorResponse
//  401: errorResponse
//  500: errorResponse

// DeleteCertificateByID handles GET requests
func (h *APICertificateHandler) DeleteCertificateByID(rw http.ResponseWriter, r *http.Request) *api.APIError {

	// Get uuid from request
	uuid, err := api.GetIDFromRequest(r)
	if err != nil {
		return &api.APIError{
			Code:    http.StatusBadRequest,
			Type:    &api.ValidationError{},
			Message: fmt.Sprintf("Cannot find id from request %s", r.URL.String()),
		}
	}

	// lookup this certificate
	err = h.certBackend.DeleteCertByID(uuid)
	if err != nil {
		if _, ok := err.(*data.DBObjectNotFound); ok {
			h.logger.Debug("DeleteCertificateByID: object not found", "uuid", uuid)

			return &api.APIError{
				Err:     err,
				Code:    http.StatusNotFound,
				Type:    &api.ObjectNotFoundError{},
				Message: err.Error(),
			}

		}

		// default
		h.logger.Debug("DeleteCertificateByID: unexpected error searching for certificate", "uuid", uuid, "err", err)

		return &api.APIError{
			Err:     err,
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: fmt.Sprintf("Error Deleting certificate id=%s", uuid.String()),
		}

	}

	// Write Status code
	rw.WriteHeader(http.StatusNoContent)

	return nil
}
