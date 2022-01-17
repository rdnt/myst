package http

import (
	"github.com/gin-gonic/gin"
)

func (api *API) initRoutes(r *gin.RouterGroup) {
	r.GET("/health", api.HealthCheck)
	r.POST("/keystores", api.CreateKeystore)
	r.POST("/keystore/:keystoreId", api.UnlockKeystore)
	r.GET("/keystore/:keystoreId", api.Keystore)
	r.POST("/keystore/:keystoreId/entries", api.CreateEntry)
}
