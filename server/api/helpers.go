package api

import (
	"github.com/gin-gonic/gin"
	"myst/server/responses"
)

// Success prints a successful response (with optional data field) and cancels
// any further handlers in the chain
func Success(c *gin.Context, data interface{}) {
	// Create response message
	msg := responses.NewSuccessResponse(data)
	// Echo it
	c.JSON(200, msg)
	// Stop the chain of handlers
	c.Abort()
}

// Error prints a custom error message and cancels any further
// handlers in the chain
func Error(c *gin.Context, code int, msg interface{}) {
	// Create custom error response
	resp := responses.NewErrorResponse(code, msg)
	// Echo it
	c.JSON(code, resp)
	// Stop the chain of handlers
	c.Abort()
}
