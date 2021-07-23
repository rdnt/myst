package api

import (
	"github.com/gin-gonic/gin"
	"myst/server/models"
	"net/http"
)

func GetCurrentUserHandler(c *gin.Context) {
	u, err := models.GetUser(c.GetString("user"))
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	Success(c, u.ToRest())
}

func CreateKeystore(c *gin.Context) {
	//name := "new-keystore"
}
