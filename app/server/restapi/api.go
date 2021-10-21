package restapi

//go:generate oapi-codegen -package generated -generate types -o generated/types.gen.go ../../../api/openapi.json

import (
	"encoding/hex"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	prometheus "github.com/zsais/go-gin-prometheus"

	"myst/app/server/core"
	"myst/app/server/restapi/generated"
	"myst/pkg/config"
	"myst/pkg/logger"
)

var log = logger.New("router", logger.Cyan)

type API struct {
	*gin.Engine
	app *core.Application
}

func (api *API) CreateKeystoreInvitation(c *gin.Context) {
	keystoreId := c.Param("keystoreId")

	//log.Debug(keystoreId)

	var params generated.CreateKeystoreInvitationRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		panic(err)
	}

	inviterKey, err := hex.DecodeString(params.PublicKey)
	if err != nil {
		panic(err)
	}

	inv, err := api.app.CreateKeystoreInvitation(
		"rdnt",
		params.InviteeId,
		keystoreId,
		inviterKey,
	)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, generated.Invitation{
		Id: inv.Id(),
	})
}

func (api *API) Run(addr string) error {
	log.Println("starting app on port :8080")

	api.app.Start()

	log.Println("app started")
	return api.Engine.Run(addr)
}

func New(app *core.Application) *API {
	api := new(API)

	api.app = app

	// Set gin mode
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DefaultWriter = ioutil.Discard

	r := gin.New()
	api.Engine = r

	// Do not redirect folders to trailing slash
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	// Custom PrintRouteFunc
	gin.DebugPrintRouteFunc = PrintRoutes

	// always use recovery middleware
	r.Use(Recovery)

	// custom logging middleware
	r.Use(LoggerMiddleware)

	// error 404 handling
	r.NoRoute(NoRoute)

	// metrics
	if config.Debug {
		p := prometheus.NewPrometheus("gin")
		p.Use(r)
	}

	// error 404 handling
	r.NoRoute(NoRoute)

	// Attach static serve middleware for / and /assets
	r.Use(static.Serve("/", static.LocalFile("static", false)))
	r.Use(static.Serve("/assets", static.LocalFile("assets", false)))

	r.Use(cors.New(cors.Options{
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
		Debug: false, // too verbose, only enable for testing CORS
	}))

	api.initRoutes(r.Group("api"))

	return api
}
