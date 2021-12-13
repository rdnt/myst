package rest

import (
	"github.com/gin-gonic/gin"
)

func (api *API) initRoutes(r *gin.RouterGroup) {
	r.POST("/keystores", api.CreateKeystore)
	r.POST("/keystore/:keystoreId", api.UnlockKeystore)
	r.GET("/keystore/:keystoreId", api.Keystore)
}
