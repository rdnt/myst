package http

import (
	"github.com/gin-gonic/gin"
)

func (api *API) initRoutes(r *gin.RouterGroup) {
	r.POST("/auth/login", api.Login)

	sec := r.Group("")
	sec.Use(Authentication())

	sec.POST("/keystores", api.CreateKeystore)
	sec.GET("/keystore/:keystoreId", api.Keystore)
	sec.GET("/keystores", api.Keystores)

	sec.POST("/keystore/:keystoreId/invitations", api.CreateInvitation)
	sec.GET("/keystore/:keystoreId/invitation/:invitationId", api.Invitation)
	sec.PATCH("/keystore/:keystoreId/invitation/:invitationId", api.AcceptInvitation)
	sec.POST("/keystore/:keystoreId/invitation/:invitationId", api.FinalizeInvitation)
}
