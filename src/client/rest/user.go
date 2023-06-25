package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"myst/src/client/rest/generated"
)

func (s *Server) Register(c *gin.Context) {
	sid := sessionId(c)

	var req generated.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest)
		return
	}

	u, err := s.app.Register(sid, req.Username, req.Password)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	restUser, err := s.userToJSON(sid, u)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, restUser)
}
