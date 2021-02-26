package tag

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/vdbulcke/cert-manager/handlers/api"
)

// MiddlewareValidateTagInput validates the tag in the request and calls next if ok
func (h *APITagHandler) MiddlewareValidateTagInput(next http.Handler) http.Handler {
	return http.Handler(api.Handler{Handler: func(rw http.ResponseWriter, r *http.Request) *api.APIError {

		// declaring the tagInput
		tagInput := &APITagInput{}

		// parsing tag intput
		err := api.FromJSON(tagInput, r.Body)
		if err != nil {
			h.logger.Error("MiddlewareValidateTagInput: error deserializing json", "err", err)

			return &api.APIError{
				Err:     err,
				Code:    http.StatusBadRequest,
				Type:    &api.ValidationError{},
				Message: "error deserializing json",
			}
		}

		// validate input
		errs := h.v.Validate(tagInput)
		if len(errs) != 0 {
			h.logger.Debug("MiddlewareValidateTagInput: invalid input", "tagInput", tagInput, "errs", errs)
			return &api.APIError{
				Code:    http.StatusBadRequest,
				Type:    &api.ValidationError{},
				Message: fmt.Sprintf("Validation error : %s", strings.Join(errs.Errors(), "\n")),
			}
		}

		// add the input to the context
		ctx := context.WithValue(r.Context(), APITagInputKey{}, *tagInput)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
		return nil

	}})
}

// MiddlewareValidateTagDescriptionInput validates the tag in the request and calls next if ok
func (h *APITagHandler) MiddlewareValidateTagDescriptionInput(next http.Handler) http.Handler {
	return http.Handler(api.Handler{Handler: func(rw http.ResponseWriter, r *http.Request) *api.APIError {

		// declaring the tag description Input
		tagDescriptionInput := &APITagDescriptionInput{}

		// parsing tag intput
		err := api.FromJSON(tagDescriptionInput, r.Body)
		if err != nil {
			h.logger.Error("MiddlewareValidateTagInput: error deserializing json", "err", err)

			return &api.APIError{
				Err:     err,
				Code:    http.StatusBadRequest,
				Type:    &api.ValidationError{},
				Message: "error deserializing json",
			}
		}

		// validate input
		errs := h.v.Validate(tagDescriptionInput)
		if len(errs) != 0 {
			h.logger.Debug("MiddlewareValidateTagInput: invalid input", "tagDescriptionInput", tagDescriptionInput, "errs", errs)
			return &api.APIError{
				Code:    http.StatusBadRequest,
				Type:    &api.ValidationError{},
				Message: fmt.Sprintf("Validation error : %s", strings.Join(errs.Errors(), "\n")),
			}
		}

		// add the input to the context
		ctx := context.WithValue(r.Context(), APITagDescriptionInputKey{}, *tagDescriptionInput)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
		return nil

	}})
}
