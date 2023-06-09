package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"myst/src/server/application"
	"myst/src/server/rest/generated"
)

func (s *Server) Register(c *gin.Context) {
	var req generated.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest)
		return
	}

	u, err := s.app.CreateUser(req.Username, req.Password, req.PublicKey)
	switch {
	case errors.Is(err, application.ErrInvalidUsername):
		Error(c, http.StatusBadRequest, "invalid-username")
		return
	case errors.Is(err, application.ErrInvalidPassword):
		Error(c, http.StatusBadRequest, "invalid-password")
		return
	case errors.Is(err, application.ErrUsernameTaken):
		Error(c, http.StatusBadRequest, "username-taken")
		return
	case err != nil:
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	token, err := s.signedToken(u.Id)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
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
		Error(c, http.StatusBadRequest)
		return
	}

	u, err := s.app.AuthenticateUser(params.Username, params.Password)
	if errors.Is(err, application.ErrAuthenticationFailed) {
		Error(c, http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	token, err := s.signedToken(u.Id)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, generated.AuthorizationResponse{
		User:  UserToJSON(u),
		Token: token,
	})
}

func (s *Server) UserByUsername(c *gin.Context) {
	var req generated.UserByUsernameParams
	err := c.ShouldBindQuery(&req)
	if err != nil {
		Error(c, http.StatusBadRequest)
		return
	}

	u, err := s.app.UserByUsername(req.Username)
	if errors.Is(err, application.ErrUserNotFound) {
		Error(c, http.StatusNotFound, "not-found")
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, UserToJSON(u))
}
