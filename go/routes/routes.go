package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sht/shtdev/go/keystore"
)

// Init creates all the HTTP routes
func Init(r *gin.Engine) {

	ks := r.Group("/entries")
	{
		ks.GET("", GetEntries)
		ks.GET("/get", GetEntry)
		ks.GET("/add", AddEntry)
		ks.GET("/remove", RemoveEntry)
	}

}

type GetEntriesData struct {
	MasterPassword string `form:"master_password" binding:"required"`
}

func GetEntries(c *gin.Context) {
	var data GetEntriesData
	err := c.ShouldBind(&data)
	if err != nil {
		HTTPError(c, 400)
		return
	}
	enc, err := keystore.Load("keystore.mst")
	if err != nil {
		fmt.Println(err)
		HTTPError(c, 500)
		return
	}

	ks, err := keystore.Decrypt(enc, data.MasterPassword)
	if err != nil {
		fmt.Println(err)
		HTTPError(c, 500)
		return
	}

	entries := ks.Entries
	Success(c, entries)
}

type GetEntryData struct {
	MasterPassword string `form:"master_password" binding:"required"`
	ID             string `form:"id" binding:"required"`
}

func GetEntry(c *gin.Context) {
	var data GetEntryData
	err := c.ShouldBind(&data)
	if err != nil {
		HTTPError(c, 400)
		return
	}
	enc, err := keystore.Load("keystore.mst")
	if err != nil {
		fmt.Println(err)
		HTTPError(c, 500)
		return
	}

	ks, err := keystore.Decrypt(enc, data.MasterPassword)
	if err != nil {
		fmt.Println(err)
		HTTPError(c, 500)
		return
	}
	entry, err := ks.Get(data.ID)
	if err == keystore.ErrNoEntry {
		HTTPError(c, 404)
		return
	} else if err != nil {
		HTTPError(c, 500)
		return
	}
	Success(c, entry)
}

type AddEntryData struct {
	MasterPassword string `form:"master_password" binding:"required"`
	Email          string `form:"email" binding:"required"`
	Password       string `form:"password" binding:"required"`
}

func AddEntry(c *gin.Context) {
	var data AddEntryData
	err := c.ShouldBind(&data)
	if err != nil {
		HTTPError(c, 400)
		return
	}
	enc, err := keystore.Load("keystore.mst")
	if err != nil {
		fmt.Println(err)
		HTTPError(c, 500)
		return
	}

	ks, err := keystore.Decrypt(enc, data.MasterPassword)
	if err != nil {
		fmt.Println(err)
		HTTPError(c, 500)
		return
	}
	entry, err := ks.Add(data.Email, data.Password)
	if err != nil {
		HTTPError(c, 500)
		return
	}
	err = ks.Save("keystore.mst", data.MasterPassword)
	if err != nil {
		HTTPError(c, 500)
		return
	}
	Success(c, entry)
}

type RemoveEntryData struct {
	MasterPassword string `form:"master_password" binding:"required"`
	ID             string `form:"id" binding:"required"`
}

func RemoveEntry(c *gin.Context) {
	var data RemoveEntryData
	err := c.ShouldBind(&data)
	if err != nil {
		HTTPError(c, 400)
		return
	}
	enc, err := keystore.Load("keystore.mst")
	if err != nil {
		fmt.Println(err)
		HTTPError(c, 500)
		return
	}

	ks, err := keystore.Decrypt(enc, data.MasterPassword)
	if err != nil {
		fmt.Println(err)
		HTTPError(c, 500)
		return
	}
	removed, err := ks.Remove(data.ID)
	if err == keystore.ErrNoEntry {
		HTTPError(c, 404)
		return
	} else if err != nil || removed == false {
		HTTPError(c, 500)
		return
	}
	err = ks.Save("keystore.mst", data.MasterPassword)
	if err != nil {
		HTTPError(c, 500)
		return
	}
	Success(c, "Resource removed successfully.")
}
