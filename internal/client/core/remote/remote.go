package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"myst/internal/client/core/domain/invitation"
	"myst/internal/client/core/domain/keystore"
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
	Keystores() ([]*keystore.Keystore, error)
	CreateInvitation(keystoreId, inviteeId, publicKey string) (*invitation.Invitation, error)
	AcceptInvitation(invitationId string, publicKey []byte) (*invitation.Invitation, error)
	FinalizeInvitation(invitationId string, keystoreKey []byte) (*invitation.Invitation, error)
}

type remote struct {
	client *generated.ClientWithResponses

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

func (r *remote) CreateInvitation(keystoreId, inviteeId, publicKey string) (*invitation.Invitation, error) {
	fmt.Println("CreateInvitation", keystoreId, inviteeId, publicKey)

	if r.bearerToken == "" {
		return nil, fmt.Errorf("not signed in")
	}

	res, err := r.client.CreateInvitationWithResponse(
		context.Background(), keystoreId, generated.CreateInvitationJSONRequestBody{
			InviteeId: inviteeId,
			PublicKey: publicKey,
		},
		r.authenticate(),
	)
	if err != nil {
		return nil, err
	}

	if res.JSON200 == nil {
		return nil, ErrInvalidResponse
	}

	inv := *res.JSON200

	return invitation.New(
		invitation.WithId(inv.Id),
	)
}

func (r *remote) AcceptInvitation(invitationId string, publicKey []byte) (*invitation.Invitation, error) {
	panic("implement me")
}

func (r *remote) FinalizeInvitation(invitationId string, keystoreKey []byte) (*invitation.Invitation, error) {
	panic("implement me")
}

func New() (*remote, error) {
	c, err := generated.NewClientWithResponses("http://localhost:8080")
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
