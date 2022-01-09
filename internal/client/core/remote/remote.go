package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	CreateKeystoreInvitation(keystoreId, inviteeId, publicKey string) error
}

type remote struct {
	client *generated.Client

	apiKey string
}

func (r *remote) SignIn(username, password string) error {
	if r.apiKey != "" {
		return fmt.Errorf("already signed in")
	}

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

	fmt.Println("received", string(b))

	var token generated.AuthToken
	err = json.Unmarshal(b, &token)
	if err != nil {
		return err
	}

	if token == "" {
		return fmt.Errorf("invalid token")
	}

	r.apiKey = string(token)

	return nil
}

func (r *remote) SignOut() error {
	if r.apiKey == "" {
		return fmt.Errorf("already signed out")
	}

	r.apiKey = ""

	return nil
}

func (r *remote) Keystores() ([]*keystore.Keystore, error) {
	return nil, nil
}

func (r *remote) CreateKeystoreInvitation(keystoreId, inviteeId, publicKey string) error {
	if r.apiKey == "" {
		return fmt.Errorf("not signed in")
	}

	resp, err := r.client.CreateKeystoreInvitation(
		context.Background(), keystoreId, generated.CreateKeystoreInvitationJSONRequestBody{
			InviteeId: inviteeId,
			PublicKey: publicKey,
		},
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+r.apiKey)

			return nil
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

	fmt.Println("received", string(b))

	var inv generated.Invitation
	err = json.Unmarshal(b, &inv)
	if err != nil {
		return err
	}

	fmt.Println("@@@@@@@@@@@@@@@@@")
	fmt.Println(inv)

	return nil
}

func New() (*remote, error) {
	c, err := generated.NewClient("http://localhost:8081")
	if err != nil {
		return nil, err
	}

	return &remote{
		client: c,
	}, nil
}
