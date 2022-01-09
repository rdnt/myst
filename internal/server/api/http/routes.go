package http

import (
	"github.com/gin-gonic/gin"
)

func (api *API) initRoutes(r *gin.RouterGroup) {
	r.POST("/keystore/:keystoreId/invitations", api.CreateKeystoreInvitation)
}
