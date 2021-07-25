package api

import (
	"github.com/gin-gonic/gin"
	"myst/cmd/server/internal/keystore"
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

	store := keystore.New(data.Keystore)
	err = store.Save()
	if err != nil {
		panic(err)
	}

	uk := userkeystore.New(uid, store.ID, data.Key)
	err = uk.Save()
	if err != nil {
		panic(err)
	}
}
