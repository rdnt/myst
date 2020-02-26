package router

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-errors/errors"
	"github.com/sht/myst/go/config"
	"github.com/sht/myst/go/logger"
	"github.com/sht/myst/go/regex"
	"github.com/sht/myst/go/responses"
	"github.com/sht/myst/go/routes"
	"github.com/unrolled/secure"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"net/http/httputil"
	// "reflect"
	"strings"
)

var validateRegex validator.Func = func(fl validator.FieldLevel) bool {
	// name := fl.FieldName()
	rex := fl.Param()
	val := fl.Field().String()
	match := regex.Match(rex, val)
	if !match {
		// routes.Invalidate(
		// 	c,
		// 	name,
		// 	"Validation failed for field "+name,
		// )
		return false
	}
	return true
}

func Init() *gin.Engine {
	// Disable console color by default
	gin.DisableConsoleColor()
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
	if config.Debug {
		gin.DefaultWriter = &logger.StdoutLogger
		gin.DefaultErrorWriter = &logger.StderrLogger
	} else {
		gin.DefaultWriter = &logger.AccessLogger
		gin.DefaultErrorWriter = &logger.ErrorLogger
	}
	// Custom PrintRouteFunc
	gin.DebugPrintRouteFunc = logger.PrintRoutes
	// Always use recovery middleware
	r.Use(Recovery(RecoveryHandler))
	// Initialize custom console & file logging middleware
	consoleLogger := gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: logger.ConsoleFormatter,
		Output:    &logger.StdoutLogger,
		SkipPaths: []string{},
	})
	fileLogger := gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: logger.LogFormatter,
		Output:    &logger.AccessLogger,
		SkipPaths: []string{},
	})
	// Enable logging to console & logfiles if debugging, else only log to logfiles
	if config.Debug {
		r.Use(consoleLogger)
		r.Use(fileLogger)
	} else {
		r.Use(fileLogger)
	}

	r.Use(HTTPSRedirect())

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		// gin.DisableBindValidation()
		v.RegisterValidation("regex", validateRegex)
	}

	// Attach client error collecting middleware
	r.Use(Pusher())
	r.Use(ClientErrors())
	// 404 handling
	r.NoRoute(func(c *gin.Context) {
		data := responses.NewHTTPError(404)
		c.JSON(404, data)
		c.Abort()
	})

	// Initialize the rest of the routes
	routes.Init(r)
	return r
}

func HTTPSRedirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		sec := secure.New(secure.Options{
			SSLRedirect: true,
			// Don't redirect when in debug mode
			IsDevelopment: config.Debug,
		})
		err := sec.Process(c.Writer, c.Request)
		// If there was an error abort handlers and return
		if err != nil {
			c.Abort()
			return
		}
		// Prevents errors due to headers already being sent
		if status := c.Writer.Status(); status > 300 && status < 399 {
			c.Abort()
		}
	}
}

func RecoveryHandler(c *gin.Context, err interface{}) {
	if strings.HasPrefix(c.Request.URL.Path, "/api/") {
		// If error occured in an API route print JSON error response
		data := responses.NewHTTPError(500)
		c.JSON(500, data)
		c.Abort()
	} else {
		// Otherwise just return error 500 status code
		// TODO: Render HTML error 500 template instead
		c.Status(500)
		c.Abort()
	}
}

// ClientErrorsMiddleware initializes the validation errors slice and stores
// it in the context
func ClientErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize validation errors slice
		c.Set("errors", []responses.Error{})
		// Pass onto the next handler
		c.Next()
	}
}

func Pusher() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/" {
			PushAssets(c)
		}
	}
}

func PushAssets(c *gin.Context) {
	p := c.Writer.Pusher()
	if p == nil {
		return
	}
	p.Push("/assets/images/logo.svg", nil)
}

func Recovery(f func(c *gin.Context, err interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				goErr := errors.Wrap(err, 3)
				logger.Errorf("RECOVERY", "Panic recovered:\n\n%s%s\n%s", httprequest, goErr.Error(), goErr.Stack())
				f(c, err)
			}
		}()
		c.Next()
	}
}
