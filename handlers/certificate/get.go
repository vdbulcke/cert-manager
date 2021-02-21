package certificate

import (
	"fmt"
	"net/http"

	"github.com/vdbulcke/cert-manager/data"
	"github.com/vdbulcke/cert-manager/handlers/api"
)

// swagger:route GET /certificate/GetCertificateByID/{id} Certificate GetCertificateByID
// Return certificate from the database
// responses:
//	200: certificateResponse
//	404: errorResponse
//  400: errorResponse
//  500: errorResponse

// GetCertificateByID handles GET requests
func (h *APICertificateHandler) GetCertificateByID(rw http.ResponseWriter, r *http.Request) *api.APIError {

	uuid, err := api.GetIDFromRequest(r)
	if err != nil {
		return &api.APIError{
			Code:    http.StatusBadRequest,
			Type:    &api.ValidationError{},
			Message: fmt.Sprintf("Cannot find id from request %s", r.URL.String()),
		}
	}

	// lookup this certificate
	cert, err := h.certBackend.GetCertByID(uuid)
	if err != nil {
		if _, ok := err.(*data.DBObjectNotFound); ok {
			h.logger.Debug("GetCertificateByID: object not found", "uuid", uuid)

			return &api.APIError{
				Err:     err,
				Code:    http.StatusNotFound,
				Type:    &api.ObjectNotFoundError{},
				Message: err.Error(),
			}

		}

		// default
		h.logger.Debug("GetCertificateByID: unexpected error searching for certificate", "uuid", uuid, "err", err)

		return &api.APIError{
			Err:     err,
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: fmt.Sprintf("Error searching for certificate id=%s", uuid.String()),
		}

	}

	h.logger.Debug("GetCertificateByID: Found cert", "cert", cert)

	// Write Status code
	rw.WriteHeader(http.StatusOK)
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

// swagger:route GET /certificate/GetCertificateByFingerprint/{id} Certificate GetCertificateByFingerprint
// Return certificate from the database
// responses:
//	200: certificateResponse
//	404: errorResponse
//  400: errorResponse
//  500: errorResponse

// GetCertificateByFingerprint handles GET requests
func (h *APICertificateHandler) GetCertificateByFingerprint(rw http.ResponseWriter, r *http.Request) *api.APIError {

	fingerprint, err := getFingerprintFromRequest(r)
	if err != nil {
		return &api.APIError{
			Err:     err,
			Code:    http.StatusBadRequest,
			Type:    &api.ValidationError{},
			Message: fmt.Sprintf("Invalid SHA256 fingerprint found in request %s", r.URL.String()),
		}
	}

	// lookup this certificate
	cert, err := h.certBackend.GetCertByFingerprint(fingerprint)

	if err != nil {
		if _, ok := err.(*data.DBObjectNotFound); ok {
			h.logger.Debug("GetCertificateByFingerprint: object not found", "fingerprint", fingerprint)

			return &api.APIError{
				Err:     err,
				Code:    http.StatusNotFound,
				Type:    &api.ObjectNotFoundError{},
				Message: err.Error(),
			}

		}

		// default
		h.logger.Debug("GetCertificateByFingerprint: unexpected error searching for certificate", "fingerprint", fingerprint, "err", err)

		return &api.APIError{
			Err:     err,
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: fmt.Sprintf("Error searching for certificate id=%s", fingerprint),
		}

	}

	h.logger.Debug("GetCertificateByFingerprint: Found cert", "cert", cert)

	// Write Status code
	rw.WriteHeader(http.StatusOK)
	err = api.ToJSON(cert, rw)
	if err != nil {
		// we should never be here but log the error just incase
		h.logger.Error("GetCertificateByFingerprint: Error Serializing JSON", "cert", cert, "err", err)
		return &api.APIError{
			Err:     err,
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: "error formating json",
		}
	}

	return nil
}

// swagger:route GET /certificate/ListCerts Certificate ListCerts
// Return a list of Certificates from the database
// responses:
//	200: certificateListResponse
//  404: errorResponse
//  500: errorResponse

// ListCerts handles GET requests
func (h *APICertificateHandler) ListCerts(rw http.ResponseWriter, r *http.Request) *api.APIError {

	// lookup all certificate
	certs, err := h.certBackend.ListCerts()
	if err != nil {
		if _, ok := err.(*data.DBObjectNotFound); ok {
			h.logger.Debug("ListCerts: object not found")

			return &api.APIError{
				Err:     err,
				Code:    http.StatusNotFound,
				Type:    &api.ObjectNotFoundError{},
				Message: err.Error(),
			}

		}
		// default
		h.logger.Debug("ListCerts: unexpected error searching for certificate", "err", err)

		return &api.APIError{
			Err:     err,
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: fmt.Sprintf("Error searching for certificate"),
		}

	}

	h.logger.Debug("ListCerts: Found certs", "certs", certs)

	// Write Status code
	rw.WriteHeader(http.StatusOK)
	err = api.ToJSON(certs, rw)
	if err != nil {
		// we should never be here but log the error just incase
		h.logger.Error("ListCerts: Error Serializing JSON", "certs", certs, "err", err)
		return &api.APIError{
			Err:     err,
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: "error formating json",
		}
	}

	return nil
}
