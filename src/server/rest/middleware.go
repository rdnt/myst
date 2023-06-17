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
	"myst/pkg/router"
)

// noRouteMiddleware is the middleware that processes http 404 errors.
func noRouteMiddleware(c *gin.Context) {
	Error(c, http.StatusNotFound)
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
	log.Error(err)
	Error(c, http.StatusInternalServerError)

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
