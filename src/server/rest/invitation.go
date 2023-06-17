package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"myst/src/server/application"
	"myst/src/server/rest/generated"
)

func (s *Server) CreateInvitation(c *gin.Context) {
	userId := CurrentUser(c)
	keystoreId := c.Param("keystoreId")

	var params generated.CreateInvitationRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		Error(c, http.StatusBadRequest)
		return
	}

	inv, err := s.app.CreateInvitation(
		keystoreId,
		userId,
		params.Invitee,
	)
	if errors.Is(err, application.ErrInviterNotFound) {
		Error(c, http.StatusNotFound, "inviter-not-found")
		return
	} else if errors.Is(err, application.ErrKeystoreNotFound) {
		Error(c, http.StatusNotFound, "keystore-not-found")
		return
	} else if errors.Is(err, application.ErrInviteeNotFound) {
		Error(c, http.StatusNotFound, "invitee-not-found")
		return
	} else if errors.Is(err, application.ErrInvalidInvitee) {
		Error(c, http.StatusBadRequest, "invalid-invitee")
		return
	} else if errors.Is(err, application.ErrForbidden) {
		Error(c, http.StatusForbidden)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	restInv, err := s.invitationToJSON(inv)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, restInv)
}

func (s *Server) Invitation(c *gin.Context) {
	userId := CurrentUser(c)
	invitationId := c.Param("invitationId")

	inv, err := s.app.UserInvitation(userId, invitationId)
	if errors.Is(err, application.ErrInvitationNotFound) {
		Error(c, http.StatusNotFound)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	restInv, err := s.invitationToJSON(inv)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
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
	if errors.Is(err, application.ErrInvitationNotFound) {
		Error(c, http.StatusNotFound)
		return
	} else if errors.Is(err, application.ErrForbidden) {
		Error(c, http.StatusForbidden)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	restInv, err := s.invitationToJSON(inv)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restInv)
}

func (s *Server) DeleteInvitation(c *gin.Context) {
	userId := CurrentUser(c)
	invitationId := c.Param("invitationId")

	inv, err := s.app.DeleteInvitation(
		userId, invitationId,
	)
	if errors.Is(err, application.ErrInvitationNotFound) {
		Error(c, http.StatusNotFound)
		return
	} else if errors.Is(err, application.ErrForbidden) {
		Error(c, http.StatusForbidden)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	restInv, err := s.invitationToJSON(inv)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restInv)
}

func (s *Server) FinalizeInvitation(c *gin.Context) {
	userId := CurrentUser(c)
	invitationId := c.Param("invitationId")

	var req generated.FinalizeInvitationRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest)
		return
	}

	inv, err := s.app.FinalizeInvitation(userId, invitationId, req.KeystoreKey)
	if errors.Is(err, application.ErrInvitationNotFound) {
		Error(c, http.StatusNotFound)
		return
	} else if errors.Is(err, application.ErrForbidden) {
		Error(c, http.StatusForbidden)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	restInv, err := s.invitationToJSON(inv)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restInv)
}

func (s *Server) Invitations(c *gin.Context) {
	userId := CurrentUser(c)

	invs, err := s.app.UserInvitations(userId, application.UserInvitationsOptions{})
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	gen := []generated.Invitation{}

	for _, inv := range invs {
		restInv, err := s.invitationToJSON(inv)
		if err != nil {
			log.Error(err)
			Error(c, http.StatusInternalServerError)
			return
		}

		gen = append(gen, restInv)
	}

	c.JSON(http.StatusOK, gen)
}
