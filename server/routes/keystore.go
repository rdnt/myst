package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sht/myst/server/database"
	"github.com/sht/myst/server/keystore"
)

func GetKeystore(c *gin.Context) {
	data, err := database.GetKeystore("b2168f97-f2c7-4a25-a8f5-6d985cce9a65")
	if err != nil {
		fmt.Println(err)
		HTTPError(c, 403)
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
		HTTPError(c, 403)
		return
	}
	err = database.SetKeystore("b2168f97-f2c7-4a25-a8f5-6d985cce9a65", data)
	if err != nil {
		fmt.Println(err)
		HTTPError(c, 403)
		return
	}
	Success(c, "succ")
}

func NewUUIDv4(c *gin.Context) {
	uid := uuid.New()
	Success(c, uid)
}
