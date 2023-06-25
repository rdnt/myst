package rest

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"

	"myst/pkg/logger"
	"myst/pkg/server"
	"myst/src/client/application"
	"myst/src/client/rest/generated"
)

//go:generate oapi-codegen --config oapi-codegen-models.yaml openapi.json
//go:generate oapi-codegen --config oapi-codegen-client.yaml openapi.json
//go:generate oapi-codegen --config oapi-codegen-models.yaml openapi.json
//go:generate npx openapi-typescript-codegen --input openapi.json --output ../../../ui/src/api/generated --client fetch --useOptions --useUnionTypes

var log = logger.New("router", logger.Cyan)

type Server struct {
	*gin.Engine
	app    application.Application
	server *server.Server
}

func NewServer(app application.Application, ui fs.FS) *Server {
	s := new(Server)

	s.app = app

	// Set gin mode
	gin.SetMode(gin.DebugMode)

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
	r.NoRoute(noRouteMiddleware("/", embedFolder(ui, "static")))

	r.Use(
		cors.New(
			cors.Options{
				AllowOriginFunc: func(origin string) bool {
					return true
				},
				AllowedHeaders: []string{"*"},
				AllowedOrigins: []string{
					"http://localhost:80",
					"http://localhost:8082",
					"http://localhost:9092",
				},
				// // TODO allow more methods (DELETE?)
				AllowedMethods: []string{http.MethodHead, http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
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

	api := r.Group("api")

	api.POST("/authenticate", s.Authenticate)
	api.POST("/enclave", s.CreateEnclave)

	sec := api.Group("/")
	sec.Use(s.authenticationMiddleware())

	sec.GET("/enclave", s.Enclave)
	sec.GET("/health", s.HealthCheck)
	sec.POST("/auth/register", s.Register)
	sec.POST("/keystores", s.CreateKeystore)
	sec.GET("/keystores", s.Keystores)
	sec.GET("/keystore/:keystoreId", s.Keystore)
	sec.DELETE("/keystore/:keystoreId", s.DeleteKeystore)
	sec.POST("/keystore/:keystoreId/entries", s.CreateEntry)
	sec.PATCH("/keystore/:keystoreId/entry/:entryId", s.UpdateEntry)
	sec.DELETE("/keystore/:keystoreId/entry/:entryId", s.DeleteEntry)
	sec.GET("/invitations", s.GetInvitations)
	sec.GET("/invitation/:invitationId", s.GetInvitation)
	sec.POST("/keystore/:keystoreId/invitations", s.CreateInvitation)
	sec.PATCH("/invitation/:invitationId", s.AcceptInvitation)
	sec.DELETE("/invitation/:invitationId", s.DeclineOrCancelInvitation)
	sec.POST("/invitation/:invitationId", s.FinalizeInvitation)
	sec.GET("/user", s.CurrentUser)
	//sec.POST("/import", s.DebugImport)

	return s
}
func (s *Server) Start(addr string) error {
	log.Println("app started on", addr)

	httpServer, err := server.New(addr, s.Engine)
	if err != nil {
		return err
	}

	s.server = httpServer
	return nil
}

func (s *Server) Stop() error {
	return s.server.Stop()
}

func (s *Server) HealthCheck(c *gin.Context) {
	sid := sessionId(c)

	err := s.app.HealthCheck(sid)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(_ string, path string) bool {
	_, err := e.Open(path)
	return err == nil
}

func embedFolder(fsEmbed fs.FS, targetPath string) static.ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
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
