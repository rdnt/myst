package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"myst/src/server/rest/generated"
)

func (s *Server) CreateInvitation(c *gin.Context) {
	// userId := CurrentUser(c)
	keystoreId := c.Param("keystoreId")

	log.Println("SERVER Creating invitation", keystoreId)

	var params generated.CreateInvitationRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}

	inv, err := s.app.CreateInvitation(
		keystoreId,
		params.InviterId,
		params.InviteeId,
	)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	restInv, err := s.ToJSONInvitation(inv)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, restInv)
}

func (s *Server) Invitation(c *gin.Context) {
	userId := CurrentUser(c)
	invitationId := c.Param("invitationId")

	inv, err := s.app.UserInvitation(userId, invitationId)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	restInv, err := s.ToJSONInvitation(inv)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restInv)
}

func (s *Server) AcceptInvitation(c *gin.Context) {
	userId := CurrentUser(c)
	invitationId := c.Param("invitationId")

	inv, err := s.app.AcceptInvitation(
		userId,
		invitationId,
	)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	restInv, err := s.ToJSONInvitation(inv)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restInv)
}

func (s *Server) DeclineOrCancelInvitation(c *gin.Context) {
	userId := CurrentUser(c)
	invitationId := c.Param("invitationId")

	inv, err := s.app.DeclineOrCancelInvitation(
		userId, invitationId,
	)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	restInv, err := s.ToJSONInvitation(inv)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restInv)
}

func (s *Server) FinalizeInvitation(c *gin.Context) {
	invitationId := c.Param("invitationId")

	var params generated.FinalizeInvitationRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}

	keystoreKey := params.KeystoreKey

	inv, err := s.app.FinalizeInvitation(
		invitationId,
		keystoreKey,
	)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	restInv, err := s.ToJSONInvitation(inv)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restInv)
}

func (s *Server) Invitations(c *gin.Context) {
	userId := CurrentUser(c)

	invs, err := s.app.UserInvitations(userId, nil)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	gen := []generated.Invitation{}

	for _, inv := range invs {
		restInv, err := s.ToJSONInvitation(inv)
		if err != nil {
			log.Error(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		gen = append(gen, restInv)
	}

	c.JSON(http.StatusOK, gen)
}
