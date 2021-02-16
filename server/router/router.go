package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"myst/server/api"
	"myst/server/config"
	"myst/server/logger"
	"myst/server/responses"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	cors "github.com/rs/cors/wrapper/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

var (
	log = logger.NewLogger("router", logger.GreenFg)
)

func Init() *gin.Engine {
	// Set gin mode
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	// Discard gin's startup messages
	gin.DefaultWriter = ioutil.Discard
	// Create gin router instance
	r := gin.New()
	// Do not redirect folders to trailing slash
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true
	// Log to stdout and stderr by default
	// Custom PrintRouteFunc
	gin.DebugPrintRouteFunc = PrintRoutes
	// always use recovery middleware
	r.Use(Recovery(RecoveryHandler))
	// custom logging middleware
	r.Use(LoggerMiddleware)
	// metrics
	if config.Debug {
		p := ginprometheus.NewPrometheus("gin")
		p.Use(r)
	}
	// error 404 handling
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			// serve a json 404 error if it's an API call
			data := responses.NewErrorResponse(404, "Route not found")
			c.JSON(404, data)
			c.Abort()
		} else {
			// serve ui and let it handle the error otherwise
			c.File("static/index.html")
			c.Abort()
		}
	})

	// Attach static serve middleware for / and /assets
	r.Use(static.Serve("/", static.LocalFile("static", false)))
	r.Use(static.Serve("/assets", static.LocalFile("assets", false)))

	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:80", "http://localhost:8081"},
		// TODO allow more methods (DELETE?)
		AllowedMethods: []string{http.MethodGet, http.MethodPost},
		// TODO expose ratelimiting headers
		ExposedHeaders: []string{},
		// TODO check if we can disable this on release mode so that no
		// authorization tokens are passed on to the frontend.
		// No harm, but no need either.
		// Required to pass authentication headers on development environment
		AllowCredentials: true,
		Debug:            false, // too verbose, only enable for testing CORS
	}))

	// Pass fizz to route handlers
	api.Init(r)

	return r
}

// RecoveryHandler sends a 500 response if a panic occurs during serving
func RecoveryHandler(c *gin.Context, err interface{}) {
	if strings.HasPrefix(c.Request.URL.Path, "/api/") {
		// If error occurred in an API route print JSON error response
		data := responses.NewErrorResponse(500, nil)
		c.JSON(500, data)
		c.Abort()
	} else {
		// Otherwise just return error 500 status code
		// TODO: Render HTML error 500 template instead
		c.Status(500)
		c.Abort()
	}
}

// Recovery is a panic recovery middleware
func Recovery(f func(c *gin.Context, err interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				goErr := errors.Wrap(err, 3)
				log.Errorf("Panic recovered:\n\n%s%s\n%s", httprequest, goErr.Error(), goErr.Stack())
				f(c, err)
			}
		}()
		c.Next()
	}
}

// PrintRoutes prints all active routes to the console on startup
func PrintRoutes(httpMethod, absolutePath, handlerName string, _ int) {
	if handlerName == "" {
		return
	}
	log.Debugf("%-7s %-50s --> %3s\n", httpMethod, absolutePath, handlerName)
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
		return logger.GreenBg | logger.BlackFg
	case status >= http.StatusMultipleChoices && status < http.StatusBadRequest:
		return logger.WhiteBg | logger.BlackFg
	case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
		return logger.YellowBg | logger.BlackFg
	default:
		return logger.RedBg | logger.BlackFg
	}
}

func MethodColor(method string) logger.Color {
	switch method {
	case http.MethodGet:
		return logger.GreenBg | logger.BlackFg
	case http.MethodPost:
		return logger.YellowBg | logger.BlackFg
	case http.MethodPut:
		return logger.BlueBg | logger.BlackFg
	case http.MethodPatch:
		return logger.CyanBg | logger.BlackFg
	case http.MethodDelete:
		return logger.RedBg | logger.BlackFg
	default:
		return logger.BlackBg
	}
}
