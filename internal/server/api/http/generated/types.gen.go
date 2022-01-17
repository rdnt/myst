// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version (devel) DO NOT EDIT.
package generated

// AcceptInvitationRequest defines model for AcceptInvitationRequest.
type AcceptInvitationRequest struct {
	PublicKey string `json:"publicKey"`
}

// AuthToken defines model for AuthToken.
type AuthToken string

// CreateInvitationRequest defines model for CreateInvitationRequest.
type CreateInvitationRequest struct {
	InviteeId string `json:"inviteeId"`
	PublicKey string `json:"publicKey"`
}

// CreateKeystoreRequest defines model for CreateKeystoreRequest.
type CreateKeystoreRequest struct {
	Name    string `json:"name"`
	Payload string `json:"payload"`
}

// Error defines model for Error.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// FinalizeInvitationRequest defines model for FinalizeInvitationRequest.
type FinalizeInvitationRequest struct {
	KeystoreKey string `json:"keystoreKey"`
}

// Invitation defines model for Invitation.
type Invitation struct {
	Id string `json:"id"`
}

// Keystore defines model for Keystore.
type Keystore struct {
	CreatedAt int    `json:"createdAt"`
	Id        string `json:"id"`
	Name      string `json:"name"`
	OwnerId   string `json:"ownerId"`
	Payload   string `json:"payload"`
	UpdatedAt int    `json:"updatedAt"`
}

// LoginRequest defines model for LoginRequest.
type LoginRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// LoginJSONBody defines parameters for Login.
type LoginJSONBody LoginRequest

// AcceptInvitationJSONBody defines parameters for AcceptInvitation.
type AcceptInvitationJSONBody AcceptInvitationRequest

// FinalizeInvitationJSONBody defines parameters for FinalizeInvitation.
type FinalizeInvitationJSONBody FinalizeInvitationRequest

// CreateInvitationJSONBody defines parameters for CreateInvitation.
type CreateInvitationJSONBody CreateInvitationRequest

// CreateKeystoreJSONBody defines parameters for CreateKeystore.
type CreateKeystoreJSONBody CreateKeystoreRequest

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody LoginJSONBody

// AcceptInvitationJSONRequestBody defines body for AcceptInvitation for application/json ContentType.
type AcceptInvitationJSONRequestBody AcceptInvitationJSONBody

// FinalizeInvitationJSONRequestBody defines body for FinalizeInvitation for application/json ContentType.
type FinalizeInvitationJSONRequestBody FinalizeInvitationJSONBody

// CreateInvitationJSONRequestBody defines body for CreateInvitation for application/json ContentType.
type CreateInvitationJSONRequestBody CreateInvitationJSONBody

// CreateKeystoreJSONRequestBody defines body for CreateKeystore for application/json ContentType.
type CreateKeystoreJSONRequestBody CreateKeystoreJSONBody
