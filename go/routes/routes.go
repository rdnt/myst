package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sht/myst/go/keystore"
	"github.com/sht/myst/go/logger"
)

// Init creates all the HTTP routes
func Init(r *gin.Engine) {

	g := r.Group("/entries")
	{
		g.GET("", GetEntries)
		g.PUT("", AddEntry)
	}
	g = r.Group("/entry/:id")
	{
		g.GET("", GetEntry)
		g.DELETE("", RemoveEntry)
	}

}

type RequireMasterPassword struct {
	MasterPassword string `form:"master_password" binding:"required,regex=master_password"`
}

type GetEntriesData struct {
	RequireMasterPassword
}

func GetEntries(c *gin.Context) {
	var data GetEntriesData
	err := c.ShouldBind(&data)
	if err != nil {
		fmt.Println(err)
		ValidationError(c, err)
		return
	}
	raw, err := keystore.Load("keystore.mst")
	if err != nil {
		logger.Errorf("KEYSTORE", err)
		HTTPError(c, 500)
		return
	}

	ks, err := keystore.Decrypt(raw, data.MasterPassword)
	if err != nil {
		logger.Errorf("KEYSTORE", err)
		HTTPError(c, 403)
		return
	}

	entries := ks.Entries
	Success(c, entries)
}

type GetEntryUri struct {
	ID string `uri:"id" binding:"required,regex=uuid"`
}

type GetEntryData struct {
	RequireMasterPassword
}

func GetEntry(c *gin.Context) {
	var uri GetEntryUri
	var data GetEntryData
	err1 := c.ShouldBindUri(&uri)
	err2 := c.ShouldBind(&data)
	if err1 != nil || err2 != nil {
		fmt.Println("dafuq")
		ValidationError(c, err1, err2)
		return
	}
	raw, err := keystore.Load("keystore.mst")
	if err != nil {
		logger.Errorf("KEYSTORE", err)
		HTTPError(c, 500)
		return
	}
	ks, err := keystore.Decrypt(raw, data.MasterPassword)
	if err != nil {
		logger.Errorf("KEYSTORE", err)
		HTTPError(c, 403)
		return
	}
	entry, err := ks.Get(uri.ID)
	if err == keystore.ErrNoEntry {
		HTTPError(c, 404)
		return
	} else if err != nil {
		logger.Errorf("KEYSTORE", err)
		HTTPError(c, 500)
		return
	}
	Success(c, entry)
}

type AddEntryData struct {
	RequireMasterPassword
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func AddEntry(c *gin.Context) {
	var data AddEntryData
	err := c.ShouldBind(&data)
	if err != nil {
		ValidationError(c, err)
		return
	}
	enc, err := keystore.Load("keystore.mst")
	if err != nil {
		logger.Errorf("KEYSTORE", err)
		HTTPError(c, 500)
		return
	}

	ks, err := keystore.Decrypt(enc, data.MasterPassword)
	if err != nil {
		logger.Errorf("KEYSTORE", err)
		HTTPError(c, 403)
		return
	}
	entry, err := ks.Add(data.Email, data.Password)
	if err != nil {
		logger.Errorf("KEYSTORE", err)
		HTTPError(c, 500)
		return
	}
	err = ks.Save("keystore.mst", data.MasterPassword)
	if err != nil {
		logger.Errorf("KEYSTORE", err)
		HTTPError(c, 500)
		return
	}
	Success(c, entry)
}

type RemoveEntryUri struct {
	ID string `uri:"id" binding:"required,regex=uuid"`
}

type RemoveEntryData struct {
	RequireMasterPassword
}

func RemoveEntry(c *gin.Context) {
	var uri RemoveEntryUri
	var data RemoveEntryData
	err1 := c.ShouldBindUri(&uri)
	err2 := c.ShouldBind(&data)
	if err1 != nil || err2 != nil {
		ValidationError(c, err1, err2)
		return
	}
	raw, err := keystore.Load("keystore.mst")
	if err != nil {
		logger.Errorf("KEYSTORE", err)
		HTTPError(c, 500)
		return
	}

	ks, err := keystore.Decrypt(raw, data.MasterPassword)
	if err != nil {
		logger.Errorf("KEYSTORE", err)
		HTTPError(c, 403)
		return
	}
	removed, err := ks.Remove(uri.ID)
	if err == keystore.ErrNoEntry {
		HTTPError(c, 404)
		return
	} else if err != nil || removed == false {
		logger.Errorf("KEYSTORE", err)
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
