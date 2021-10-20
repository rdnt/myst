package api

//go:generate oapi-codegen -package generated -generate gin-server -o generated/server.gen.go ../../../api/openapi.json
//go:generate oapi-codegen -package generated -generate types -o generated/types.gen.go ../../../api/openapi.json
//go:generate oapi-codegen -package generated -generate client -o generated/client.gen.go ../../../api/openapi.json

import (
	"github.com/gin-gonic/gin"

	cors "github.com/rs/cors/wrapper/gin"

	"myst/app/server"
	"myst/app/server/api/generated"
)

type API struct {
	*gin.Engine
	app *server.Application
}

func (api *API) CreateKeystoreInvitation(c *gin.Context, keystoreId string) {
	var params generated.CreateKeystoreInvitationRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		panic(err)
	}

	inv, err := api.app.CreateKeystoreInvitation(
		"rdnt",
		params.InviteeId,
		keystoreId,
		params.PublicKey,
	)
	if err != nil {
		panic(err)
	}

	c.JSON(200, generated.Invitation{
		Id: inv.Id(),
	})
}

func New(app *server.Application) *API {
	api := new(API)
	api.app = app

	api.Engine = gin.Default()

	generated.HandlerWithOptions(api, generated.GinServerOptions{
		BaseURL:     "/api",
		BaseRouter:  api.Engine,
		Middlewares: nil,
	})

	api.Engine.Use(cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		//AllowedOrigins: []string{"http://localhost:80", "http://localhost:8081"},
		//// TODO allow more methods (DELETE?)
		//AllowedMethods: []string{http.MethodGet, http.MethodPost},
		//// TODO expose ratelimiting headers
		//ExposedHeaders: []string{},
		//// TODO check if we can disable this on release mode so that no
		//// authorization tokens are passed on to the frontend.
		//// No harm, but no need either.
		//// Required to pass authentication headers on development environment
		//AllowCredentials: true,
		Debug: true, // too verbose, only enable for testing CORS
	}))

	return api
}
