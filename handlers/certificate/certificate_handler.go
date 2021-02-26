package certificate

import (
	"github.com/hashicorp/go-hclog"
	"github.com/vdbulcke/cert-manager/data"
)

// APICertificateHandler handler for certificate API
type APICertificateHandler struct {
	logger      hclog.Logger
	v           *data.Validation
	certBackend *data.CertBackend
}

// NewAPICertificateHandler create a new APICertificateHandler
func NewAPICertificateHandler(l hclog.Logger, v *data.Validation, certBackend *data.CertBackend) *APICertificateHandler {
	return &APICertificateHandler{
		logger:      l,
		v:           v,
		certBackend: certBackend,
	}
}

// APICertificateInput input struct for certificate API
type APICertificateInput struct {
	Pem  string   `json:"pem" validate:"required"`
	Tags []string `json:"tags" `
}

// APICertificateTagInput input struct for certificate API
type APICertificateTagInput struct {
	Tags []string `json:"tags" validate:"required"`
}

// APICertificateKey place holder for storing APICertificateInput in request context
type APICertificateKey struct{}

// APICertificateTagKey place holder for storing APICertificateTagInput in request context
type APICertificateTagKey struct{}
