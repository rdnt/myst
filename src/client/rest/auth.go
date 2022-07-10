package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"myst/src/client/application"
	"myst/src/client/rest/generated"
)

func (s *Server) Authenticate(c *gin.Context) {
	var req generated.AuthenticateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := s.app.Authenticate(req.Password); err == application.ErrAuthenticationFailed {
		c.Status(http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) Login(c *gin.Context) {
	var req generated.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if _, err := s.app.SignIn(req.Username, req.Password); err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) Register(c *gin.Context) {
	var req generated.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	u, err := s.app.Register(req.Username, req.Password)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, UserToRest(u))
}
