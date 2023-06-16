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

func (s *Server) CurrentUser(c *gin.Context) {
	u, err := s.app.CurrentUser()
	if errors.Is(err, application.ErrCredentialsNotFound) {
		Error(c, http.StatusNotFound)
		return
	} else if u == nil {
		Error(c, http.StatusUnauthorized)
		return
	}

	restUser, err := s.userToJSON(*u)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restUser)
}
