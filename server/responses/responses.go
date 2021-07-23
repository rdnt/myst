package responses

import (
	"net/http"
)

// SuccessResponse is a success response with embedded data (if passed)
type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

// ErrorResponse is an error response. Adds a default http status message if
// none was supplied
type ErrorResponse struct {
	Status  string      `json:"status"`
	Message interface{} `json:"message"`
}

// NewSuccessResponse returns a success response message
func NewSuccessResponse(data interface{}) *SuccessResponse {
	// Format response message
	return &SuccessResponse{
		Status: "success",
		Data:   data,
	}
}

// NewErrorResponse returns an HTTP error message
func NewErrorResponse(code int, msg interface{}) *ErrorResponse {
	if msg == nil {
		msg = http.StatusText(code)
	}
	err, ok := msg.(error)
	if ok {
		msg = err.Error()
	}
	// Format response message
	return &ErrorResponse{
		Status:  "error",
		Message: msg,
	}
}
