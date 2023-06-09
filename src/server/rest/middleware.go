package rest

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	goerrors "github.com/go-errors/errors"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"

	"myst/pkg/logger"
	"myst/src/server/rest/generated"
)

// noRouteMiddleware is the middleware that processes 404 - not found errors.
// if the route is prefixed with /api/ it will return a json error response,
// otherwise the UI's index.html will be served so that the frontend can render
// the error.
func noRouteMiddleware(c *gin.Context) {
	if strings.HasPrefix(c.Request.URL.Path, "/api/") {
		Error(c, http.StatusNotFound, "not-found", http.StatusText(http.StatusNotFound))
	} else {
		// serve ui and let it handle the error otherwise
		// TODO @rdnt @@@: serve from in-memory filesystem
		c.File("static/index.html")
	}

	c.Abort()
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

// printRoutes prints all active routes to the console on startup
func printRoutes(httpMethod, absolutePath, handlerName string, _ int) {
	if handlerName == "" {
		return
	}

	log.Printf("%-7s %s\n\t -> %s\n", httpMethod, absolutePath, handlerName)
}

func loggerMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
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
		logger.Colorize(fmt.Sprintf(" %d ", status), statusColor(status)),
		latency,
		c.ClientIP(),
		logger.Colorize(fmt.Sprintf(" %s ", method), methodColor(method)),
		path,
		c.Errors.ByType(gin.ErrorTypePrivate).String(),
	)
}

func statusColor(status int) logger.Color {
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

func methodColor(method string) logger.Color {
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
		return errors.New("authentication required")
	}

	parts := strings.Split(auth, "Bearer")
	if len(parts) != 2 {
		return errors.New("authentication failed")
	}
	auth = strings.TrimSpace(parts[1])

	if auth == "" {
		return errors.New("authentication required")
	}

	// Check if the authentication token is valid
	token, err := jwt.Parse(
		auth, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.Errorf("unexpected token signing method: %v", token.Header["alg"])
			}

			return s.jwtSigningKey, nil
		},
	)
	if err != nil {
		return errors.Wrap(err, "invalid jwt format")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid claims type")
	}

	err = claims.Valid()
	if err != nil {
		return errors.Wrap(err, "invalid claims")
	}

	userId, ok := claims["usr"].(string)
	if !ok {
		return errors.New("invalid user claim")
	}

	u, err := s.app.User(userId)
	if err != nil {
		return errors.Wrap(err, "user not found")
	}

	c.Set("userId", u.Id)

	return nil
}
