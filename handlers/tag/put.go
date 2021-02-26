package tag

import (
	"fmt"
	"net/http"

	"github.com/vdbulcke/cert-manager/data"
	"github.com/vdbulcke/cert-manager/handlers/api"
)

// swagger:route PUT /tag/UpdateTagDescription/{id} Tag CreateTag
// Return newly Updated Tag from the database
// responses:
//	202: tagResponse
//	404: errorResponse
//  400: errorResponse
//  500: errorResponse

// UpdateTagDescription handles POST requests
func (h *APITagHandler) UpdateTagDescription(rw http.ResponseWriter, r *http.Request) *api.APIError {

	uuid, err := api.GetIDFromRequest(r)
	if err != nil {
		return &api.APIError{
			Code:    http.StatusBadRequest,
			Type:    &api.ValidationError{},
			Message: fmt.Sprintf("Cannot find id from request %s", r.URL.String()),
		}
	}

	// look up tag input (set by the middleware) in request context
	tagDescriptionInput := r.Context().Value(APITagDescriptionInputKey{}).(APITagDescriptionInput)

	// checking tagDescriptionInput contains a Name
	if tagDescriptionInput.Description == "" {
		return &api.APIError{
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: "Error parsing input",
		}
	}

	// updating tag description
	tag, err := h.certBackend.SetTagDescriptionByID(uuid, tagDescriptionInput.Description)
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

	h.logger.Debug("UpdateTagDescription: Found tag", "tag", tag)

	// Write Status code
	rw.WriteHeader(http.StatusAccepted)
	err = api.ToJSON(tag, rw)
	if err != nil {
		// we should never be here but log the error just incase
		h.logger.Error("UpdateTagDescription: Error Serializing JSON", "tag", tag, "err", err)
		return &api.APIError{
			Err:     err,
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: "error formating json",
		}
	}

	return nil
}
