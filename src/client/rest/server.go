package rest

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io/fs"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/static"

	"myst/pkg/config"
	"myst/pkg/logger"
	"myst/src/client/application"
	"myst/src/client/application/domain/enclave"
	"myst/src/client/application/domain/entry"
	"myst/src/client/application/domain/keystore"
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
//go:generate openapi --input openapi.json --output ../../../../ui/src/api/generated --client fetch --useOptions --useUnionTypes

var log = logger.New("router", logger.Cyan)

type Server struct {
	*gin.Engine
	app application.Application
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

	gin.DefaultWriter = ioutil.Discard

	r := gin.New()
	s.Engine = r

	// Do not redirect folders to trailing slash
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	// Custom PrintRouteFunc
	gin.DebugPrintRouteFunc = PrintRoutes

	// always use recovery middleware
	r.Use(gin.CustomRecovery(recoveryHandler))

	// custom logging middleware
	r.Use(LoggerMiddleware)

	// error 404 handling
	r.NoRoute(NoRoute)

	// metrics
	if config.Debug {
		// p := prometheus.NewPrometheus("gin")
		// p.Use(r)
	}

	// error 404 handling
	r.NoRoute(NoRoute)

	// serve the UI
	if config.Debug || ui == nil {
		r.Use(static.Serve("/", static.LocalFile("static", false)))
	} else {
		r.Use(static.Serve("/", EmbedFolder(ui, "static")))
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

	c.JSON(
		http.StatusOK, generated.User{
			Id:       u.Id,
			Username: u.Username,
		},
	)
}

func (s *Server) CreateKeystore(c *gin.Context) {
	var req generated.CreateKeystoreRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	var k keystore.Keystore
	if req.Password != nil {
		k, err = s.app.CreateKeystore(
			keystore.New(keystore.WithName(req.Name)),
		)
		if err != nil {
			log.Error(err)
			Error(c, http.StatusInternalServerError, err)
			return
		}
	} else {
		k, err = s.app.CreateKeystore(keystore.New(keystore.WithName(req.Name)))
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
			Id:       k.Id,
			RemoteId: k.RemoteId,
			Name:     k.Name,
			Entries:  entries,
		},
	)
}

func (s *Server) DeleteKeystore(c *gin.Context) {
	keystoreId := c.Param("keystoreId")

	err := s.app.DeleteKeystore(keystoreId)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) CreateEnclave(c *gin.Context) {
	var req generated.CreateEnclaveRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest, err)
		return
	}

	err = s.app.CreateEnclave(req.Password)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) Enclave(c *gin.Context) {
	err := s.app.Enclave()
	if errors.Is(err, application.ErrInitializationRequired) {
		Error(c, http.StatusNotFound, err)
		return
	} else if errors.Is(err, application.ErrAuthenticationRequired) {
		Error(c, http.StatusUnauthorized, err)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) Import(c *gin.Context) {

	keystoreName := "All"
	b := []byte("")

	csvReader := csv.NewReader(bytes.NewReader(b))

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	k, err := s.app.CreateKeystore(keystore.New(keystore.WithName(keystoreName)))
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
			Id:       k.Id,
			RemoteId: k.RemoteId,
			Name:     k.Name,
			Entries:  entries,
		},
	)
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

	Success(c, nil)
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
		restInvs = append(restInvs, InvitationToRest(inv))
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

	Success(c, InvitationToRest(inv))
}

func (s *Server) AcceptInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	inv, err := s.app.AcceptInvitation(invitationId)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	Success(c, InvitationToRest(inv))
}

func (s *Server) DeclineOrCancelInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	inv, err := s.app.DeclineOrCancelInvitation(invitationId)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	Success(c, InvitationToRest(inv))
}

func (s *Server) FinalizeInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	inv, err := s.app.FinalizeInvitation(invitationId)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError, err)
		return
	}

	Success(c, InvitationToRest(inv))
}

func (s *Server) Keystores(c *gin.Context) {
	ks, err := s.app.Keystores()
	if errors.Is(err, application.ErrInitializationRequired) {
		Error(c, http.StatusNotFound, err)
		return
	} else if errors.Is(err, application.ErrAuthenticationRequired) {
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
				Id:       k.Id,
				RemoteId: k.RemoteId,
				Name:     k.Name,
				Entries:  entries,
			},
		)
	}

	Success(c, keystores)
}

func (s *Server) HealthCheck(_ *gin.Context) {
	s.app.HealthCheck()
}

func (s *Server) Run(addr string) error {
	log.Println("starting app on", addr)

	s.app.Start()

	log.Println("app started")
	return s.Engine.Run(addr)
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
