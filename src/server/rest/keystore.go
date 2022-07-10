package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"myst/src/server/rest/generated"
)

func (s *Server) CreateKeystore(c *gin.Context) {
	userId := CurrentUser(c)

	var req generated.CreateKeystoreRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		panic(err)
	}

	k, err := s.app.CreateKeystore(req.Name, userId, req.Payload)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusCreated, ToJSONKeystore(k))
}

func (s *Server) Keystore(c *gin.Context) {
	userId := CurrentUser(c)
	keystoreId := c.Param("keystoreId")

	k, err := s.app.UserKeystore(userId, keystoreId)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, ToJSONKeystore(k))
}

func (s *Server) Keystores(c *gin.Context) {
	userId := CurrentUser(c)

	ks, err := s.app.UserKeystores(userId)
	if err != nil {
		panic(err)
	}

	gen := []generated.Keystore{}

	for _, k := range ks {
		gen = append(gen, ToJSONKeystore(k))
	}

	c.JSON(http.StatusOK, gen)
}
