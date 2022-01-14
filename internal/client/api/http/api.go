package http

//go:generate oapi-codegen -package generated -generate types -o generated/types.gen.go openapi.json

import (
	"errors"
	"io/ioutil"
	"net/http"

	application "myst/internal/client"
	"myst/internal/client/api/http/generated"
	"myst/internal/client/core/domain/keystore/entry"
	"myst/internal/client/core/keystoreservice"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	prometheus "github.com/zsais/go-gin-prometheus"

	"myst/pkg/config"
	"myst/pkg/logger"
)

var log = logger.New("router", logger.Cyan)

type API struct {
	*gin.Engine
	app application.Application
}

func (api *API) CreateKeystore(c *gin.Context) {
	var req generated.CreateKeystoreRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	k, err := api.app.CreateKeystore(
		req.Name,
		req.Passphrase,
	)
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	entries := make([]generated.Entry, len(k.Entries()))

	for i, e := range k.Entries() {
		entries[i] = generated.Entry{
			Id:       e.Id(),
			Label:    e.Label(),
			Username: e.Username(),
			Password: e.Password(),
		}
	}

	Success(
		c, generated.Keystore{
			Id:      k.Id(),
			Name:    k.Name(),
			Entries: entries,
		},
	)
}

func (api *API) UnlockKeystore(c *gin.Context) {
	keystoreId := c.Param("keystoreId")

	var req generated.UnlockKeystoreRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	k, err := api.app.UnlockKeystore(keystoreId, req.Passphrase)
	if errors.Is(err, keystoreservice.ErrAuthenticationFailed) {
		Error(c, http.StatusForbidden, err)
		return
	}
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	entries := make([]generated.Entry, len(k.Entries()))

	for i, e := range k.Entries() {
		entries[i] = generated.Entry{
			Id:       e.Id(),
			Label:    e.Label(),
			Username: e.Username(),
			Password: e.Password(),
		}
	}

	Success(
		c, generated.Keystore{
			Id:      k.Id(),
			Name:    k.Name(),
			Entries: entries,
		},
	)
}

func (api *API) CreateEntry(c *gin.Context) {
	keystoreId := c.Param("keystoreId")

	var req generated.CreateEntryRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	k, err := api.app.Keystore(keystoreId)
	if errors.Is(err, keystoreservice.ErrAuthenticationRequired) {
		Error(c, http.StatusForbidden, err)
		return
	} else if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	e, err := entry.New(
		entry.WithLabel(req.Label),
		entry.WithUsername(req.Username),
		entry.WithPassword(req.Password),
	)
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	err = k.AddEntry(e)
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	err = api.app.UpdateKeystore(k)
	if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	entries := make([]generated.Entry, len(k.Entries()))

	for i, e := range k.Entries() {
		entries[i] = generated.Entry{
			Id:       e.Id(),
			Label:    e.Label(),
			Username: e.Username(),
			Password: e.Password(),
		}
	}

	Success(
		c, generated.Keystore{
			Id:      k.Id(),
			Name:    k.Name(),
			Entries: entries,
		},
	)

}

func (api *API) Keystore(c *gin.Context) {
	keystoreId := c.Param("keystoreId")

	k, err := api.app.Keystore(keystoreId)
	if errors.Is(err, keystoreservice.ErrAuthenticationRequired) {
		Error(c, http.StatusForbidden, err)
		return
	} else if errors.Is(err, keystoreservice.ErrAuthenticationFailed) {
		Error(c, http.StatusForbidden, err)
		return
	} else if err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}

	entries := make([]generated.Entry, len(k.Entries()))

	for i, e := range k.Entries() {
		entries[i] = generated.Entry{
			Id:       e.Id(),
			Label:    e.Label(),
			Username: e.Username(),
			Password: e.Password(),
		}
	}

	Success(
		c, generated.Keystore{
			Id:      k.Id(),
			Name:    k.Name(),
			Entries: entries,
		},
	)
}

func (api *API) HealthCheck(_ *gin.Context) {
	api.app.HealthCheck()
}

func (api *API) Run(addr string) error {
	log.Println("starting app on port :8081")

	api.app.Start()

	log.Println("app started")
	return api.Engine.Run(addr)
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
		p := prometheus.NewPrometheus("gin")
		p.Use(r)
	}

	// error 404 handling
	r.NoRoute(NoRoute)

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
