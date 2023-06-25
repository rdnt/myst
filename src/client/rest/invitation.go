package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"myst/src/client/application"
	"myst/src/client/rest/generated"
)

func (s *Server) GetInvitations(c *gin.Context) {
	sid := sessionId(c)

	invs, err := s.app.Invitations(sid)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	restInvs := generated.Invitations{}
	for _, inv := range invs {
		restInv, err := s.invitationToJSON(sid, inv)
		if err != nil {
			log.Error(err)
			Error(c, http.StatusInternalServerError)
			return
		}

		restInvs = append(restInvs, restInv)
	}

	c.JSON(http.StatusOK, restInvs)
}

func (s *Server) CreateInvitation(c *gin.Context) {
	sid := sessionId(c)
	keystoreId := c.Param("keystoreId")

	var req generated.CreateInvitationRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest)
		return
	}

	inv, err := s.app.CreateInvitation(sid, keystoreId, req.Invitee)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	restInv, err := s.invitationToJSON(sid, inv)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, restInv)
}

func (s *Server) AcceptInvitation(c *gin.Context) {
	sid := sessionId(c)
	invitationId := c.Param("invitationId")

	inv, err := s.app.AcceptInvitation(sid, invitationId)
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

	restInv, err := s.invitationToJSON(sid, inv)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restInv)
}

func (s *Server) DeclineOrCancelInvitation(c *gin.Context) {
	sid := sessionId(c)
	invitationId := c.Param("invitationId")

	inv, err := s.app.DeleteInvitation(sid, invitationId)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	restInv, err := s.invitationToJSON(sid, inv)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restInv)
}

func (s *Server) GetInvitation(c *gin.Context) {
	sid := sessionId(c)
	invitationId := c.Param("invitationId")

	inv, err := s.app.Invitation(sid, invitationId)
	if errors.Is(err, application.ErrInvitationNotFound) {
		Error(c, http.StatusNotFound)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	restInv, err := s.invitationToJSON(sid, inv)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restInv)
}

func (s *Server) FinalizeInvitation(c *gin.Context) {
	sid := sessionId(c)
	invitationId := c.Param("invitationId")

	var req generated.FinalizeInvitationRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest)
		return
	}

	inv, err := s.app.FinalizeInvitation(sid, invitationId, req.RemoteKeystoreId, req.InviteePublicKey)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	restInv, err := s.invitationToJSON(sid, inv)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restInv)
}
