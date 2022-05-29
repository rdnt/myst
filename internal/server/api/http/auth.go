package http

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"myst/internal/server/api/http/generated"
)

func (api *API) Register(c *gin.Context) {
	var req generated.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	u, err := api.app.CreateUser(req.Username, req.Password)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, ToJSONUser(u))
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

	token, err := api.loginUser(params.Username)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, generated.AuthorizationResponse{
		Token:  token,
		UserId: params.Username, // TODO: query user on login and return proper UserID (uuid)
	})
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
