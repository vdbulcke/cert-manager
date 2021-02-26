package tag

import (
	"fmt"
	"net/http"

	"github.com/vdbulcke/cert-manager/data"
	"github.com/vdbulcke/cert-manager/handlers/api"
)

// swagger:route GET /tag/GetTagByID/{id} Tag GetTagByID
// Return tag from the database
// responses:
//	200: tagResponse
//	404: errorResponse
//  400: errorResponse
//  500: errorResponse

// GetTagByID handles GET requests
func (h *APITagHandler) GetTagByID(rw http.ResponseWriter, r *http.Request) *api.APIError {

	uuid, err := api.GetIDFromRequest(r)
	if err != nil {
		return &api.APIError{
			Code:    http.StatusBadRequest,
			Type:    &api.ValidationError{},
			Message: fmt.Sprintf("Cannot find id from request %s", r.URL.String()),
		}
	}

	// lookup this tag
	tag, err := h.certBackend.GetTagByID(uuid)
	if err != nil {
		if _, ok := err.(*data.DBObjectNotFound); ok {
			h.logger.Debug("GetTagByID: object not found", "uuid", uuid)

			return &api.APIError{
				Err:     err,
				Code:    http.StatusNotFound,
				Type:    &api.ObjectNotFoundError{},
				Message: err.Error(),
			}

		}

		// default
		h.logger.Debug("GetTagByID: unexpected error searching for tag", "uuid", uuid, "err", err)

		return &api.APIError{
			Err:     err,
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: fmt.Sprintf("Error searching for tag id=%s", uuid.String()),
		}

	}

	h.logger.Debug("GetTagByID: Found tag", "tag", tag)

	// Write Status code
	rw.WriteHeader(http.StatusOK)
	err = api.ToJSON(tag, rw)
	if err != nil {
		// we should never be here but log the error just incase
		h.logger.Error("GetTagByID: Error Serializing JSON", "tag", tag, "err", err)
		return &api.APIError{
			Err:     err,
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: "error formating json",
		}
	}

	return nil
}

// swagger:route GET /tag/GetTagByName/{name} Tag GetTagByName
// Return tag from the database
// responses:
//	200: tagResponse
//	404: errorResponse
//  400: errorResponse
//  500: errorResponse

// GetTagByName handles GET requests
func (h *APITagHandler) GetTagByName(rw http.ResponseWriter, r *http.Request) *api.APIError {

	tagName, err := getTagNameFromRequest(r)
	if err != nil {
		return &api.APIError{
			Code:    http.StatusBadRequest,
			Type:    &api.ValidationError{},
			Message: fmt.Sprintf("Cannot find tag Name from request %s", r.URL.String()),
		}
	}

	// lookup this tag
	tag, err := h.certBackend.GetTagByName(tagName)
	if err != nil {
		if _, ok := err.(*data.DBObjectNotFound); ok {
			h.logger.Debug("GetTagByName: object not found", "tagName", tagName)

			return &api.APIError{
				Err:     err,
				Code:    http.StatusNotFound,
				Type:    &api.ObjectNotFoundError{},
				Message: err.Error(),
			}

		}

		// default
		h.logger.Debug("GetTagByName: unexpected error searching for tag", "tagName", tagName, "err", err)

		return &api.APIError{
			Err:     err,
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: fmt.Sprintf("Error searching for tag name=%s", tagName),
		}

	}

	h.logger.Debug("GetTagByName: Found tag", "tag", tag)

	// Write Status code
	rw.WriteHeader(http.StatusOK)
	err = api.ToJSON(tag, rw)
	if err != nil {
		// we should never be here but log the error just incase
		h.logger.Error("GetTagByName: Error Serializing JSON", "tag", tag, "err", err)
		return &api.APIError{
			Err:     err,
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: "error formating json",
		}
	}

	return nil
}

// swagger:route GET /tag/ListTags Tag ListTags
// Return tag from the database
// responses:
//	200: tagListResponse
//	404: errorResponse
//  400: errorResponse
//  500: errorResponse

// ListTags handles GET requests
func (h *APITagHandler) ListTags(rw http.ResponseWriter, r *http.Request) *api.APIError {

	// lookup this tag
	tags, err := h.certBackend.ListTags()
	if err != nil {
		if _, ok := err.(*data.DBObjectNotFound); ok {
			h.logger.Debug("ListTags: object not found", "tags", tags)

			return &api.APIError{
				Err:     err,
				Code:    http.StatusNotFound,
				Type:    &api.ObjectNotFoundError{},
				Message: err.Error(),
			}

		}

		// default
		h.logger.Debug("ListTags: unexpected error searching for tag", "err", err)

		return &api.APIError{
			Err:     err,
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: fmt.Sprintf("Error searching for tags"),
		}

	}

	h.logger.Debug("ListTags: Found tag", "tags", tags)

	// Write Status code
	rw.WriteHeader(http.StatusOK)
	err = api.ToJSON(tags, rw)
	if err != nil {
		// we should never be here but log the error just incase
		h.logger.Error("GetTagByName: Error Serializing JSON", "tag", tags, "err", err)
		return &api.APIError{
			Err:     err,
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: "error formating json",
		}
	}

	return nil
}
