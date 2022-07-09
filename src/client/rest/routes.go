package rest

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) initRoutes(r *gin.RouterGroup) {
	r.GET("/health", s.HealthCheck)
	r.POST("/authenticate", s.Authenticate)
	r.POST("/auth/login", s.Login)
	r.POST("/auth/register", s.Register)
	r.POST("/keystores", s.CreateKeystore)
	r.GET("/keystores", s.Keystores)
	r.GET("/keystore/:keystoreId", s.Keystore)
	r.DELETE("/keystore/:keystoreId", s.DeleteKeystore)
	r.POST("/keystore/:keystoreId/entries", s.CreateEntry)
	r.PATCH("/keystore/:keystoreId/entry/:entryId", s.UpdateEntry)
	r.DELETE("/keystore/:keystoreId/entry/:entryId", s.DeleteEntry)
	r.GET("/invitations", s.GetInvitations)
	r.POST("/keystore/:keystoreId/invitations", s.CreateInvitation)
	r.PATCH("/invitation/:invitationId", s.AcceptInvitation)
	r.DELETE("/invitation/:invitationId", s.DeclineOrCancelInvitation)
	r.POST("/invitation/:invitationId", s.FinalizeInvitation)
	r.GET("/user", s.CurrentUser)
	r.POST("/enclave", s.CreateEnclave)
	r.GET("/enclave", s.Enclave)
	r.POST("/import", s.Import)
}
