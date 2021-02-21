// Package certificate  of Certificate API
//
// Documentation for Tenant API
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
package certificate

import (
	"github.com/vdbulcke/cert-manager/data"
	"github.com/vdbulcke/cert-manager/handlers/api"
)

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

// A list of certificates
// swagger:response certificateListResponse
type certListResponseWrapper struct {
	// All current certificates
	// in: body
	Body []data.Certificates
}

// Data structure representing a single certificate
// swagger:response certificateResponse
type certResponseWrapper struct {
	// a certificate
	// in: body
	Body data.Certificate
}

// No content is returned by this API endpoint
// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// swagger:parameters CreateCertificate
type certificateCreateParamsWrapper struct {
	// Certificate data structure to Create
	// in: body
	// required: true
	Body APICertificateInput
}

// swagger:parameters UpdateCertificateTag DeleteCertificateTagsByID
type certificateUpdateTagsParamsWrapper struct {
	// Tags data structure to Add or Remove from Certificate
	// in: body
	// required: true
	Body APICertificateTagInput
}

// swagger:parameters GetCertificateByID GetCertificateByFingerprint UpdateCertificateTag DeleteCertificateTagsByID DeleteCertificateByID
type certificateIDParamsWrapper struct {
	// The id of the certificate for which the operation relates
	// could be uuid or sha256 fingerprint
	// in: path
	// required: true
	ID string `json:"id"`
}
