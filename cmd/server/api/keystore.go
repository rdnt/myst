package api

import (
	"net/http"

	"myst/cmd/server/keystoreinvitation"

	"myst/user"

	"github.com/gin-gonic/gin"

	"myst/cmd/server/keystore"
	"myst/userkey"
)

type RestUserKeystoreKey struct {
	Keystore *keystore.RestKeystore `json:"keystore"`
	Key      *userkey.RestUserKey   `json:"key"`
}

func CreateKeystore(c *gin.Context) {
	uid := GetCurrentUserID(c)

	var data struct {
		Name     string `json:"name" binding:"required"`
		Key      []byte `json:"key" binding:"required"`
		Keystore []byte `json:"keystore" binding:"required"`
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
	uid := GetCurrentUserID(c)
	keystoreName := c.Param("keystore")

	store, err := keystore.Get("name", keystoreName)
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
		"keystore_id": store.ID,
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

func CreateKeystoreInvitation(c *gin.Context) {
	inviterID := GetCurrentUserID(c)
	keystoreName := c.Param("keystore")

	var data struct {
		Username  string `json:"username" binding:"required"`
		PublicKey []byte `json:"public_key" binding:"required"`
	}
	err := c.ShouldBind(&data)
	if err != nil {
		Error(c, 400, err)
		return
	}

	store, err := keystore.Get("name", keystoreName)
	if err == keystore.ErrNotFound {
		log.Error(err)
		Error(c, http.StatusNotFound, err)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	invitee, err := user.Get("username", data.Username)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	inv, err := keystoreinvitation.New(inviterID, store.ID, invitee.ID, data.PublicKey)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	Success(c, inv.ToRest())
}
