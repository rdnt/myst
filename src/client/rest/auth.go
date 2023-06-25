package rest

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"myst/src/client/application"
	"myst/src/client/rest/generated"
)

func sessionId(c *gin.Context) []byte {
	// GetCurrentUserID returns the username of the currently logged-in user
	sid, ok := c.Get("sessionId")
	if !ok {
		return nil
	}

	return sid.([]byte)
}

func (s *Server) Authenticate(c *gin.Context) {
	var req generated.AuthenticateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest)
		return
	}

	sessionId, err := s.app.Authenticate(req.Password)
	if errors.Is(err, application.ErrInitializationRequired) {
		Error(c, http.StatusConflict)
		return
	} else if errors.Is(err, application.ErrAuthenticationFailed) {
		Error(c, http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	sid := base64.StdEncoding.EncodeToString(sessionId)

	c.JSON(http.StatusOK, sid)
}

func (s *Server) CurrentUser(c *gin.Context) {
	sid := sessionId(c)

	u, err := s.app.CurrentUser(sid)
	if errors.Is(err, application.ErrCredentialsNotFound) {
		Error(c, http.StatusNotFound)
		return
	} else if u == nil {
		Error(c, http.StatusUnauthorized)
		return
	}

	restUser, err := s.userToJSON(sid, *u)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restUser)
}
