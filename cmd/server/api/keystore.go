package api

import (
	"myst/pkg/keystore"

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
		Error(c, 400, nil)
		return
	}

	e := keystore.NewEncrypted(data.Payload, data.Password)
	e.Save()
}
