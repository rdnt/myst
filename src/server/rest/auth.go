package rest

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"myst/src/server/rest/generated"
)

func (s *Server) Register(c *gin.Context) {
	var req generated.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	u, err := s.app.CreateUser(req.Username, req.Password, req.PublicKey)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	token, err := s.loginUser(u.Id)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, generated.AuthorizationResponse{
		User: generated.User{
			Id:        u.Id,
			Username:  u.Username,
			PublicKey: u.PublicKey,
		},
		Token: token,
	})
}

func (s *Server) Login(c *gin.Context) {
	var params generated.LoginRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	u, err := s.app.UserByUsername(params.Username)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// // TODO: proper password hash check
	// if !((params.Username == "rdnt" && params.Password == "1234") || (params.Username == "abcd" && params.Password == "5678")) {
	// 	panic("invalid username or password")
	// }

	err = s.app.AuthorizeUser(params.Username, params.Password)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	token, err := s.loginUser(u.Id)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, generated.AuthorizationResponse{
		User: generated.User{
			Id:       u.Id,
			Username: u.Username,
		},
		Token: token,
	})
}

func (s *Server) loginUser(userId string) (string, error) {
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
			"usr": userId,
		},
	)

	key, err := base64.StdEncoding.DecodeString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return token.SignedString(key)
}

type UserByUsernameRequest struct {
	Username *string `form:"username"`
}

func (s *Server) User(c *gin.Context) {
	var req UserByUsernameRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if req.Username == nil {
		c.Status(http.StatusBadRequest)
		return
	}

	u, err := s.app.UserByUsername(*req.Username)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, UserToJSON(u))
}
