package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"myst/src/client/rest/generated"
)

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

	restUser, err := s.userToJSON(u)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, restUser)
}
