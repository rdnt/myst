package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"myst/pkg/rest"
)

var (
	ErrUnauthenticatedRoute = fmt.Errorf("unauthenticated route")
)

func Paginate(c *gin.Context, total int) {
	c.Set("pagination", true)
	c.Set("total", total)
}

func WithLinks(c *gin.Context, prev, next interface{}) {
	c.Set("links", true)
	c.Set("prev", prev)
	c.Set("next", next)
}

// Success prints a successful response (with optional data field) and cancels
// any further handlers in the chain
func Success(c *gin.Context, data interface{}) {
	msg := rest.NewSuccessResponse(data)
	// attach pagination, if applicable
	if c.GetBool("pagination") {
		msg.Pagination = &rest.Pagination{
			Total:   c.GetInt("total"),
			PerPage: 10, // get from config
		}
	}
	if c.GetBool("links") {
		prev, ok1 := c.Get("prev")
		next, ok2 := c.Get("next")
		if ok1 && ok2 {
			msg.Links = map[string]interface{}{
				"prev": prev,
				"next": next,
			}
		}
	}
	// Echo it
	c.JSON(200, msg)
	// Stop the chain of handlers
	c.Abort()
}

// Error prints a custom error message and cancels any further
// handlers in the chain
func Error(c *gin.Context, code int, msg interface{}) {
	// Create custom error response
	resp := rest.NewErrorResponse(code, msg)
	// Echo it
	c.JSON(code, resp)
	// Stop the chain of handlers
	c.Abort()
}

// GetCurrentUserID returns the username of the currently logged-in user
func GetCurrentUserID(c *gin.Context) string {
	uid, ok := c.Get("user")
	if !ok {
		panic(ErrUnauthenticatedRoute)
	}
	return uid.(string)
}
