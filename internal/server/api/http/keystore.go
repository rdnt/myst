package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"myst/internal/server/api/http/generated"
)

func (api *API) CreateKeystore(c *gin.Context) {
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

func (api *API) Keystore(c *gin.Context) {
	userId := CurrentUser(c)
	keystoreId := c.Param("keystoreId")

	k, err := api.app.Keystores.UserKeystore(userId, keystoreId)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, ToJSONKeystore(k))
}

func (api *API) Keystores(c *gin.Context) {
	userId := CurrentUser(c)

	ks, err := api.app.Keystores.UserKeystores(userId)
	if err != nil {
		panic(err)
	}

	gen := []generated.Keystore{}

	for _, k := range ks {
		gen = append(gen, ToJSONKeystore(k))
	}

	c.JSON(http.StatusOK, gen)
}
