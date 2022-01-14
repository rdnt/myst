package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"myst/internal/client/core/domain/invitation"

	"myst/internal/server/api/http/generated"

	"myst/internal/client/core/domain/keystore"
	"myst/pkg/enclave"
)

var (
	ErrAuthenticationFailed = enclave.ErrAuthenticationFailed
)

type Client interface {
	SignIn(username, password string) error
	SignOut() error
	Keystores() ([]*keystore.Keystore, error)
	CreateKeystoreInvitation(keystoreId, inviteeId, publicKey string) (*invitation.Invitation, error)
}

type remote struct {
	client *generated.Client

	bearerToken string
}

func (r *remote) SignIn(username, password string) error {
	fmt.Println("SignIn", username, password)

	resp, err := r.client.Login(
		context.Background(), generated.LoginJSONRequestBody{
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var token generated.AuthToken
	err = json.Unmarshal(b, &token)
	if err != nil {
		return err
	}

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

func (r *remote) Keystores() ([]*keystore.Keystore, error) {
	return nil, nil
}

func (r *remote) CreateKeystoreInvitation(keystoreId, inviteeId, publicKey string) (*invitation.Invitation, error) {
	fmt.Println("CreateKeystoreInvitation", keystoreId, inviteeId, publicKey)

	if r.bearerToken == "" {
		return nil, fmt.Errorf("not signed in")
	}

	resp, err := r.client.CreateKeystoreInvitation(
		context.Background(), keystoreId, generated.CreateKeystoreInvitationJSONRequestBody{
			InviteeId: inviteeId,
			PublicKey: publicKey,
		},
		r.authenticate(),
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jinv generated.Invitation
	err = json.Unmarshal(b, &jinv)
	if err != nil {
		return nil, err
	}

	fmt.Println("Invitation created", jinv)

	inv, err := invitation.New(
		invitation.WithId(jinv.Id),
	)
	if err != nil {
		return nil, err
	}

	return inv, nil
}

func New() (*remote, error) {
	c, err := generated.NewClient("http://localhost:8080")
	if err != nil {
		return nil, err
	}

	return &remote{
		client: c,
	}, nil
}

func (r *remote) authenticate() generated.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+r.bearerToken)
		return nil
	}
}
