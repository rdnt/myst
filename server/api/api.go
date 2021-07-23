package api

import (
	"github.com/gin-gonic/gin"
)

// Init creates all the HTTP routes
func Init(r *gin.Engine) {
	// r.GET("/", LoadUI)
	api := r.Group("/api")

	//api.GET("/create", func(c *gin.Context) {
	//	err := database.New("1234567890")
	//	fmt.Println(err)
	//})

	auth := api.Group("auth")
	auth.
		POST("/login", LoginHandler).
		POST("/register", RegisterHandler)

	sec := api.Group("").Use(Authentication())
	sec.GET("/user", GetCurrentUserHandler)

	//api.GET("/load", func(c *gin.Context) {
	//	err := database.Load(pass)
	//	if err != nil {
	//		panic(err)
	//	}
	//	Success(c, nil)
	//})

	//api.GET("/get", func(c *gin.Context) {
	//	Success(c, database.Debug())
	//})

	//api.POST("/sites", func(c *gin.Context) {
	//	err := database.AddSite(
	//		keystore.AddSiteOptions{
	//			Label:    "test",
	//			Email:    "test@test.test",
	//			Password: "123123123123",
	//		})
	//	if err != nil {
	//		panic(err)
	//	}
	//	Success(c, nil)
	//})

	//api.GET("/uuid4", NewUUIDv4)

	//g := api.Group("/entries")
	//g.Use(AuthenticationRequired())
	//{
	//	g.GET("", GetEntries)
	//	// g.PUT("", AddEntry)
	//}
	//g = api.Group("/entry/:id")
	//g.Use(AuthenticationRequired())
	//{
	//	// g.GET("", GetEntry)
	//	// g.DELETE("", RemoveEntry)
	//}
	//
	//g = api.Group("/keystore/:id")
	//// g.Use(AuthenticationRequired())
	//{
	//	g.GET("", GetKeystore)
	//	g.POST("", UpdateKeystore)
	//}
}
