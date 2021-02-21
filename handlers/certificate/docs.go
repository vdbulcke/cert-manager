// Package certificate  of Certificate API
//
// Documentation for Tenant API
//
//	Schemes: http
//	BasePath: /api/beta2/
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package certificate

import "github.com/vdbulcke/cert-manager/handlers/api"

//
// NOTE: Types defined here are purely for documentation purposes
// these types are not used by any of the handers

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body api.GenericAPIError
}
