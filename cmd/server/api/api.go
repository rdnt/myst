package api

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"myst/pkg/logger"
)

var log = logger.New("router", logger.GreenFg)

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

	AuthenticatedAPI(api.Group("").Use(Authentication()))
}

func AuthenticatedAPI(api gin.IRoutes) {
	{
		api.GET("/pong", func(c *gin.Context) {
			c.String(http.StatusOK, "PING!")
		})
	}
}
