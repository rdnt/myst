package http

//go:generate oapi-codegen -package generated -generate types -o generated/types.gen.go openapi.json
//go:generate oapi-codegen -package generated -generate client -o generated/client.gen.go openapi.json

import (
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	prometheus "github.com/zsais/go-gin-prometheus"

	application "myst/internal/server"
	"myst/internal/server/api/http/generated"
	"myst/pkg/config"
	"myst/pkg/logger"
)

var log = logger.New("router", logger.Cyan)

type API struct {
	*gin.Engine
	app *application.Application
}

func (api *API) CreateKeystore(c *gin.Context) {
	var req generated.CreateKeystoreRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		panic(err)
	}

	payload, err := hex.DecodeString(req.Payload)
	if err != nil {
		panic(err)
	}

	k, err := api.app.CreateKeystore(req.Name, "rdnt", payload)
	if err != nil {
		panic(err)
	}

	c.JSON(
		http.StatusOK, generated.Keystore{
			Id:        k.Id(),
			Name:      k.Name(),
			OwnerId:   k.OwnerId(),
			Payload:   hex.EncodeToString(k.Payload()),
			CreatedAt: int(k.CreatedAt().Unix()),
			UpdatedAt: int(k.UpdatedAt().Unix()),
		},
	)
}

func (api *API) CreateInvitation(c *gin.Context) {
	keystoreId := c.Param("keystoreId")

	//log.Debug(keystoreId)

	var params generated.CreateInvitationRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		panic(err)
	}

	inviterKey, err := hex.DecodeString(params.PublicKey)
	if err != nil {
		panic(err)
	}

	inv, err := api.app.CreateInvitation(
		keystoreId,
		"rdnt",
		params.InviteeId,
		inviterKey,
	)
	if err != nil {
		panic(err)
	}

	c.JSON(
		http.StatusOK, generated.Invitation{
			Id: inv.Id(),
		},
	)
}

func (api *API) GetInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	// TODO: verify client is allowed to accept invitation for that keystore

	inv, err := api.app.GetInvitation(invitationId)
	if err != nil {
		panic(err)
	}

	c.JSON(
		http.StatusOK, generated.Invitation{
			Id: inv.Id(),
		},
	)
}

func (api *API) AcceptInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	// TODO: verify client is allowed to accept invitation for that keystore

	var params generated.AcceptInvitationRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		panic(err)
	}

	inviteeKey, err := hex.DecodeString(params.PublicKey)
	if err != nil {
		panic(err)
	}

	inv, err := api.app.AcceptInvitation(
		invitationId,
		inviteeKey,
	)
	if err != nil {
		panic(err)
	}

	c.JSON(
		http.StatusOK, generated.Invitation{
			Id: inv.Id(),
		},
	)
}

func (api *API) FinalizeInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	// TODO: verify client is allowed to accept invitation for that keystore

	var params generated.FinalizeInvitationRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		panic(err)
	}

	keystoreKey, err := hex.DecodeString(params.KeystoreKey)
	if err != nil {
		panic(err)
	}

	inv, err := api.app.FinalizeInvitation(
		invitationId,
		keystoreKey,
	)
	if err != nil {
		panic(err)
	}

	c.JSON(
		http.StatusOK, generated.Invitation{
			Id: inv.Id(),
		},
	)
}

func (api *API) Login(c *gin.Context) {
	var params generated.LoginRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		panic(err)
	}

	if params.Username == "rdnt" && params.Password != "1234" {
		panic(err)
	}

	key, err := api.loginUser(params.Username)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, generated.AuthToken(key))
}

var jwtCookieLifetime = 604800

//var jwtSecretKey = os.Getenv("JWT_SECRET_KEY")
var jwtSecretKey = "dzl6JEMmRilKQE1jUWZUalduWnI0dTd4IUElRCpHLUs="

func (api *API) loginUser(username string) (string, error) {
	now := time.Now()

	exp := now.Unix() + int64(jwtCookieLifetime)
	iat := now.Unix()
	nbf := now.Unix()

	// if password matches hash
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": exp,
			"iat": iat,
			"nbf": nbf,
			"usr": username,
		},
	)
	key, err := base64.StdEncoding.DecodeString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return token.SignedString(key)
}

func (api *API) Run(addr string) error {
	log.Println("starting app on port :8080")

	api.app.Start()

	log.Println("app started")
	return api.Engine.Run(addr)
}

func New(app *application.Application) *API {
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