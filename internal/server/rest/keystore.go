package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"myst/internal/server/rest/generated"
)

func (api *Server) CreateKeystore(c *gin.Context) {
	userId := CurrentUser(c)

	var req generated.CreateKeystoreRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		panic(err)
	}

	k, err := api.app.CreateKeystore(req.Name, userId, req.Payload)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, ToJSONKeystore(k))
}

func (api *Server) Keystore(c *gin.Context) {
	userId := CurrentUser(c)
	keystoreId := c.Param("keystoreId")

	k, err := api.app.UserKeystore(userId, keystoreId)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, ToJSONKeystore(k))
}

func (api *Server) Keystores(c *gin.Context) {
	userId := CurrentUser(c)

	ks, err := api.app.UserKeystores(userId)
	if err != nil {
		panic(err)
	}

	gen := []generated.Keystore{}

	for _, k := range ks {
		gen = append(gen, ToJSONKeystore(k))
	}

	c.JSON(http.StatusOK, gen)
}
