// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.10.1 DO NOT EDIT.
package generated

// AcceptInvitationRequest defines model for AcceptInvitationRequest.
type AcceptInvitationRequest struct {
	PublicKey []byte `json:"publicKey"`
}

// AuthorizationResponse defines model for AuthorizationResponse.
type AuthorizationResponse struct {
	Token  string `json:"token"`
	UserId string `json:"userId"`
}

// CreateInvitationRequest defines model for CreateInvitationRequest.
type CreateInvitationRequest struct {
	InviteeId string `json:"inviteeId"`
	PublicKey []byte `json:"publicKey"`
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
	CreatedAt    int64  `json:"createdAt"`
	Id           string `json:"id"`
	InviteeId    string `json:"inviteeId"`
	InviteeKey   []byte `json:"inviteeKey"`
	InviterId    string `json:"inviterId"`
	InviterKey   []byte `json:"inviterKey"`
	KeystoreId   string `json:"keystoreId"`
	KeystoreKey  []byte `json:"keystoreKey"`
	KeystoreName string `json:"keystoreName"`
	Status       string `json:"status"`
	UpdatedAt    int64  `json:"updatedAt"`
}

// Keystore defines model for Keystore.
type Keystore struct {
	CreatedAt int64  `json:"createdAt"`
	Id        string `json:"id"`
	Name      string `json:"name"`
	OwnerId   string `json:"ownerId"`
	Payload   []byte `json:"payload"`
	UpdatedAt int64  `json:"updatedAt"`
	Version   int    `json:"version"`
}

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

// User defines model for User.
type User struct {
	CreatedAt int64  `json:"createdAt"`
	Id        string `json:"id"`
	UpdatedAt int64  `json:"updatedAt"`
	Username  string `json:"username"`
}

// LoginJSONBody defines parameters for Login.
type LoginJSONBody LoginRequest

// RegisterJSONBody defines parameters for Register.
type RegisterJSONBody RegisterRequest

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

// RegisterJSONRequestBody defines body for Register for application/json ContentType.
type RegisterJSONRequestBody RegisterJSONBody

// AcceptInvitationJSONRequestBody defines body for AcceptInvitation for application/json ContentType.
type AcceptInvitationJSONRequestBody AcceptInvitationJSONBody

// FinalizeInvitationJSONRequestBody defines body for FinalizeInvitation for application/json ContentType.
type FinalizeInvitationJSONRequestBody FinalizeInvitationJSONBody

// CreateInvitationJSONRequestBody defines body for CreateInvitation for application/json ContentType.
type CreateInvitationJSONRequestBody CreateInvitationJSONBody

// CreateKeystoreJSONRequestBody defines body for CreateKeystore for application/json ContentType.
type CreateKeystoreJSONRequestBody CreateKeystoreJSONBody
