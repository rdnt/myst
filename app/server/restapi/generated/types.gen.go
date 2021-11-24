// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version (devel) DO NOT EDIT.
package generated

// CreateKeystoreInvitationRequest defines model for CreateKeystoreInvitationRequest.
type CreateKeystoreInvitationRequest struct {
	InviteeId string `json:"inviteeId"`
	PublicKey string `json:"publicKey"`
}

// Error defines model for Error.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Invitation defines model for Invitation.
type Invitation struct {
	Id string `json:"id"`
}

// CreateKeystoreInvitationJSONBody defines parameters for CreateKeystoreInvitation.
type CreateKeystoreInvitationJSONBody CreateKeystoreInvitationRequest

// CreateKeystoreInvitationJSONRequestBody defines body for CreateKeystoreInvitation for application/json ContentType.
type CreateKeystoreInvitationJSONRequestBody CreateKeystoreInvitationJSONBody