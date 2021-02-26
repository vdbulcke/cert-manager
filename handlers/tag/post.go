package tag

import (
	"net/http"

	"github.com/vdbulcke/cert-manager/data"
	"github.com/vdbulcke/cert-manager/handlers/api"
)

// swagger:route POST /tag/CreateTag Tag CreateTag
// Return newly Created Tag from the database
// responses:
//	200: tagResponse
//  409: tagResponse
//	404: errorResponse
//  400: errorResponse
//  500: errorResponse

// CreateTag handles POST requests
func (h *APITagHandler) CreateTag(rw http.ResponseWriter, r *http.Request) *api.APIError {

	// look up tag input (set by the middleware) in request context
	tagInput := r.Context().Value(APITagInputKey{}).(APITagInput)

	// checking tagInput contains a Name
	if tagInput.Name == "" {
		return &api.APIError{
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: "Error parsing input",
		}
	}

	h.logger.Debug("CreateTag: tagInput", "tagInput", tagInput)
	// setting default status code
	statusCode := http.StatusCreated
	var tag *data.Tag
	var err error

	// checking if the input contains description
	if tagInput.Description == "" {
		tag, err = h.certBackend.CreateTag(tagInput.Name)
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
		tag, err = h.certBackend.CreateTagWithDescription(tagInput.Name, tagInput.Description)
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

	h.logger.Debug("CreateTag: Found tag", "tag", tag)

	// Write Status code
	rw.WriteHeader(statusCode)
	err = api.ToJSON(tag, rw)
	if err != nil {
		// we should never be here but log the error just incase
		h.logger.Error("CreateTag: Error Serializing JSON", "tag", tag, "err", err)
		return &api.APIError{
			Err:     err,
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: "error formating json",
		}
	}

	return nil
}
