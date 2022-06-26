package http

import (
	"github.com/gin-gonic/gin"
)

func (api *API) initRoutes(r *gin.RouterGroup) {
	r.GET("/health", api.HealthCheck)
	r.POST("/authenticate", api.Authenticate)
	r.POST("/auth/login", api.Login)
	r.POST("/auth/register", api.Register)
	r.POST("/keystores", api.CreateKeystore)
	r.GET("/keystores", api.Keystores)
	r.GET("/keystore/:keystoreId", api.Keystore)
	r.DELETE("/keystore/:keystoreId", api.DeleteKeystore)
	r.POST("/keystore/:keystoreId/entries", api.CreateEntry)
	r.PATCH("/keystore/:keystoreId/entry/:entryId", api.UpdateEntry)
	r.DELETE("/keystore/:keystoreId/entry/:entryId", api.DeleteEntry)
	r.GET("/invitations", api.GetInvitations)
	r.POST("/keystore/:keystoreId/invitations", api.CreateInvitation)
	r.PATCH("/invitation/:invitationId", api.AcceptInvitation)
	r.DELETE("/invitation/:invitationId", api.DeclineOrCancelInvitation)
	r.POST("/invitation/:invitationId", api.FinalizeInvitation)
	r.GET("/user", api.CurrentUser)
	r.POST("/enclave", api.CreateEnclave)
	r.GET("/enclave", api.Enclave)
}
