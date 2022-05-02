package remote

import (
	"context"
	"fmt"
	"myst/internal/client/application/domain/keystore"
	"myst/pkg/crypto"
	"net/http"

	"github.com/pkg/errors"

	"golang.org/x/crypto/curve25519"

	"myst/internal/server/api/http/generated"
	"myst/pkg/enclave"
)

var (
	ErrAuthenticationFailed = enclave.ErrAuthenticationFailed
	ErrInvalidResponse      = fmt.Errorf("invalid response")
)

type Client interface {
	SignIn(username, password string) error
	SignOut() error
	CreateKeystore(name string, keystoreKey []byte, keystore *keystore.Keystore) (*generated.Keystore, error)
	Keystore(id string) (*generated.Keystore, error)
	Keystores() ([]*generated.Keystore, error)
	CreateInvitation(keystoreId, inviteeId string) (*generated.Invitation, error)
	AcceptInvitation(keystoreId, invitationId string) (*generated.Invitation, error)
	FinalizeInvitation(keystoreId, invitationId string, keystoreKey []byte) (*generated.Invitation, error)
}

type remote struct {
	client *generated.ClientWithResponses

	bearerToken string

	publicKey  []byte
	privateKey []byte
}

func New() (Client, error) {
	r := &remote{
		client:      nil,
		bearerToken: "",
		publicKey:   nil,
		privateKey:  nil,
	}

	c, err := generated.NewClientWithResponses("http://localhost:8080",
		generated.WithRequestEditorFn(r.authenticate()),
	)
	if err != nil {
		return nil, err
	}

	pub, key, err := newKeypair()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate public/private keypair")
	}

	r.client = c
	r.publicKey = pub
	r.privateKey = key

	return r, nil
}

func (r *remote) SignIn(username, password string) error {
	fmt.Println("Signing in to remote...")

	res, err := r.client.LoginWithResponse(
		context.Background(), generated.LoginJSONRequestBody{
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		return err
	}

	if res.JSON200 == nil {
		return ErrInvalidResponse
	}

	token := *res.JSON200

	if token == "" {
		return fmt.Errorf("invalid token")
	}

	r.bearerToken = string(token)
	fmt.Println("Signed in.")

	return nil
}

func (r *remote) SignOut() error {
	r.bearerToken = ""

	return nil
}

func (r *remote) authenticate() generated.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+r.bearerToken)
		return nil
	}
}

func newKeypair() ([]byte, []byte, error) {
	var pub [32]byte
	var key [32]byte

	b, err := crypto.GenerateRandomBytes(32)
	if err != nil {
		return nil, nil, err
	}
	copy(key[:], b)

	curve25519.ScalarBaseMult(&pub, &key)

	return pub[:], key[:], nil
}
