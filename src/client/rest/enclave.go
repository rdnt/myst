package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"myst/src/client/application"
	"myst/src/client/rest/generated"
)

func (s *Server) CreateEnclave(c *gin.Context) {
	var req generated.CreateEnclaveRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest)
		return
	}

	err = s.app.Initialize(req.Password)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func (s *Server) Enclave(c *gin.Context) {
	exists, err := s.app.IsInitialized()
	if errors.Is(err, application.ErrAuthenticationRequired) {
		Error(c, http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	if !exists {
		Error(c, http.StatusNotFound)
		return
	}

	c.Status(http.StatusOK)
}
