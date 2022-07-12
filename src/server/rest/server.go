package rest

//go:generate oapi-codegen --config oapi-codegen-models.yaml openapi.json
//go:generate oapi-codegen --config oapi-codegen-client.yaml openapi.json

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	// prometheus "github.com/zsais/go-gin-prometheus"

	"myst/pkg/config"
	"myst/pkg/logger"
	"myst/pkg/server"
	"myst/src/server/application"
)

var (
	log               = logger.New("router", logger.Cyan)
	jwtCookieLifetime = 604800
	jwtSecretKey      = "dzl6JEMmRilKQE1jUWZUalduWnI0dTd4IUElRCpHLUs=" //  TODO: os.Getenv("JWT_SECRET_KEY")
)

type Server struct {
	*gin.Engine
	app    application.Application
	server *server.Server
}

func New(app application.Application) *Server {
	s := new(Server)

	s.app = app

	// Set gin mode
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DefaultWriter = ioutil.Discard

	r := gin.New()
	s.Engine = r

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
		// p := prometheus.NewPrometheus("gin")
		// p.Use(r)
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
				// AllowedOrigins: []string{"http://localhost:80", "http://localhost:8082"},
				// // TODO allow more methods (DELETE?)
				// AllowedMethods: []string{rest.MethodGet, rest.MethodPost},
				// // TODO expose ratelimiting headers
				// ExposedHeaders: []string{},
				// // TODO check if we can disable this on release mode so that no
				// // authorization tokens are passed on to the frontend.
				// // No harm, but no need either.
				// // Required to pass authentication headers on development environment
				// AllowCredentials: true,
				Debug: false, // too verbose, only enable for testing CORS
			},
		),
	)

	s.initRoutes(r.Group("api"))

	return s
}

func (s *Server) initRoutes(g *gin.RouterGroup) {
	g.GET("/debug", func(c *gin.Context) {
		data, err := s.app.Debug()
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, data)
	})
	g.POST("/auth/login", s.Login)
	g.POST("/auth/register", s.Register)

	sec := g.Group("")
	sec.Use(s.Authentication())

	sec.GET("/user", s.User)
	// keystore
	sec.POST("/keystores", s.CreateKeystore)
	sec.GET("/keystore/:keystoreId", s.Keystore)
	sec.GET("/keystores", s.Keystores)

	// invitation
	sec.POST("/keystore/:keystoreId/invitations", s.CreateInvitation)
	sec.GET("/invitation/:invitationId", s.Invitation)
	sec.PATCH("/invitation/:invitationId", s.AcceptInvitation)
	sec.POST("/invitation/:invitationId", s.FinalizeInvitation)
	sec.DELETE("/invitation/:invitationId", s.DeclineOrCancelInvitation)
	sec.GET("/invitations", s.Invitations)
}

func (s *Server) Start(addr string) error {
	log.Println("starting app on", addr)

	s.app.Start()
	log.Println("app started")

	httpServer, err := server.New(addr, s.Engine)
	if err != nil {
		return err
	}

	s.server = httpServer
	return nil
}

func (s *Server) Stop() error {
	var firstErr error

	err := s.server.Stop()
	if err != nil && firstErr == nil {
		firstErr = err
	}

	err = s.app.Stop()
	if err != nil && firstErr == nil {
		firstErr = err
	}

	return firstErr
}
