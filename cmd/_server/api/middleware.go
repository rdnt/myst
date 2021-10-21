package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"myst/pkg/regex"
	"myst/pkg/user"
)

var (
	jwtSecretKey = os.Getenv("JWT_SECRET_KEY")

	ErrAuthenticationRequired = fmt.Errorf("authentication required")
	ErrAuthenticationFailed   = fmt.Errorf("authentication failed")
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := TokenAuthentication(c)
		if err != nil {
			Error(c, http.StatusForbidden, err)
			c.Abort()
			return
		}
	}
}

func TokenAuthentication(c *gin.Context) error {
	auth := c.GetHeader("Authorization")
	if auth != "" {
		// Remove the "Bearer" prefix
		parts := strings.Split(auth, "Bearer")
		if len(parts) != 2 {
			return ErrAuthenticationFailed
		}
		auth = strings.TrimSpace(parts[1])
	} else {
		cookie, err := c.Cookie("auth_token")
		if err == nil {
			auth = cookie
		}
	}

	if auth == "" {
		return ErrAuthenticationRequired
	}

	// Validate token format
	match := regex.Match("jwt", auth)
	if !match {
		return ErrAuthenticationFailed
	}

	// Check if authentication token is in the valid format
	token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
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
	})

	if err != nil {
		return ErrAuthenticationFailed
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ErrAuthenticationFailed
	}

	username, ok := claims["usr"].(string)
	if !ok {
		return ErrAuthenticationFailed
	}

	u, err := user.Get("id", username)
	if err == user.ErrNotFound {
		return ErrAuthenticationFailed
	} else if err != nil {
		return err
	}

	err = claims.Valid()
	if err != nil {
		return ErrAuthenticationFailed
	}

	c.Set("user", u.ID)
	return nil
}
