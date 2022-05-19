package http

import (
	"io/ioutil"
	"net/http"

	"myst/internal/client/api/http/generated"
	"myst/internal/client/application"
	"myst/internal/client/application/domain/entry"
	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/application/keystoreservice"
	"myst/pkg/config"
	"myst/pkg/logger"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	prometheus "github.com/zsais/go-gin-prometheus"
)

//go:generate oapi-codegen -package generated -generate types -o generated/types.gen.go openapi.json
//go:generate oapi-codegen -package generated -generate client -o generated/client.gen.go openapi.json
// TODO: remove redundant go:generate for old ui
////go:generate openapi-generator-cli generate -i openapi.json -o ../../../../ui/src/api/generated -g typescript-fetch --additional-properties=supportsES6=true,npmVersion=8.1.2,typescriptThreePlus=true
//go:generate openapi-generator-cli generate -i openapi.json -o ../../../../ui/src/api/generated -g typescript-fetch --additional-properties=supportsES6=true,npmVersion=8.1.2,typescriptThreePlus=true

var log = logger.New("router", logger.Cyan)

type API struct {
	*gin.Engine
	app application.Application
}

func (api *API) Authenticate(c *gin.Context) {
	var req generated.AuthenticateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := api.app.Authenticate(req.Password); err == application.ErrAuthenticationFailed {
		c.Status(http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (api *API) CreateKeystore(c *gin.Context) {
	var req generated.CreateKeystoreRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	var k keystore.Keystore
	if req.Password != nil {
		k, err = api.app.CreateFirstKeystore(
			req.Name,
			*req.Password,
		)
		if err != nil {
			log.Error(err)
			Error(c, http.StatusInternalServerError, err)
			return
		}
	} else {
		k, err = api.app.CreateKeystore(req.Name)
		if err != nil {
			log.Error(err)
			Error(c, http.StatusInternalServerError, err)
			return
		}
	}

	entries := []generated.Entry{}

	for _, e := range k.Entries {
		entries = append(entries, generated.Entry{
			Id:       e.Id,
			Website:  e.Website,
			Username: e.Username,
			Password: e.Password,
			Notes:    e.Notes,
		})
	}

	c.JSON(
		http.StatusCreated, generated.Keystore{
			Id:      k.Id,
			Name:    k.Name,
			Entries: entries,
		},
	)
}

//func (api *API) UnlockKeystore(c *gin.Context) {
//	keystoreId := c.Param("keystoreId")
//
//	var req generated.UnlockKeystoreRequest
//
//	err := c.ShouldBindJSON(&req)
//	if err != nil {
//		Error(c, http.StatusBadRequest, err)
//		return
//	}
//
//	k, err := api.app.UnlockKeystore(keystoreId, req.Password)
//	//if errors.Is(err, keystoreservice.ErrAuthenticationFailed) {
//	//	Error(c, http.StatusForbidden, err)
//	//	return
//	//}
//	if err != nil {
//		Error(c, http.StatusInternalServerError, err)
//		return
//	}
//
//	entries := make([]generated.Entry, len(k.Entries()))
//
//	for i, e := range k.Entries() {
//		entries[i] = generated.Entry{
//			Id:       e.Id(),
//			Label:    e.Label(),
//			Username: e.Username(),
//			Password: e.Password(),
//		}
//	}
//
//	Success(
//		c, generated.Keystore{
//			Id:      k.Id(),
//			Name:    k.Name(),
//			Entries: entries,
//		},
//	)
//}

func (api *API) CreateEntry(c *gin.Context) {
	keystoreId := c.Param("keystoreId")

	var req generated.CreateEntryRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	k, err := api.app.Keystore(keystoreId)
	//if errors.Is(err, keystoreservice.ErrAuthenticationRequired) {
	//	Error(c, http.StatusForbidden, err)
	//	return
	//}
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	e, err := api.app.CreateKeystoreEntry(
		k.Id,
		entry.WithWebsite(req.Website),
		entry.WithUsername(req.Username),
		entry.WithPassword(req.Password),
		entry.WithNotes(req.Notes),
	)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	Success(
		c, generated.Entry{
			Id:       e.Id,
			Website:  e.Website,
			Username: e.Username,
			Password: e.Password,
			Notes:    e.Notes,
		},
	)
}

func (api *API) Keystore(c *gin.Context) {
	keystoreId := c.Param("keystoreId")

	k, err := api.app.Keystore(keystoreId)
	//if errors.Is(err, keystoreservice.ErrAuthenticationRequired) {
	//	Error(c, http.StatusForbidden, err)
	//	return
	//} else if errors.Is(err, keystoreservice.ErrAuthenticationFailed) {
	//	Error(c, http.StatusForbidden, err)
	//	return
	//}
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	entries := []generated.Entry{}

	for _, e := range k.Entries {
		entries = append(entries, generated.Entry{
			Id:       e.Id,
			Website:  e.Website,
			Username: e.Username,
			Password: e.Password,
			Notes:    e.Notes,
		})
	}

	Success(
		c, generated.Keystore{
			Id:      k.Id,
			Name:    k.Name,
			Entries: entries,
		},
	)
}

func (api *API) UpdateEntry(c *gin.Context) {
	keystoreId := c.Param("keystoreId")
	entryId := c.Param("entryId")

	var req generated.UpdateEntryRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	e, err := api.app.UpdateKeystoreEntry(keystoreId, entryId, req.Password, req.Notes)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	// TODO: change entries returned to be a map, implemennt the rest
	Success(
		c, generated.Entry{
			Id:       e.Id,
			Website:  e.Website,
			Username: e.Username,
			Password: e.Password,
			Notes:    e.Notes,
		},
	)
}

func (api *API) DeleteEntry(c *gin.Context) {
	keystoreId := c.Param("keystoreId")
	entryId := c.Param("entryId")

	err := api.app.DeleteKeystoreEntry(keystoreId, entryId)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	Success(c, nil)
}

func (api *API) Keystores(c *gin.Context) {
	ks, err := api.app.Keystores()
	if err == keystoreservice.ErrInitializationRequired {
		Success(c, []generated.Keystore{})
		return
	} else if err == keystoreservice.ErrAuthenticationRequired {
		Error(c, http.StatusUnauthorized, err)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	keystores := []generated.Keystore{}

	for _, k := range ks {
		entries := []generated.Entry{}

		for _, e := range k.Entries {
			entries = append(
				entries, generated.Entry{
					Id:       e.Id,
					Website:  e.Website,
					Username: e.Username,
					Password: e.Password,
					Notes:    e.Notes,
				},
			)
		}

		keystores = append(
			keystores, generated.Keystore{
				Id:      k.Id,
				Name:    k.Name,
				Entries: entries,
			},
		)
	}

	Success(c, keystores)
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

	// Attach static serve middleware
	r.Use(static.Serve("/", static.LocalFile("static", false)))

	r.Use(
		cors.New(
			cors.Options{
				AllowOriginFunc: func(origin string) bool {
					return true
				},
				AllowedHeaders: []string{"*"},
				AllowedOrigins: []string{"http://localhost:80", "http://localhost:8082"},
				//// TODO allow more methods (DELETE?)
				AllowedMethods: []string{http.MethodHead, http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
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
