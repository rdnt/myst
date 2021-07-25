package api

import (
	"github.com/gin-gonic/gin"
	"myst/cmd/server/internal/keystore"
	"myst/userkey"
	"net/http"
)

type RestUserKeystoreKey struct {
	Keystore *keystore.RestKeystore `json:"keystore"`
	Key      *userkey.RestUserKey   `json:"key"`
}

func CreateKeystore(c *gin.Context) {
	uid := GetCurrentUser(c)

	var data struct {
		Name     string `form:"name" binding:"required"`
		Key      []byte `form:"key" binding:"required"`
		Keystore []byte `form:"keystore" binding:"required"`
	}
	err := c.ShouldBind(&data)
	if err != nil {
		Error(c, 400, err)
		return
	}

	store, err := keystore.New(data.Name, data.Keystore)
	if err != nil {
		panic(err)
	}

	key, err := userkey.New(uid, store.ID, data.Key)
	if err != nil {
		panic(err)
	}

	restStore := RestUserKeystoreKey{
		Keystore: store.ToRest(),
		Key:      key.ToRest(),
	}

	Success(c, restStore)
}

func GetKeystore(c *gin.Context) {
	uid := GetCurrentUser(c)
	keystoreID := c.Param("keystore_id")

	store, err := keystore.Get("id", keystoreID)
	if err == keystore.ErrNotFound {
		log.Error(err)
		Error(c, http.StatusNotFound, err)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	key, err := userkey.Get(map[string]string{
		"user_id":     uid,
		"keystore_id": keystoreID,
	})
	if err == userkey.ErrNotFound {
		log.Error(err)
		Error(c, http.StatusNotFound, err)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	restStore := RestUserKeystoreKey{
		Keystore: store.ToRest(),
		Key:      key.ToRest(),
	}

	Success(c, restStore)
}
