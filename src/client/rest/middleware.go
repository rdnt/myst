package rest

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"

	"myst/pkg/logger"
	"myst/src/server/rest/generated"
)

func NoRoute(c *gin.Context) {
	if strings.HasPrefix(c.Request.URL.Path, "/api/") {
		// serve a json 404 error if it's an Server call
		resp := generated.Error{
			Code:    "not-found",
			Message: http.StatusText(http.StatusNotFound),
		}

		c.JSON(http.StatusNotFound, resp)
	} else {
		// serve ui and let it handle the error otherwise
		c.File("static/index.html")
	}

	c.Abort()
}

// Recovery is a panic recovery middleware
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				goErr := errors.Wrap(err, 0)
				log.Errorf("Panic recovered:\n\n%s%s\n%s", httprequest, goErr.Error(), goErr.Stack())
				recoveryHandler(c, err)
			}
		}()

		c.Next()
	}
}

// recoveryHandler sends a 500 response if a panic occurs during serving
func recoveryHandler(c *gin.Context, err interface{}) {
	if strings.HasPrefix(c.Request.URL.Path, "/api/") {
		// If error occurred in an Server route print JSON error response
		log.Error(err)

		resp := generated.Error{
			Code:    "internal-error",
			Message: http.StatusText(http.StatusInternalServerError),
		}

		c.JSON(http.StatusInternalServerError, resp)
	} else {
		// Otherwise just return error 500 status code
		// TODO: Render HTML error 500 template instead
		c.Status(http.StatusInternalServerError)
	}

	c.Abort()
}

// PrintRoutes prints all active routes to the console on startup
func PrintRoutes(httpMethod, absolutePath, handlerName string, _ int) {
	if handlerName == "" {
		return
	}
	log.Printf("%-7s %-50s --> %3s\n", httpMethod, absolutePath, handlerName)
}

func LoggerMiddleware(c *gin.Context) {
	// Start timer
	start := time.Now()
	// Process request
	c.Next()
	// calculate latency
	latency := time.Since(start)

	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	if query != "" {
		path = path + "?" + query
	}

	status := c.Writer.Status()
	method := c.Request.Method

	log.Printf(
		"%5s  %13v  %15s  %-21s  %s\n%s",
		logger.Colorize(fmt.Sprintf(" %d ", status), StatusColor(status)),
		latency,
		c.ClientIP(),
		logger.Colorize(fmt.Sprintf(" %s ", method), MethodColor(method)),
		path,
		c.Errors.ByType(gin.ErrorTypePrivate).String(),
	)
}

func StatusColor(status int) logger.Color {
	switch {
	case status >= http.StatusOK && status < http.StatusMultipleChoices:
		return logger.GreenBg | logger.Black
	case status >= http.StatusMultipleChoices && status < http.StatusBadRequest:
		return logger.WhiteBg | logger.Black
	case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
		return logger.YellowBg | logger.Black
	default:
		return logger.RedBg | logger.Black
	}
}

func MethodColor(method string) logger.Color {
	switch method {
	case http.MethodGet:
		return logger.GreenBg | logger.Black
	case http.MethodPost:
		return logger.BlueBg | logger.Black
	case http.MethodPut:
		return logger.CyanBg | logger.Black
	case http.MethodPatch:
		return logger.YellowBg | logger.Black
	case http.MethodDelete:
		return logger.RedBg | logger.Black
	default:
		return logger.MagentaBg | logger.Black
	}
}
