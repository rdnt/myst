package api

import (
	"github.com/gin-gonic/gin"

	"myst/pkg/logger"
)

var (
	log = logger.New("router", logger.GreenFg)
)

// Init creates all the API endpoints for the server
func Init(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		Success(c, "Pong!")
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
		Success(c, "PING!")
	})

	api.POST("/keystorerepo", CreateKeystore)
	api.GET("/domain/:domain", GetKeystore)

	// domain invitations
	api.POST("/domain/:domain/invitations", CreateKeystoreInvitation)
}
