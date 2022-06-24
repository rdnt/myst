package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"myst/internal/client/api/http/generated"
	"myst/internal/client/application"
)

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

func (api *API) Login(c *gin.Context) {
	var req generated.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := api.app.SignIn(req.Username, req.Password); err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (api *API) Register(c *gin.Context) {
	var req generated.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	u, err := api.app.Register(req.Username, req.Password)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, UserToRest(u))
}
