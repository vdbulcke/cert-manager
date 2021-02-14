package api

import "fmt"

// APIError an API error struct for error handling
type APIError struct {
	Err     error
	Message string
	Code    int
	Type    APIErrorType
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s - %s", e.Type, e.Message)
}

// APIErrorType generic api error type
type APIErrorType interface{}

// InternalServerError default error type
type InternalServerError struct {
	APIErrorType
}

// ObjectNotFoundError when object not found
type ObjectNotFoundError struct {
	APIErrorType
}

// ValidationError when object validation failed
type ValidationError struct {
	APIErrorType
}

// ObjectAlreadyExistError when object already exist
type ObjectAlreadyExistError struct {
	APIErrorType
}

// GenericAPIError is a generic error message returned by a server
type GenericAPIError struct {
	Message string `json:"message"`
	Error   string `json:"error"` // the error type
}
