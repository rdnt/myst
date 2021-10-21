package api

import (
	"github.com/gin-gonic/gin"

	"myst/cmd/server/api"
)

func initRoutes(r *gin.RouterGroup) {
	r.POST("/keystore/:keystoreId/invitations", api.CreateKeystoreInvitation)
}
