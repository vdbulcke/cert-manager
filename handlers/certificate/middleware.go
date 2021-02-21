package certificate

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/vdbulcke/cert-manager/handlers/api"
)

// MiddlewareValidateCertificateInput validates the tenant in the request and calls next if ok
func (h *APICertificateHandler) MiddlewareValidateCertificateInput(next http.Handler) http.Handler {
	return http.Handler(api.Handler{Handler: func(rw http.ResponseWriter, r *http.Request) *api.APIError {

		// declaring the certInput
		certInput := &APICertificateInput{}

		// parsing cert intput
		err := api.FromJSON(certInput, r.Body)
		if err != nil {
			h.logger.Error("MiddlewareValidateCertificateInput: error deserializing json", "err", err)

			return &api.APIError{
				Err:     err,
				Code:    http.StatusBadRequest,
				Type:    &api.ValidationError{},
				Message: "error deserializing json",
			}
		}

		// validate input
		errs := h.v.Validate(certInput)
		if len(errs) != 0 {
			h.logger.Debug("MiddlewareValidateCertificateInput: invalid input", "certInput", certInput, "errs", errs)
			return &api.APIError{
				Code:    http.StatusBadRequest,
				Type:    &api.ValidationError{},
				Message: fmt.Sprintf("Validation error : %s", strings.Join(errs.Errors(), "\n")),
			}
		}

		// add the input to the context
		ctx := context.WithValue(r.Context(), APICertificateKey{}, *certInput)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
		return nil

	}})
}

// MiddlewareValidateCertificateTagInput validates the tenant in the request and calls next if ok
func (h *APICertificateHandler) MiddlewareValidateCertificateTagInput(next http.Handler) http.Handler {
	return http.Handler(api.Handler{Handler: func(rw http.ResponseWriter, r *http.Request) *api.APIError {

		// declaring the certInput
		certTagInput := &APICertificateTagInput{}

		// parsing cert intput
		err := api.FromJSON(certTagInput, r.Body)
		if err != nil {
			h.logger.Error("MiddlewareValidateCertificateTagInput: error deserializing json", "err", err)

			return &api.APIError{
				Err:     err,
				Code:    http.StatusBadRequest,
				Type:    &api.ValidationError{},
				Message: "error deserializing json",
			}
		}

		// validate input
		errs := h.v.Validate(certTagInput)
		if len(errs) != 0 {
			h.logger.Debug("MiddlewareValidateCertificateTagInput: invalid input", "certTagInput", certTagInput, "errs", errs)
			return &api.APIError{
				Code:    http.StatusBadRequest,
				Type:    &api.ValidationError{},
				Message: fmt.Sprintf("Validation error : %s", strings.Join(errs.Errors(), "\n")),
			}
		}

		// add the input to the context
		ctx := context.WithValue(r.Context(), APICertificateTagKey{}, *certTagInput)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
		return nil

	}})
}
