package api

//go:generate oapi-codegen -package generated -generate gin-server -o generated/server.gen.go ../../../api/openapi.json
//go:generate oapi-codegen -package generated -generate types -o generated/types.gen.go ../../../api/openapi.json
//go:generate oapi-codegen -package generated -generate client -o generated/client.gen.go ../../../api/openapi.json

import (
	"github.com/gin-gonic/gin"

	"myst/app/server"
	"myst/app/server/api/generated"
)

type API struct {
	*gin.Engine
	app server.Application
}

func (api *API) CreateKeystoreInvitation(c *gin.Context, keystoreId string) {
	panic("implement me")
}

func New() *API {
	api := new(API)

	api.Engine = gin.Default()

	generated.HandlerWithOptions(api, generated.GinServerOptions{
		BaseURL:     "/api",
		BaseRouter:  api.Engine,
		Middlewares: nil,
	})

	return api
}
