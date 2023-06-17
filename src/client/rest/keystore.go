package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"myst/src/client/application"
	"myst/src/client/rest/generated"
)

func (s *Server) CreateKeystore(c *gin.Context) {
	var req generated.CreateKeystoreRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest)
		return
	}

	k, err := s.app.CreateKeystore(req.Name)
	if errors.Is(err, application.ErrInvalidKeystoreName) {
		Error(c, http.StatusBadRequest)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, keystoreToJSON(k))
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
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, keystoreToJSON(k))
}

func (s *Server) DeleteKeystore(c *gin.Context) {
	keystoreId := c.Param("keystoreId")

	err := s.app.DeleteKeystore(keystoreId)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) Keystores(c *gin.Context) {
	ks, err := s.app.Keystores()
	if errors.Is(err, application.ErrInitializationRequired) {
		Error(c, http.StatusUnauthorized)
		return
	} else if errors.Is(err, application.ErrAuthenticationRequired) {
		Error(c, http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	keystores := generated.Keystores{}

	for _, k := range ks {
		keystores = append(
			keystores, keystoreToJSON(k),
		)
	}

	c.JSON(http.StatusOK, keystores)
}

func (s *Server) CreateEntry(c *gin.Context) {
	keystoreId := c.Param("keystoreId")

	var req generated.CreateEntryRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest)
		return
	}

	k, err := s.app.Keystore(keystoreId)
	// if errors.Is(err, keystoreservice.ErrAuthenticationRequired) {
	//	Error(c, rest.StatusForbidden, err)
	//	return
	// }
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	e, err := s.app.CreateEntry(
		k.Id,
		req.Website,
		req.Username,
		req.Password,
		req.Notes,
	)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
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

func (s *Server) UpdateEntry(c *gin.Context) {
	keystoreId := c.Param("keystoreId")
	entryId := c.Param("entryId")

	var req generated.UpdateEntryRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest)
		return
	}

	e, err := s.app.UpdateEntry(keystoreId, entryId,
		application.UpdateEntryOptions{
			Password: req.Password, Notes: req.Notes,
		},
	)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, generated.Entry{
		Id:       e.Id,
		Website:  e.Website,
		Username: e.Username,
		Password: e.Password,
		Notes:    e.Notes,
	})
}

func (s *Server) DeleteEntry(c *gin.Context) {
	keystoreId := c.Param("keystoreId")
	entryId := c.Param("entryId")

	err := s.app.DeleteEntry(keystoreId, entryId)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
