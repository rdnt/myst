package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"myst/internal/server/api/http/generated"
)

func (api *API) CreateInvitation(c *gin.Context) {
	userId := CurrentUser(c)
	keystoreId := c.Param("keystoreId")

	var params generated.CreateInvitationRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		panic(err)
	}

	inviterKey := params.PublicKey

	inv, err := api.app.CreateInvitation(
		keystoreId,
		userId,
		params.InviteeId,
		inviterKey,
	)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, ToJSONInvitation(inv))
}

func (api *API) Invitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	// TODO: verify client is allowed to accept invitation for that keystore

	inv, err := api.app.GetInvitation(invitationId)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, ToJSONInvitation(inv))
}

func (api *API) AcceptInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	// TODO: verify client is allowed to accept invitation for that keystore

	var params generated.AcceptInvitationRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}

	inviteeKey := params.PublicKey

	inv, err := api.app.AcceptInvitation(
		invitationId,
		inviteeKey,
	)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, ToJSONInvitation(inv))
}

func (api *API) FinalizeInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	// TODO: verify client is allowed to accept invitation for that keystore

	var params generated.FinalizeInvitationRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		panic(err)
	}

	keystoreKey := params.KeystoreKey

	inv, err := api.app.FinalizeInvitation(
		invitationId,
		keystoreKey,
	)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, ToJSONInvitation(inv))
}

func (api *API) Invitations(c *gin.Context) {
	userId := CurrentUser(c)

	invs, err := api.app.UserInvitations(userId)
	if err != nil {
		panic(err)
	}

	gen := []generated.Invitation{}

	for _, inv := range invs {
		gen = append(gen, ToJSONInvitation(inv))
	}

	c.JSON(http.StatusOK, gen)
}
