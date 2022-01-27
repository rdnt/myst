// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version (devel) DO NOT EDIT.
package generated

// AuthenticateRequest defines model for AuthenticateRequest.
type AuthenticateRequest struct {
	Password string `json:"password"`
}

// CreateEntryRequest defines model for CreateEntryRequest.
type CreateEntryRequest struct {
	Label    string `json:"label"`
	Password string `json:"password"`
	Username string `json:"username"`
}

// CreateKeystoreRequest defines model for CreateKeystoreRequest.
type CreateKeystoreRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Entry defines model for Entry.
type Entry struct {
	Id       string `json:"id"`
	Label    string `json:"label"`
	Password string `json:"password"`
	Username string `json:"username"`
}

// Error defines model for Error.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Keystore defines model for Keystore.
type Keystore struct {
	Entries []Entry `json:"entries"`
	Id      string  `json:"id"`
	Name    string  `json:"name"`
}

// UnlockKeystoreRequest defines model for UnlockKeystoreRequest.
type UnlockKeystoreRequest struct {
	Password string `json:"password"`
}

// AuthenticateJSONBody defines parameters for Authenticate.
type AuthenticateJSONBody AuthenticateRequest

// UnlockKeystoreJSONBody defines parameters for UnlockKeystore.
type UnlockKeystoreJSONBody UnlockKeystoreRequest

// CreateEntryJSONBody defines parameters for CreateEntry.
type CreateEntryJSONBody CreateEntryRequest

// CreateKeystoreJSONBody defines parameters for CreateKeystore.
type CreateKeystoreJSONBody CreateKeystoreRequest

// AuthenticateJSONRequestBody defines body for Authenticate for application/json ContentType.
type AuthenticateJSONRequestBody AuthenticateJSONBody

// UnlockKeystoreJSONRequestBody defines body for UnlockKeystore for application/json ContentType.
type UnlockKeystoreJSONRequestBody UnlockKeystoreJSONBody

// CreateEntryJSONRequestBody defines body for CreateEntry for application/json ContentType.
type CreateEntryJSONRequestBody CreateEntryJSONBody

// CreateKeystoreJSONRequestBody defines body for CreateKeystore for application/json ContentType.
type CreateKeystoreJSONRequestBody CreateKeystoreJSONBody
