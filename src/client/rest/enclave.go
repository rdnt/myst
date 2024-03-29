package rest

import (
	"encoding/base64"
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

	sessionId, err := s.app.Initialize(req.Password)
	if errors.Is(err, application.ErrEnclaveExists) {
		Error(c, http.StatusConflict)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	sid := base64.StdEncoding.EncodeToString(sessionId)

	c.JSON(http.StatusCreated, sid)
}

func (s *Server) Enclave(c *gin.Context) {
	sid := sessionId(c)

	exists, err := s.app.IsInitialized(sid)
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
