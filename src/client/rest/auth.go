package rest

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"myst/src/client/application"
	"myst/src/client/rest/generated"
)

func (s *Server) Authenticate(c *gin.Context) {
	var req generated.AuthenticateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest)
		return
	}

	err = s.app.Authenticate(req.Password)
	if errors.Is(err, application.ErrAuthenticationFailed) {
		Error(c, http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) Register(c *gin.Context) {
	var req generated.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest)
		return
	}

	u, err := s.app.Register(req.Username, req.Password)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	restUser, err := s.userToRest(u)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, restUser)
}
