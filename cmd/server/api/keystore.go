package api

import (
	"myst/pkg/keystore"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateKeystore(c *gin.Context) {

}

func CreateEncryptedKeystore(c *gin.Context) {
	var data struct {
		Password string `form:"password" binding:"required"`
		Payload  string `form:"payload" binding:"required"`
	}
	err := c.ShouldBind(&data)
	if err != nil {
		Error(c, 400, err)
		return
	}

	e, err := keystore.NewEncrypted(data.Payload, data.Password)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	e.Save()
}
