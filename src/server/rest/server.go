package rest

//go:generate oapi-codegen --config oapi-codegen-models.yaml openapi.json
//go:generate oapi-codegen --config oapi-codegen-client.yaml openapi.json

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	prometheus "github.com/zsais/go-gin-prometheus"

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
		p := prometheus.NewPrometheus("rest-api")
		p.Use(r)
	}

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, ":)")
	})

	api := r.Group("api")

	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]any{
			"status": "healthy",
		})
	})

	api.POST("/auth/login", s.Login)
	api.POST("/auth/register", s.Register)

	sec := api.Group("")
	sec.Use(s.authenticationMiddleware())

	sec.GET("/user", s.UserByUsername)

	sec.POST("/keystores", s.CreateKeystore)
	sec.GET("/keystore/:keystoreId", s.Keystore)
	sec.PATCH("/keystore/:keystoreId", s.UpdateKeystore)
	sec.DELETE("/keystore/:keystoreId", s.DeleteKeystore)
	sec.GET("/keystores", s.Keystores)

	sec.POST("/keystore/:keystoreId/invitations", s.CreateInvitation)
	sec.GET("/invitation/:invitationId", s.Invitation)
	sec.PATCH("/invitation/:invitationId", s.AcceptInvitation)
	sec.POST("/invitation/:invitationId", s.FinalizeInvitation)
	sec.DELETE("/invitation/:invitationId", s.DeleteInvitation)
	sec.GET("/invitations", s.Invitations)

	return s
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
