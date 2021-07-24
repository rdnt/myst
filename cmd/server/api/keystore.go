package api

import (
	"github.com/gin-gonic/gin"
	"myst/pkg/userkeystore"
)

func CreateKeystore(c *gin.Context) {
	uid := GetCurrentUser(c)

	var data struct {
		Key   []byte `form:"key" binding:"required"`
		Store []byte `form:"store" binding:"required"`
	}
	err := c.ShouldBind(&data)
	if err != nil {
		Error(c, 400, err)
		return
	}

	uk := userkeystore.New(uid, data.Key, data.Store)

	uk.Save()
}
