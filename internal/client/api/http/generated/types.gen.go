// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.1 DO NOT EDIT.
package generated

// Defines values for InvitationStatus.
const (
	InvitationStatusAccepted InvitationStatus = "accepted"

	InvitationStatusFinalized InvitationStatus = "finalized"

	InvitationStatusPending InvitationStatus = "pending"

	InvitationStatusRejected InvitationStatus = "rejected"
)

// AuthenticateRequest defines model for AuthenticateRequest.
type AuthenticateRequest struct {
	Password string `json:"password"`
}

// CreateEntryRequest defines model for CreateEntryRequest.
type CreateEntryRequest struct {
	Notes    string `json:"notes"`
	Password string `json:"password"`
	Username string `json:"username"`
	Website  string `json:"website"`
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
	CreatedAt   int64            `json:"createdAt"`
	Id          string           `json:"id"`
	InviteeId   string           `json:"inviteeId"`
	InviteeKey  []byte           `json:"inviteeKey"`
	InviterId   string           `json:"inviterId"`
	InviterKey  []byte           `json:"inviterKey"`
	KeystoreId  string           `json:"keystoreId"`
	KeystoreKey []byte           `json:"keystoreKey"`
	Status      InvitationStatus `json:"status"`
	UpdatedAt   int64            `json:"updatedAt"`
}

// InvitationStatus defines model for Invitation.Status.
type InvitationStatus string

// Invitations defines model for Invitations.
type Invitations []Invitation

// Keystore defines model for Keystore.
type Keystore struct {
	Entries []Entry `json:"entries"`
	Id      string  `json:"id"`
	Name    string  `json:"name"`
}

// Keystores defines model for Keystores.
type Keystores []Keystore

// UpdateEntryRequest defines model for UpdateEntryRequest.
type UpdateEntryRequest struct {
	Notes    *string `json:"notes,omitempty"`
	Password *string `json:"password,omitempty"`
}

// AuthenticateJSONBody defines parameters for Authenticate.
type AuthenticateJSONBody AuthenticateRequest

// CreateEntryJSONBody defines parameters for CreateEntry.
type CreateEntryJSONBody CreateEntryRequest

// UpdateEntryJSONBody defines parameters for UpdateEntry.
type UpdateEntryJSONBody UpdateEntryRequest

// CreateKeystoreJSONBody defines parameters for CreateKeystore.
type CreateKeystoreJSONBody CreateKeystoreRequest

// AuthenticateJSONRequestBody defines body for Authenticate for application/json ContentType.
type AuthenticateJSONRequestBody AuthenticateJSONBody

// CreateEntryJSONRequestBody defines body for CreateEntry for application/json ContentType.
type CreateEntryJSONRequestBody CreateEntryJSONBody

// UpdateEntryJSONRequestBody defines body for UpdateEntry for application/json ContentType.
type UpdateEntryJSONRequestBody UpdateEntryJSONBody

// CreateKeystoreJSONRequestBody defines body for CreateKeystore for application/json ContentType.
type CreateKeystoreJSONRequestBody CreateKeystoreJSONBody
