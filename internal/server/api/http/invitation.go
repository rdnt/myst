package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"myst/internal/server/api/http/generated"
)

func (api *API) CreateInvitation(c *gin.Context) {
	//userId := CurrentUser(c)
	keystoreId := c.Param("keystoreId")

	log.Println("SERVER Creating invitation", keystoreId)

	var params generated.CreateInvitationRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}

	inv, err := api.app.CreateInvitation(
		keystoreId,
		params.InviterId,
		params.InviteeId,
	)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	restInv, err := api.ToJSONInvitation(inv)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, restInv)
}

func (api *API) Invitation(c *gin.Context) {
	userId := CurrentUser(c)
	invitationId := c.Param("invitationId")

	inv, err := api.app.UserInvitation(userId, invitationId)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	restInv, err := api.ToJSONInvitation(inv)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restInv)
}

func (api *API) AcceptInvitation(c *gin.Context) {
	userId := CurrentUser(c)
	invitationId := c.Param("invitationId")

	inv, err := api.app.AcceptInvitation(
		userId,
		invitationId,
	)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	restInv, err := api.ToJSONInvitation(inv)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restInv)
}

func (api *API) DeclineOrCancelInvitation(c *gin.Context) {
	userId := CurrentUser(c)
	invitationId := c.Param("invitationId")

	inv, err := api.app.DeclineOrCancelInvitation(
		userId, invitationId,
	)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	restInv, err := api.ToJSONInvitation(inv)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restInv)
}

func (api *API) FinalizeInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	var params generated.FinalizeInvitationRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}

	keystoreKey := params.KeystoreKey

	inv, err := api.app.FinalizeInvitation(
		invitationId,
		keystoreKey,
	)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	restInv, err := api.ToJSONInvitation(inv)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restInv)
}

func (api *API) Invitations(c *gin.Context) {
	userId := CurrentUser(c)

	invs, err := api.app.UserInvitations(userId, nil)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	gen := []generated.Invitation{}

	for _, inv := range invs {
		restInv, err := api.ToJSONInvitation(inv)
		if err != nil {
			log.Error(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		gen = append(gen, restInv)
	}

	c.JSON(http.StatusOK, gen)
}
