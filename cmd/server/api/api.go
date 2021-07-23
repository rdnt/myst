package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Init creates all the API endpoints for the server
func Init(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "Pong!")
		})
	}
}
