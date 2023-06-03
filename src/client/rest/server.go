package rest

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"

	"myst/src/client/enclaverepo/enclave"

	"github.com/gin-contrib/static"

	"myst/pkg/config"
	"myst/pkg/logger"
	"myst/pkg/server"
	"myst/src/client/application"
	"myst/src/client/application/domain/keystore/entry"
	"myst/src/client/rest/generated"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	// prometheus "github.com/zsais/go-gin-prometheus"
)

//go:generate oapi-codegen --config oapi-codegen-models.yaml openapi.json
//go:generate oapi-codegen --config oapi-codegen-client.yaml openapi.json
// TODO: remove redundant go:generate for old ui
// //go:generate openapi-generator-cli generate -i openapi.json -o ../../../../ui/src/api/generated -g typescript-fetch --additional-properties=supportsES6=true,npmVersion=8.1.2,typescriptThreePlus=true
// //go:generate openapi-generator-cli generate -i openapi.json -o ../../../../ui/src/api/generated -g typescript-fetch --additional-properties=supportsES6=true,npmVersion=8.1.2,typescriptThreePlus=true,withInterfaces=true
//go:generate npx openapi-typescript-codegen --input openapi.json --output ../../../ui/src/api/generated --client fetch --useOptions --useUnionTypes

//go:generate oapi-codegen --config oapi-codegen-models.yaml openapi.json

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
	gin.DebugPrintRouteFunc = PrintRoutes

	// always use recovery middleware
	// r.Use(gin.CustomRecovery(recoveryHandler))

	// custom logging middleware
	r.Use(LoggerMiddleware)

	// error 404 handling
	r.NoRoute(NoRoute("/", EmbedFolder(ui, "static")))

	// metrics
	if config.Debug {
		// p := prometheus.NewPrometheus("gin")
		// p.Use(r)
	}

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

	s.initRoutes(r.Group("api"))

	return s
}

func (s *Server) CurrentUser(c *gin.Context) {
	u, err := s.app.CurrentUser()
	if errors.Is(err, enclave.ErrRemoteNotSet) {
		c.Status(http.StatusNotFound)
		return
	} else if u == nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	restUser, err := s.userToRest(*u)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, restUser)
}

func (s *Server) CreateKeystore(c *gin.Context) {
	var req generated.CreateKeystoreRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	k, err := s.app.CreateKeystore(req.Name)
	if errors.Is(err, application.ErrInvalidKeystoreName) {
		Error(c, http.StatusBadRequest, err)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, KeystoreToRest(k))
}

func (s *Server) DeleteKeystore(c *gin.Context) {
	keystoreId := c.Param("keystoreId")

	err := s.app.DeleteKeystore(keystoreId)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) CreateEnclave(c *gin.Context) {
	var req generated.CreateEnclaveRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	err = s.app.Initialize(req.Password)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusCreated)
}

func (s *Server) Enclave(c *gin.Context) {
	exists, err := s.app.IsInitialized()
	if errors.Is(err, application.ErrAuthenticationRequired) {
		Error(c, http.StatusUnauthorized, err)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	if !exists {
		Error(c, http.StatusNotFound, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) Import(c *gin.Context) {

	keystoreName := "Passwords"
	b := []byte("")

	csvReader := csv.NewReader(bytes.NewReader(b))

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	k, err := s.app.CreateKeystore(keystoreName)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	for _, row := range records {
		if len(row) != 4 {
			continue
		}

		website := row[0]
		// url := row[1]
		username := row[2]
		password := row[3]

		_, err := s.app.CreateKeystoreEntry(
			k.Id,
			entry.WithWebsite(website),
			entry.WithUsername(username),
			entry.WithPassword(password),
		)
		if err != nil {
			log.Error(err)
			Error(c, http.StatusInternalServerError, err)
			return
		}
		fmt.Println("Imported: ", website, username, "***")

	}

	c.JSON(http.StatusOK, nil)
}

// func (api *Server) UnlockKeystore(c *gin.Context) {
//	keystoreId := c.Param("keystoreId")
//
//	var req generated.UnlockKeystoreRequest
//
//	err := c.ShouldBindJSON(&req)
//	if err != nil {
//		Error(c, rest.StatusBadRequest, err)
//		return
//	}
//
//	k, err := api.app.UnlockKeystore(keystoreId, req.Password)
//	//if errors.Is(err, keystoreservice.ErrAuthenticationFailed) {
//	//	Error(c, rest.StatusForbidden, err)
//	//	return
//	//}
//	if err != nil {
//		Error(c, rest.StatusInternalServerError, err)
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
// }

func (s *Server) CreateEntry(c *gin.Context) {
	keystoreId := c.Param("keystoreId")

	var req generated.CreateEntryRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	k, err := s.app.Keystore(keystoreId)
	// if errors.Is(err, keystoreservice.ErrAuthenticationRequired) {
	//	Error(c, rest.StatusForbidden, err)
	//	return
	// }
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	e, err := s.app.CreateKeystoreEntry(
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

	c.JSON(http.StatusCreated, generated.Entry{
		Id:       e.Id,
		Website:  e.Website,
		Username: e.Username,
		Password: e.Password,
		Notes:    e.Notes,
	})
}

func (s *Server) Keystore(c *gin.Context) {
	keystoreId := c.Param("keystoreId")

	k, err := s.app.Keystore(keystoreId)
	// if errors.Is(err, keystoreservice.ErrAuthenticationRequired) {
	//	Error(c, rest.StatusForbidden, err)
	//	return
	// } else if errors.Is(err, keystoreservice.ErrAuthenticationFailed) {
	//	Error(c, rest.StatusForbidden, err)
	//	return
	// }
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, KeystoreToRest(k))
}

func (s *Server) UpdateEntry(c *gin.Context) {
	keystoreId := c.Param("keystoreId")
	entryId := c.Param("entryId")

	var req generated.UpdateEntryRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	e, err := s.app.UpdateKeystoreEntry(keystoreId, entryId, req.Password, req.Notes)
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

func (s *Server) DeleteEntry(c *gin.Context) {
	keystoreId := c.Param("keystoreId")
	entryId := c.Param("entryId")

	err := s.app.DeleteKeystoreEntry(keystoreId, entryId)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) GetInvitations(c *gin.Context) {
	invs, err := s.app.Invitations()
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	restInvs := generated.Invitations{}
	for _, inv := range invs {
		restInv, err := s.InvitationToRest(inv)
		if err != nil {

			log.Error(err)
			Error(c, http.StatusInternalServerError, err)
			return
		}

		restInvs = append(restInvs, restInv)
	}

	Success(c, restInvs)
}

func (s *Server) CreateInvitation(c *gin.Context) {
	keystoreId := c.Param("keystoreId")

	var req generated.CreateInvitationRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	inv, err := s.app.CreateInvitation(keystoreId, req.Invitee)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	restInv, err := s.InvitationToRest(inv)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, restInv)
}

func (s *Server) AcceptInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	inv, err := s.app.AcceptInvitation(invitationId)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	restInv, err := s.InvitationToRest(inv)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	Success(c, restInv)
}

func (s *Server) DeclineOrCancelInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	inv, err := s.app.DeclineOrCancelInvitation(invitationId)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	restInv, err := s.InvitationToRest(inv)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	Success(c, restInv)
}

func (s *Server) GetInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	inv, err := s.app.Invitation(invitationId)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	restInv, err := s.InvitationToRest(inv)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	Success(c, restInv)
}

func (s *Server) FinalizeInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	var req generated.FinalizeInvitationRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	inv, err := s.app.FinalizeInvitation(invitationId, req.RemoteKeystoreId, req.InviteePublicKey)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	restInv, err := s.InvitationToRest(inv)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	Success(c, restInv)
}

func (s *Server) Keystores(c *gin.Context) {
	ks, err := s.app.Keystores()
	if errors.Is(err, application.ErrInitializationRequired) {
		Error(c, http.StatusUnauthorized, err)
		return
	} else if errors.Is(err, application.ErrAuthenticationRequired) {
		Error(c, http.StatusUnauthorized, err)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	keystores := generated.Keystores{}

	for _, k := range ks {
		keystores = append(
			keystores, KeystoreToRest(k),
		)
	}

	Success(c, keystores)
}

func (s *Server) HealthCheck(_ *gin.Context) {
	s.app.HealthCheck()
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
	// TODO: find way to return all errors, maybe (find) an errgroup-esque package
	// that just runs all functions sequentially and returns all errors as one?
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

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	if err != nil {
		return false
	}
	return true
}

func EmbedFolder(fsEmbed fs.FS, targetPath string) static.ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}
