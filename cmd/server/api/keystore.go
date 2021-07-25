package api

import (
	"github.com/gin-gonic/gin"
	"myst/userkeystore"
)

func CreateKeystore(c *gin.Context) {
	uid := GetCurrentUser(c)

	var data struct {
		Key      []byte `form:"key" binding:"required"`
		Keystore []byte `form:"keystore" binding:"required"`
	}
	err := c.ShouldBind(&data)
	if err != nil {
		Error(c, 400, err)
		return
	}

	uk := userkeystore.New(uid, data.Key, data.Keystore)

	uk.Save()
}
