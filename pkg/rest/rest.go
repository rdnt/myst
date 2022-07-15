package rest

import (
	"net/http"
)

// SuccessResponse is a success response with embedded data (if passed)
type SuccessResponse struct {
	Status     string                 `json:"status"`
	Data       interface{}            `json:"data,omitempty"`
	Links      map[string]interface{} `json:"links,omitempty"`
	Pagination *Pagination            `json:"pagination,omitempty"`
}

type Link struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Pagination adds optional pagination metadata to a response
type Pagination struct {
	Total   int `json:"total"`
	PerPage int `json:"per-page"`
}

// ErrorResponse is an error response. Adds a default rest status message if
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
