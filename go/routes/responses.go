package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sht/shtdev/go/responses"
	"strconv"
)

// Success prints a successful response message and cancels any further
// handlers in the chain
func Success(c *gin.Context, data interface{}) {
	// Create response message
	msg := responses.NewResponse(data)
	// Echo it
	c.JSON(200, msg)
	// Stop the chain of handlers
	c.Abort()
}

// HTTPError prints an HTTP error message and cancels any further
// handlers in the chain
func HTTPError(c *gin.Context, id int) {
	// Create HTTO error message
	msg := responses.NewHTTPError(id)
	// Echo it
	c.JSON(id, msg)
	// Stop the chain of handlers
	c.Abort()
}

// Error prints a custom error message and cancels any further
// handlers in the chasin
func Error(c *gin.Context, code string, message interface{}) {
	// Create custom error message
	msg := responses.NewError(code, message)
	// Echo it
	c.JSON(200, msg)
	// Stop the chain of handlers
	c.Abort()
}

// Invalidate invalidates a specific data field and sets a reason
func Invalidate(c *gin.Context, name, desc string) {
	// Format fail message
	err := responses.AddError(name, desc)
	errors, exists := c.Get("errors")
	if !exists {

	}
	errors = append(errors.([]responses.Error), err)
	c.Set("errors", errors)
}

// DataValid reports invalid data with reasoning and returns whether the
// request should be dropped or not
func DataValid(c *gin.Context) bool {
	// Format fail message
	errors, exists := c.Get("errors")
	if !exists {

	}
	errarr := errors.([]responses.Error)
	if len(errarr) > 0 {
		msg := responses.NewFail(errarr)
		// Send response
		c.JSON(200, msg)
		// Stop the chain of middleware
		c.Abort()
		return false
	}
	return true
}

// WithPagination sets the correct HTTP headers based on the pagination provided
func WithPagination(c *gin.Context, count int) {
	c.Header("Total-Count", strconv.Itoa(count))
}

// ValidatePagination returns sanitized limit & offset arguments
func ValidatePagination(perPage, page int) (*Pagination, error) {
	// Calculate limit
	var limit, offset int
	// Invalidate negative pagination
	if perPage < 0 || page < 0 {
		return nil, fmt.Errorf("invalid pagination: per_page = %d, page = %d", perPage, page)
	} else if perPage == 0 && page > 0 {
		return nil, fmt.Errorf("invalid pagination: per_page = %d,", perPage)
	} else if perPage > 0 && page == 0 {
		return nil, fmt.Errorf("invalid pagination: page = %d,", page)
	}
	// both are either 0 or >0
	if perPage == 0 && page == 0 {
		limit = int(^uint(0) >> 1) // maxInt
		offset = 0
	} else if perPage > 0 && page > 0 {
		limit = perPage
		offset = (page - 1) * perPage
	}
	// Return pagination struct
	return &Pagination{
		Limit:  limit,
		Offset: offset,
	}, nil
}

// Pagination asd
type Pagination struct {
	Limit  int
	Offset int
}
