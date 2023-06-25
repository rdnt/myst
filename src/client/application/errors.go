package application

import "github.com/pkg/errors"

var (
	ErrKeystoreNotFound       = errors.New("keystore not found")
	ErrAuthenticationRequired = errors.New("authentication required")
	ErrAuthenticationFailed   = errors.New("authentication failed")
	ErrInitializationRequired = errors.New("initialization required")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrInvalidWebsite         = errors.New("invalid website")
	ErrorInvalidUsername      = errors.New("invalid username")
	ErrEntryNotFound          = errors.New("entry not found")
	ErrInvalidKeystoreName    = errors.New("invalid keystore name")
	ErrCredentialsNotFound    = errors.New("credentials not found")
	ErrEnclaveExists          = errors.New("enclave already exists")
	ErrInvitationNotFound     = errors.New("invitation not found")
	ErrForbidden              = errors.New("forbidden")
	ErrSharedSecretNotFound   = errors.New("shared secret not found")
	ErrRemoteAddressMismatch  = errors.New("remote address mismatch")
)
