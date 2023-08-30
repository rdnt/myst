// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.4 DO NOT EDIT.
package generated

import (
	"time"
)

// Defines values for InvitationStatus.
const (
	Accepted  InvitationStatus = "accepted"
	Declined  InvitationStatus = "declined"
	Deleted   InvitationStatus = "deleted"
	Finalized InvitationStatus = "finalized"
	Pending   InvitationStatus = "pending"
)

// AuthorizationResponse defines model for AuthorizationResponse.
type AuthorizationResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// CreateInvitationRequest defines model for CreateInvitationRequest.
type CreateInvitationRequest struct {
	Invitee string `json:"invitee"`
}

// CreateKeystoreRequest defines model for CreateKeystoreRequest.
type CreateKeystoreRequest struct {
	Name    string `json:"name"`
	Payload []byte `json:"payload"`
}

// Error defines model for Error.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// FinalizeInvitationRequest defines model for FinalizeInvitationRequest.
type FinalizeInvitationRequest struct {
	KeystoreKey []byte `json:"keystoreKey"`
}

// Invitation defines model for Invitation.
type Invitation struct {
	AcceptedAt           time.Time        `json:"acceptedAt"`
	CancelledAt          time.Time        `json:"cancelledAt"`
	CreatedAt            time.Time        `json:"createdAt"`
	DeclinedAt           time.Time        `json:"declinedAt"`
	DeletedAt            time.Time        `json:"deletedAt"`
	EncryptedKeystoreKey []byte           `json:"encryptedKeystoreKey"`
	FinalizedAt          time.Time        `json:"finalizedAt"`
	Id                   string           `json:"id"`
	Invitee              User             `json:"invitee"`
	Inviter              User             `json:"inviter"`
	Keystore             KeystoreName     `json:"keystore"`
	Status               InvitationStatus `json:"status"`
	UpdatedAt            time.Time        `json:"updatedAt"`
}

// InvitationStatus defines model for Invitation.Status.
type InvitationStatus string

// Keystore defines model for Keystore.
type Keystore struct {
	CreatedAt time.Time `json:"createdAt"`
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	OwnerId   string    `json:"ownerId"`
	Payload   []byte    `json:"payload"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// KeystoreName defines model for KeystoreName.
type KeystoreName struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// LoginRequest defines model for LoginRequest.
type LoginRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// RegisterRequest defines model for RegisterRequest.
type RegisterRequest struct {
	Password  string `json:"password"`
	PublicKey []byte `json:"publicKey"`
	Username  string `json:"username"`
}

// UpdateKeystoreRequest defines model for UpdateKeystoreRequest.
type UpdateKeystoreRequest struct {
	Name    *string `json:"name,omitempty"`
	Payload *[]byte `json:"payload,omitempty"`
}

// User defines model for User.
type User struct {
	Id        string `json:"id"`
	PublicKey []byte `json:"publicKey"`
	Username  string `json:"username"`
}

// UserByUsernameParams defines parameters for UserByUsername.
type UserByUsernameParams struct {
	Username string `form:"username" json:"username"`
}

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody = LoginRequest

// RegisterJSONRequestBody defines body for Register for application/json ContentType.
type RegisterJSONRequestBody = RegisterRequest

// FinalizeInvitationJSONRequestBody defines body for FinalizeInvitation for application/json ContentType.
type FinalizeInvitationJSONRequestBody = FinalizeInvitationRequest

// UpdateKeystoreJSONRequestBody defines body for UpdateKeystore for application/json ContentType.
type UpdateKeystoreJSONRequestBody = UpdateKeystoreRequest

// CreateInvitationJSONRequestBody defines body for CreateInvitation for application/json ContentType.
type CreateInvitationJSONRequestBody = CreateInvitationRequest

// CreateKeystoreJSONRequestBody defines body for CreateKeystore for application/json ContentType.
type CreateKeystoreJSONRequestBody = CreateKeystoreRequest
