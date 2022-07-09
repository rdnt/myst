package rest

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/golang-jwt/jwt"

	"myst/pkg/logger"
	"myst/pkg/regex"
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
func Recovery(c *gin.Context) {
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

func (s *Server) Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := s.TokenAuthentication(c)
		if err != nil {
			c.JSON(403, gin.H{"error": "Forbidden"})
			c.Abort()

			return
		}
	}
}

func (s *Server) TokenAuthentication(c *gin.Context) error {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return fmt.Errorf("authentication required")
	}

	parts := strings.Split(auth, "Bearer")
	if len(parts) != 2 {
		return fmt.Errorf("authentication failed")
	}
	auth = strings.TrimSpace(parts[1])

	if auth == "" {
		return fmt.Errorf("authentication required")
	}

	// Validate token format
	match := regex.Match("jwt", auth)
	if !match {
		return fmt.Errorf("authentication failed")
	}

	// Check if authentication token is in the valid format
	token, err := jwt.Parse(
		auth, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected token signing method: %v", token.Header["alg"])
			}
			key, err := base64.StdEncoding.DecodeString(jwtSecretKey)
			if err != nil {
				return nil, fmt.Errorf("jwt secret decode failed")
			}
			// return the secret when token is valid format
			return key, nil
		},
	)

	if err != nil {
		return fmt.Errorf("authentication failed")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("authentication failed")
	}

	err = claims.Valid()
	if err != nil {
		return fmt.Errorf("authentication failed")
	}

	userId, ok := claims["usr"].(string)
	if !ok {
		return fmt.Errorf("authentication failed")
	}

	u, err := s.app.User(userId)
	if err != nil {
		return errors.New("user not found")
	}

	if u.Username != "rdnt" && u.Username != "abcd" {
		return fmt.Errorf("authentication failed")
	}

	c.Set("userId", u.Id)

	return nil
}

func CurrentUser(c *gin.Context) string {
	// GetCurrentUserID returns the username of the currently logged-in user
	userId, ok := c.Get("userId")
	if !ok {
		panic("userId not set for authenticated route")
	}

	return userId.(string)
}
