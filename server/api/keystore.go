package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"myst/server/database"
	"myst/server/keystore"
)

func GetKeystore(c *gin.Context) {
	data, err := database.GetKeystore("b2168f97-f2c7-4a25-a8f5-6d985cce9a65")
	if err != nil {
		fmt.Println(err)
		Error(c, 403, nil)
		return
	}
	// raw, err := keystore.Load("keystore.mst")
	// if err != nil {
	// 	HTTPError(c, 403)
	// 	return
	// }
	enc := fmt.Sprintf("%x", data)
	Success(c, enc)
}

func UpdateKeystore(c *gin.Context) {
	data, err := keystore.Load("keystore.mst")
	if err != nil {
		fmt.Println(err)
		Error(c, 403, nil)
		return
	}
	err = database.SetKeystore("b2168f97-f2c7-4a25-a8f5-6d985cce9a65", data)
	if err != nil {
		fmt.Println(err)
		Error(c, 403, nil)
		return
	}
	Success(c, "succ")
}

func NewUUIDv4(c *gin.Context) {
	uid := uuid.New()
	Success(c, uid)
}
