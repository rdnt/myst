package rest

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// signedToken creates a jwt token for the given user
func (s *Server) signedToken(userId string) (string, error) {
	now := time.Now()

	exp := now.Unix() + jwtLifetime
	iat := now.Unix()
	nbf := now.Unix()

	// if password matches hash
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": exp,
			"iat": iat,
			"nbf": nbf,
			"usr": userId,
		},
	)

	return token.SignedString(s.jwtSigningKey)
}

func CurrentUser(c *gin.Context) string {
	// GetCurrentUserID returns the username of the currently logged-in user
	userId, ok := c.Get("userId")
	if !ok {
		panic("userId not set")
	}

	return userId.(string)
}
