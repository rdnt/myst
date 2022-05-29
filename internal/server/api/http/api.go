package http

//go:generate oapi-codegen -package generated -generate types -o generated/types.gen.go openapi.json
//go:generate oapi-codegen -package generated -generate client -o generated/client.gen.go openapi.json

import (
	"io/ioutil"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	//prometheus "github.com/zsais/go-gin-prometheus"

	application "myst/internal/server"
	"myst/pkg/config"
	"myst/pkg/logger"
)

var (
	log               = logger.New("router", logger.Cyan)
	jwtCookieLifetime = 604800
	jwtSecretKey      = "dzl6JEMmRilKQE1jUWZUalduWnI0dTd4IUElRCpHLUs=" //  TODO: os.Getenv("JWT_SECRET_KEY")
)

type API struct {
	*gin.Engine
	app application.Application
}

func New(app application.Application) *API {
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
		//p := prometheus.NewPrometheus("gin")
		//p.Use(r)
	}

	// Attach static serve middleware for / and /assets
	r.Use(static.Serve("/", static.LocalFile("static", false)))
	r.Use(static.Serve("/assets", static.LocalFile("assets", false)))

	r.Use(
		cors.New(
			cors.Options{
				AllowOriginFunc: func(origin string) bool {
					return true
				},
				//AllowedOrigins: []string{"http://localhost:80", "http://localhost:8082"},
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
			},
		),
	)

	api.initRoutes(r.Group("api"))

	return api
}

func (api *API) initRoutes(g *gin.RouterGroup) {
	g.POST("/auth/login", api.Login)
	g.POST("/auth/register", api.Register)

	sec := g.Group("")
	sec.Use(Authentication())

	// keystore
	sec.POST("/keystores", api.CreateKeystore)
	sec.GET("/keystore/:keystoreId", api.Keystore)
	sec.GET("/keystores", api.Keystores)

	// invitation
	sec.POST("/keystore/:keystoreId/invitations", api.CreateInvitation)
	sec.GET("/invitation/:invitationId", api.Invitation)
	sec.PATCH("/invitation/:invitationId", api.AcceptInvitation)
	sec.POST("/invitation/:invitationId", api.FinalizeInvitation)
	sec.GET("/invitations", api.Invitations)
}

func (api *API) Run(addr string) error {
	log.Println("starting app on port :8080")

	api.app.Start()

	log.Println("app started")
	return api.Engine.Run(addr)
}
