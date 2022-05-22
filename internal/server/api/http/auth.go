package http

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"myst/internal/server/api/http/generated"
	"myst/internal/server/core/domain/user"
)

func (api *API) Register(c *gin.Context) {
	var req generated.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		panic(err)
	}

	u, err := api.app.CreateUser(req.Username, req.Password)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, ToJSONUser(u))
}

func ToJSONUser(u user.User) generated.User {
	return generated.User{
		Id:       u.Id,
		Username: u.Username,
	}
}

func (api *API) Login(c *gin.Context) {
	var params generated.LoginRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		panic(err)
	}

	if !((params.Username == "rdnt" && params.Password == "1234") || (params.Username == "abcd" && params.Password == "5678")) {
		panic("invalid username or password")
	}

	key, err := api.loginUser(params.Username)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, generated.AuthToken(key))
}

func (api *API) loginUser(username string) (string, error) {
	now := time.Now()

	exp := now.Unix() + int64(jwtCookieLifetime)
	iat := now.Unix()
	nbf := now.Unix()

	// if password matches hash
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": exp,
			"iat": iat,
			"nbf": nbf,
			"usr": username,
		},
	)
	key, err := base64.StdEncoding.DecodeString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return token.SignedString(key)
}
