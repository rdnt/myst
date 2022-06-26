// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.1 DO NOT EDIT.
package generated

import (
	"time"
)

// Defines values for InvitationStatus.
const (
	InvitationStatusAccepted InvitationStatus = "accepted"

	InvitationStatusDeclined InvitationStatus = "declined"

	InvitationStatusDeleted InvitationStatus = "deleted"

	InvitationStatusFinalized InvitationStatus = "finalized"

	InvitationStatusPending InvitationStatus = "pending"
)

// AuthenticateRequest defines model for AuthenticateRequest.
type AuthenticateRequest struct {
	Password string `json:"password"`
}

// CreateEnclaveRequest defines model for CreateEnclaveRequest.
type CreateEnclaveRequest struct {
	Password string `json:"password"`
}

// CreateEntryRequest defines model for CreateEntryRequest.
type CreateEntryRequest struct {
	Notes    string `json:"notes"`
	Password string `json:"password"`
	Username string `json:"username"`
	Website  string `json:"website"`
}

// CreateInvitationRequest defines model for CreateInvitationRequest.
type CreateInvitationRequest struct {
	InviteeId string `json:"inviteeId"`
}

// CreateKeystoreRequest defines model for CreateKeystoreRequest.
type CreateKeystoreRequest struct {
	Name     string  `json:"name"`
	Password *string `json:"password,omitempty"`
}

// Entry defines model for Entry.
type Entry struct {
	Id       string `json:"id"`
	Notes    string `json:"notes"`
	Password string `json:"password"`
	Username string `json:"username"`
	Website  string `json:"website"`
}

// Error defines model for Error.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Invitation defines model for Invitation.
type Invitation struct {
	AcceptedAt   time.Time        `json:"acceptedAt"`
	CreatedAt    time.Time        `json:"createdAt"`
	DeclinedAt   time.Time        `json:"declinedAt"`
	DeletedAt    time.Time        `json:"deletedAt"`
	Id           string           `json:"id"`
	InviteeId    string           `json:"inviteeId"`
	InviteeKey   []byte           `json:"inviteeKey"`
	InviterId    string           `json:"inviterId"`
	InviterKey   []byte           `json:"inviterKey"`
	KeystoreId   string           `json:"keystoreId"`
	KeystoreKey  []byte           `json:"keystoreKey"`
	KeystoreName string           `json:"keystoreName"`
	Status       InvitationStatus `json:"status"`
	UpdatedAt    time.Time        `json:"updatedAt"`
}

// InvitationStatus defines model for Invitation.Status.
type InvitationStatus string

// Invitations defines model for Invitations.
type Invitations []Invitation

// Keystore defines model for Keystore.
type Keystore struct {
	Entries  []Entry `json:"entries"`
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	RemoteId string  `json:"remoteId"`
}

// Keystores defines model for Keystores.
type Keystores []Keystore

// LoginRequest defines model for LoginRequest.
type LoginRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// RegisterRequest defines model for RegisterRequest.
type RegisterRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// UpdateEntryRequest defines model for UpdateEntryRequest.
type UpdateEntryRequest struct {
	Notes    *string `json:"notes,omitempty"`
	Password *string `json:"password,omitempty"`
}

// User defines model for User.
type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

// LoginJSONBody defines parameters for Login.
type LoginJSONBody LoginRequest

// RegisterJSONBody defines parameters for Register.
type RegisterJSONBody RegisterRequest

// AuthenticateJSONBody defines parameters for Authenticate.
type AuthenticateJSONBody AuthenticateRequest

// CreateEnclaveJSONBody defines parameters for CreateEnclave.
type CreateEnclaveJSONBody CreateEnclaveRequest

// CreateEntryJSONBody defines parameters for CreateEntry.
type CreateEntryJSONBody CreateEntryRequest

// UpdateEntryJSONBody defines parameters for UpdateEntry.
type UpdateEntryJSONBody UpdateEntryRequest

// CreateInvitationJSONBody defines parameters for CreateInvitation.
type CreateInvitationJSONBody CreateInvitationRequest

// CreateKeystoreJSONBody defines parameters for CreateKeystore.
type CreateKeystoreJSONBody CreateKeystoreRequest

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody LoginJSONBody

// RegisterJSONRequestBody defines body for Register for application/json ContentType.
type RegisterJSONRequestBody RegisterJSONBody

// AuthenticateJSONRequestBody defines body for Authenticate for application/json ContentType.
type AuthenticateJSONRequestBody AuthenticateJSONBody

// CreateEnclaveJSONRequestBody defines body for CreateEnclave for application/json ContentType.
type CreateEnclaveJSONRequestBody CreateEnclaveJSONBody

// CreateEntryJSONRequestBody defines body for CreateEntry for application/json ContentType.
type CreateEntryJSONRequestBody CreateEntryJSONBody

// UpdateEntryJSONRequestBody defines body for UpdateEntry for application/json ContentType.
type UpdateEntryJSONRequestBody UpdateEntryJSONBody

// CreateInvitationJSONRequestBody defines body for CreateInvitation for application/json ContentType.
type CreateInvitationJSONRequestBody CreateInvitationJSONBody

// CreateKeystoreJSONRequestBody defines body for CreateKeystore for application/json ContentType.
type CreateKeystoreJSONRequestBody CreateKeystoreJSONBody
