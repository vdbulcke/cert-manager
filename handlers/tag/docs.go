// Package tag  of Tag API
//
// Documentation for Tag API
//
//	Schemes: http
//	BasePath: /api/beta2/
//	Version: 0.1.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package tag

import (
	"github.com/vdbulcke/cert-manager/data"
	"github.com/vdbulcke/cert-manager/handlers/api"
)

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body api.GenericAPIError
}

// No content is returned by this API endpoint
// swagger:response noContentResponse
type noContentResponseWrapper struct{}

// A list of tags
// swagger:response tagListResponse
type tagListResponseWrapper struct {
	// All current tags
	// in: body
	Body []data.Tags
}

// Data structure representing a single tag
// swagger:response tagResponse
type tagResponseWrapper struct {
	// a Tag
	// in: body
	Body data.Tag
}

// swagger:parameters CreateTag
type tagCreateParamsWrapper struct {
	// Tag data structure to Create
	// in: body
	// required: true
	Body APITagInput
}

// swagger:parameters TagDescriptionUpdateCertificate
type tagDescriptionUpdateParamsWrapper struct {
	// Tag description to Update
	// in: body
	// required: true
	Body APITagDescriptionInput
}

// swagger:parameters DeleteTagByID GetTagByID
type tagIDParamsWrapper struct {
	// The id of the tag for which the operation relates
	// in: path
	// required: true
	ID string `json:"id"`
}

// swagger:parameters GetTagByName
type tagNameParamsWrapper struct {
	// The id of the tag for which the operation relates
	// in: path
	// required: true
	Name string `json:"name"`
}
