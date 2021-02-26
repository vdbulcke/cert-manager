package tag

import (
	"github.com/hashicorp/go-hclog"
	"github.com/vdbulcke/cert-manager/data"
)

// APITagHandler handler for certificate tag api
type APITagHandler struct {
	logger      hclog.Logger
	v           *data.Validation
	certBackend *data.CertBackend
}

// NewAPITagHandler create a new APITagHandler
func NewAPITagHandler(l hclog.Logger, v *data.Validation, certBackend *data.CertBackend) *APITagHandler {
	return &APITagHandler{
		logger:      l,
		v:           v,
		certBackend: certBackend,
	}
}

// APITagInput input for Tag Creation API
type APITagInput struct {
	Name        string `json:"name" validate:"required,max=50"`
	Description string `json:"description" validate:"max=100" `
}

// APITagInputKey key to find the  APITagInput in request context
type APITagInputKey struct{}

// APITagDescriptionInput  input for Tag Update API
type APITagDescriptionInput struct {
	Description string `json:"description" validate:"required,max=100" `
}

// APITagDescriptionInputKey key to find the  APITagDescriptionInput in request context
type APITagDescriptionInputKey struct{}
