package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"myst/logger"
)

var (
	log = logger.New("router", logger.GreenFg)
)

// Init creates all the API endpoints for the server
func Init(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "Pong!")
	})

	auth := api.Group("/auth")
	{
		auth.POST("/login", LoginHandler)
		auth.POST("/register", RegisterHandler)
	}

	api.POST("/init", InitData)

	AuthenticatedAPI(api.Group("").Use(Authentication()))
}

func AuthenticatedAPI(api gin.IRoutes) {
	api.GET("/pong", func(c *gin.Context) {
		c.String(http.StatusOK, "PING!")
	})

	api.POST("/keystores", CreateKeystore)
	api.GET("/keystore/:keystore", GetKeystore)

	// keystore invitations
	api.POST("/keystore/:keystore/invitations", CreateKeystoreInvitation)
}
