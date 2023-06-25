package rest

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	goerrors "github.com/go-errors/errors"
	"github.com/pkg/errors"

	"myst/pkg/logger"
	"myst/pkg/router"
)

// noRouteMiddleware is the middleware that processes http 404 errors.
// If the request is an API call, then a JSON 404 error is returned.
// In any other case, if a file exists on the embedded UI filesystem, it is
// served, otherwise the index.html file is served so that the UI can
// render an error page.
func noRouteMiddleware(urlPrefix string, fs static.ServeFileSystem) gin.HandlerFunc {
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}

	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			// serve a json 404 error if it's an api call
			Error(c, http.StatusNotFound)
			c.Abort()
		} else if fs.Exists(urlPrefix, c.Request.URL.Path) {
			// serve ui and let it handle the error otherwise
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		} else {
			// serve the / path
			c.Request.URL.Path = "/"
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

// recoveryMiddleware is a panic recovery middleware
func recoveryMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			httprequest, _ := httputil.DumpRequest(c.Request, false)
			goErr := goerrors.Wrap(err, 0)
			log.Errorf("Panic recovered:\n\n%s%s\n%s", httprequest, goErr.Error(), goErr.Stack())
			recoveryHandler(c, err)
		}
	}()

	c.Next()
}

// recoveryHandler sends a 500 response if a panic occurs during serving
func recoveryHandler(c *gin.Context, err interface{}) {
	if strings.HasPrefix(c.Request.URL.Path, "/api/") {
		// If error occurred in an API route print JSON error response
		log.Error(err)
		Error(c, http.StatusInternalServerError)
	} else {
		// Otherwise just return error 500 status code
		// TODO: Render HTML error 500 template instead
		c.Status(http.StatusInternalServerError)
	}

	c.Abort()
}

// printRoutes prints all active routes to the console on startup
func printRoutes(httpMethod, absolutePath, handlerName string, _ int) {
	if handlerName == "" {
		return
	}
	log.Printf("%-7s %s\n\t -> %s\n", httpMethod, absolutePath, handlerName)
}

func loggerMiddleware(c *gin.Context) {
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
		logger.Colorize(fmt.Sprintf(" %d ", status), router.StatusColor(status)),
		latency,
		c.ClientIP(),
		logger.Colorize(fmt.Sprintf(" %s ", method), router.MethodColor(method)),
		path,
		c.Errors.ByType(gin.ErrorTypePrivate).String(),
	)
}

func (s *Server) authenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := s.authenticateWithToken(c)
		if err != nil {
			log.Error(err)
			Error(c, http.StatusForbidden)
			c.Abort()
		}
	}
}

func (s *Server) authenticateWithToken(c *gin.Context) error {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return nil
	}

	parts := strings.Split(auth, "Bearer")
	if len(parts) != 2 {
		return errors.New("authentication failed")
	}
	auth = strings.TrimSpace(parts[1])

	if auth == "" {
		return errors.New("authentication required")
	}

	sid, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return errors.Wrap(err, "invalid token")
	}

	c.Set("sessionId", sid)

	return nil
}
