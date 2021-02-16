package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"myst/server/keystore"
	"myst/server/logger"
	"myst/server/responses"
	"strings"
)

func AuthenticationRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		// Remove the "Bearer" prefix
		parts := strings.Split(auth, "Bearer")
		if len(parts) != 2 {
			rsp := responses.NewErrorResponse(400, nil)
			c.JSON(400, rsp)
			c.Abort()
			return
		}
		// Trim the space that separated Bearer from the JWT
		auth = strings.TrimSpace(parts[1])

		if auth == "" {
			fmt.Println("authorization required")
			rsp := responses.NewErrorResponse(403, nil)
			c.JSON(403, rsp)
			c.Abort()
			return
		}

		raw, err := keystore.Load("keystore.mst")
		if err != nil {
			logger.Errorf("KEYSTORE", err)
			Error(c, 500, nil)
			c.Abort()
			return
		}

		_, err = keystore.DecryptOld(raw, auth)
		if err != nil {
			logger.Errorf("KEYSTORE", err)
			Error(c, 403, nil)
			c.Abort()
			return
		}

		c.Set("master_password", auth)
		// Pass onto the next handler
		c.Next()
	}
}

// Init creates all the HTTP routes
func Init(r *gin.Engine) {
	// r.GET("/", LoadUI)

	api := r.Group("/api")
	api.GET("/uuid4", NewUUIDv4)

	g := api.Group("/entries")
	g.Use(AuthenticationRequired())
	{
		g.GET("", GetEntries)
		// g.PUT("", AddEntry)
	}
	g = api.Group("/entry/:id")
	g.Use(AuthenticationRequired())
	{
		// g.GET("", GetEntry)
		// g.DELETE("", RemoveEntry)
	}

	g = api.Group("/keystore/:id")
	// g.Use(AuthenticationRequired())
	{
		g.GET("", GetKeystore)
		g.POST("", UpdateKeystore)
	}

}

func LoadUI(c *gin.Context) {
	c.File("../static")
	c.Abort()
}

// type RequireMasterPassword struct {
// 	MasterPassword string `header:"Authorization" binding:"required,regex=master_password"`
// }
//
// type GetEntriesData struct {
// 	RequireMasterPassword
// }

func GetEntries(c *gin.Context) {
	pass := c.GetString("master_password")
	raw, err := keystore.Load("keystore.mst")
	if err != nil {
		logger.Errorf("KEYSTORE", err)
		Error(c, 500, nil)
		return
	}

	ks, err := keystore.DecryptOld(raw, pass)
	if err != nil {
		logger.Errorf("KEYSTORE", err)
		Error(c, 403, nil)
		return
	}

	err = ks.Save("keystore.mst", pass)
	if err != nil {
		logger.Errorf("KEYSTORE", err)
		Error(c, 500, nil)
		return
	}

	entries := ks.Entries
	Success(c, entries)
}

// type GetEntryUri struct {
// 	ID string `uri:"id" binding:"required,regex=uuid"`
// }
//
// type GetEntryData struct {
// 	RequireMasterPassword
// }
//
// func GetEntry(c *gin.Context) {
// 	pass := c.GetString("master_password")
// 	var uri GetEntryUri
// 	err := c.ShouldBindUri(&uri)
// 	if err {
// 		fmt.Println("dafuq")
// 		ValidationError(c, err1)
// 		return
// 	}
// 	raw, err := keystore.Load("keystore.mst")
// 	if err != nil {
// 		logger.Errorf("KEYSTORE", err)
// 		HTTPError(c, 500)
// 		return
// 	}
// 	ks, err := keystore.Decrypt(raw, data.MasterPassword)
// 	if err != nil {
// 		logger.Errorf("KEYSTORE", err)
// 		HTTPError(c, 403)
// 		return
// 	}
// 	entry, err := ks.Get(uri.ID)
// 	if err == keystore.ErrNoEntry {
// 		HTTPError(c, 404)
// 		return
// 	} else if err != nil {
// 		logger.Errorf("KEYSTORE", err)
// 		HTTPError(c, 500)
// 		return
// 	}
// 	Success(c, entry)
// }
//
// type AddEntryData struct {
// 	RequireMasterPassword
// 	Email    string `form:"email" binding:"required"`
// 	Password string `form:"password" binding:"required"`
// }
//
// func AddEntry(c *gin.Context) {
// 	var data AddEntryData
// 	err := c.ShouldBind(&data)
// 	if err != nil {
// 		ValidationError(c, err)
// 		return
// 	}
// 	enc, err := keystore.Load("keystore.mst")
// 	if err != nil {
// 		logger.Errorf("KEYSTORE", err)
// 		HTTPError(c, 500)
// 		return
// 	}
//
// 	ks, err := keystore.Decrypt(enc, data.MasterPassword)
// 	if err != nil {
// 		logger.Errorf("KEYSTORE", err)
// 		HTTPError(c, 403)
// 		return
// 	}
// 	entry, err := ks.Add(data.Email, data.Password)
// 	if err != nil {
// 		logger.Errorf("KEYSTORE", err)
// 		HTTPError(c, 500)
// 		return
// 	}
// 	err = ks.Save("keystore.mst", data.MasterPassword)
// 	if err != nil {
// 		logger.Errorf("KEYSTORE", err)
// 		HTTPError(c, 500)
// 		return
// 	}
// 	Success(c, entry)
// }
//
// type RemoveEntryUri struct {
// 	ID string `uri:"id" binding:"required,regex=uuid"`
// }
//
// type RemoveEntryData struct {
// 	RequireMasterPassword
// }
//
// func RemoveEntry(c *gin.Context) {
// 	var uri RemoveEntryUri
// 	var data RemoveEntryData
// 	err1 := c.ShouldBindUri(&uri)
// 	err2 := c.ShouldBind(&data)
// 	if err1 != nil || err2 != nil {
// 		ValidationError(c, err1, err2)
// 		return
// 	}
// 	raw, err := keystore.Load("keystore.mst")
// 	if err != nil {
// 		logger.Errorf("KEYSTORE", err)
// 		HTTPError(c, 500)
// 		return
// 	}
//
// 	ks, err := keystore.Decrypt(raw, data.MasterPassword)
// 	if err != nil {
// 		logger.Errorf("KEYSTORE", err)
// 		HTTPError(c, 403)
// 		return
// 	}
// 	removed, err := ks.Remove(uri.ID)
// 	if err == keystore.ErrNoEntry {
// 		HTTPError(c, 404)
// 		return
// 	} else if err != nil || removed == false {
// 		logger.Errorf("KEYSTORE", err)
// 		HTTPError(c, 500)
// 		return
// 	}
// 	err = ks.Save("keystore.mst", data.MasterPassword)
// 	if err != nil {
// 		HTTPError(c, 500)
// 		return
// 	}
// 	Success(c, "Resource removed successfully.")
// }
//
// func IsAuthenticated(c *gin.Context) {
//
// }
