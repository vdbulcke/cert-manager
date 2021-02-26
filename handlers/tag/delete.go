package tag

import (
	"fmt"
	"net/http"

	"github.com/vdbulcke/cert-manager/data"
	"github.com/vdbulcke/cert-manager/handlers/api"
)

// swagger:route DELETE /tag/DeleteTagByID/{id} Tag DeleteTagByID
// Deletes Tags from DB (and all association)
// responses:
//	204: noContentResponse
//	404: errorResponse
//  400: errorResponse
//  500: errorResponse

// DeleteTagByID handles GET requests
func (h *APITagHandler) DeleteTagByID(rw http.ResponseWriter, r *http.Request) *api.APIError {

	uuid, err := api.GetIDFromRequest(r)
	if err != nil {
		return &api.APIError{
			Code:    http.StatusBadRequest,
			Type:    &api.ValidationError{},
			Message: fmt.Sprintf("Cannot find id from request %s", r.URL.String()),
		}
	}

	//  deletes tag
	err = h.certBackend.DeleteTagByID(uuid)
	if err != nil {
		if _, ok := err.(*data.DBObjectNotFound); ok {
			h.logger.Debug("DeleteTagByID: object not found", "uuid", uuid)

			return &api.APIError{
				Err:     err,
				Code:    http.StatusNotFound,
				Type:    &api.ObjectNotFoundError{},
				Message: err.Error(),
			}

		}

		// default
		h.logger.Debug("DeleteTagByID: unexpected error searching for tag", "uuid", uuid, "err", err)

		return &api.APIError{
			Err:     err,
			Code:    http.StatusInternalServerError,
			Type:    &api.InternalServerError{},
			Message: fmt.Sprintf("Error Deleting tag id=%s", uuid.String()),
		}

	}

	// Write Status code
	rw.WriteHeader(http.StatusNoContent)

	return nil
}
