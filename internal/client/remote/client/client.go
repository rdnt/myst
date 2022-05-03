package client

import (
	"context"
	"myst/internal/server/api/http/generated"
	"myst/pkg/enclave"
	"net/http"

	"github.com/pkg/errors"
)

var (
	ErrAuthenticationFailed = enclave.ErrAuthenticationFailed
	ErrInvalidResponse      = errors.New("invalid response")
	ErrNotSignedIn          = errors.New("not signed in")
)

type Client interface {
	SignIn(username, password string) error
	SignOut() error
	SignedIn() bool
	UploadKeystore(name string, payload []byte) (*generated.Keystore, error)
	Keystore(id string) (*generated.Keystore, error)
	Keystores() ([]*generated.Keystore, error)
	CreateInvitation(keystoreId, inviteeId string, inviterPublicKey []byte) (*generated.Invitation, error)
	AcceptInvitation(keystoreId, invitationId string, inviteePublicKey []byte) (*generated.Invitation, error)
	FinalizeInvitation(keystoreId, invitationId string, keystoreKey []byte) (*generated.Invitation, error)
	Invitation(keystoreId, invitationId string) (*generated.Invitation, error)
}

type client struct {
	client *generated.ClientWithResponses

	bearerToken string
}

func New() (*client, error) {
	c := &client{}

	cw, err := generated.NewClientWithResponses("http://localhost:8080",
		generated.WithRequestEditorFn(c.authenticate()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create remote client")
	}

	c.client = cw

	return c, nil
}

func (c *client) SignIn(username, password string) error {
	res, err := c.client.LoginWithResponse(
		context.Background(), generated.LoginJSONRequestBody{
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to sign in")
	}

	if res.JSON200 == nil {
		return ErrInvalidResponse
	}

	token := *res.JSON200

	if token == "" {
		return errors.New("invalid token")
	}

	c.bearerToken = string(token)

	return nil
}

func (c *client) SignOut() error {
	c.bearerToken = ""

	return nil
}

func (c *client) SignedIn() bool {
	return c.bearerToken != ""
}

func (c *client) authenticate() generated.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
		return nil
	}
}
