package rest

//go:generate oapi-codegen --config oapi-codegen-models.yaml openapi.json
//go:generate oapi-codegen --config oapi-codegen-client.yaml openapi.json

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	cors "github.com/rs/cors/wrapper/gin"
	// prometheus "github.com/zsais/go-gin-prometheus"

	"myst/pkg/config"
	"myst/pkg/logger"
	"myst/pkg/server"
	"myst/src/server/application"
	"myst/src/server/rest/generated"
)

var (
	log               = logger.New("router", logger.Cyan)
	jwtLifetime int64 = 604800
)

type Server struct {
	*gin.Engine
	app           application.Application
	server        *server.Server
	jwtSigningKey []byte
}

func NewServer(app application.Application, jwtSigningKey []byte) *Server {
	s := &Server{
		app:           app,
		jwtSigningKey: jwtSigningKey,
	}

	// Set gin mode
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	r := gin.New()
	s.Engine = r

	// Do not redirect folders to trailing slash
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	// Custom PrintRouteFunc
	gin.DebugPrintRouteFunc = printRoutes

	// always use recovery middleware
	r.Use(recoveryMiddleware)

	// custom logging middleware
	r.Use(loggerMiddleware)

	// error 404 handling
	r.NoRoute(noRouteMiddleware)

	// metrics
	if config.Debug {
		// p := prometheus.NewPrometheus("gin")
		// p.Use(r)
	}

	// TODO @rdnt: @@@ server also doesn't need a UI.
	// Attach static serve middleware for / and /assets
	r.Use(static.Serve("/", static.LocalFile("static", false)))
	r.Use(static.Serve("/assets", static.LocalFile("assets", false)))

	// TODO @rdnt: @@@ does the server really need CORS? I think not.
	r.Use(
		cors.New(
			cors.Options{
				AllowOriginFunc: func(origin string) bool {
					// TODO: @rdnt @@@ fix ASAP
					return true
				},
				// AllowedOrigins: []string{"http://localhost:80", "http://localhost:8082"},
				// // TODO allow more methods (DELETE?)
				// AllowedMethods: []string{rest.MethodGet, rest.MethodPost},
				// // TODO expose ratelimiting headers
				// ExposedHeaders: []string{},
				// // TODO check if we can disable this on release mode so that no
				// // authentication tokens are passed on to the frontend.
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
	g.POST("/auth/login", s.Login)
	g.POST("/auth/register", s.Register)

	sec := g.Group("")
	sec.Use(s.authenticationMiddleware())

	sec.GET("/user", s.UserByUsername)

	sec.POST("/keystores", s.CreateKeystore)
	sec.PATCH("/keystore/:keystoreId", s.UpdateKeystore)
	sec.DELETE("/keystore/:keystoreId", s.DeleteKeystore)
	sec.GET("/keystores", s.Keystores)

	sec.POST("/keystore/:keystoreId/invitations", s.CreateInvitation)
	sec.GET("/invitation/:invitationId", s.Invitation)
	sec.PATCH("/invitation/:invitationId", s.AcceptInvitation)
	sec.POST("/invitation/:invitationId", s.FinalizeInvitation)
	sec.DELETE("/invitation/:invitationId", s.DeleteInvitation)
	sec.GET("/invitations", s.Invitations)
}

func (s *Server) Start(addr string) error {
	log.Println("app started on", addr)

	httpServer, err := server.New(addr, s.Engine)
	if err != nil {
		return errors.WithMessage(err, "failed to create http server")
	}

	s.server = httpServer

	return nil
}

func (s *Server) Stop() error {
	return s.server.Stop()
}

func Error(c *gin.Context, statusCode int, errorCodeAndOptionalMessage ...string) {
	code := ""
	msg := ""

	if len(errorCodeAndOptionalMessage) > 0 {
		code = errorCodeAndOptionalMessage[0]
	}

	if len(errorCodeAndOptionalMessage) > 1 {
		msg = errorCodeAndOptionalMessage[1]
	}

	if code == "" {
		code = fmt.Sprintf("%d", statusCode)
		msg = http.StatusText(statusCode)
	}

	c.JSON(statusCode, generated.Error{
		Code:    code,
		Message: msg,
	})
}
