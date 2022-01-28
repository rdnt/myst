package http

import (
	"github.com/gin-gonic/gin"
)

func (api *API) initRoutes(r *gin.RouterGroup) {
	r.POST("/authenticate", api.Authenticate)
	r.GET("/health", api.HealthCheck)
	r.POST("/keystores", api.CreateKeystore)
	r.GET("/keystores", api.Keystores)
	r.GET("/keystore/:keystoreId", api.Keystore)
	r.POST("/keystore/:keystoreId/entries", api.CreateEntry)
}
