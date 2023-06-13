package remote

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"myst/src/client/application"
	"myst/src/client/application/domain/user"
	"myst/src/client/pkg/enclave"
	"myst/src/server/rest/generated"
)

var (
	ErrAuthenticationFailed = enclave.ErrAuthenticationFailed
	ErrInvalidResponse      = fmt.Errorf("invalid response")
)

type remote struct {
	address string

	client *generated.ClientWithResponses

	bearerToken string
	user        *user.User
}

func New(opts ...Option) (application.Remote, error) {
	r := &remote{
		client:      nil,
		bearerToken: "",
		user:        nil,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(r)
		}
	}

	var err error
	r.client, err = generated.NewClientWithResponses(
		r.address,
		generated.WithRequestEditorFn(r.authenticate()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create remote client")
	}

	return r, nil
}

// func (r *remote) UploadKeystore(id string) (keystore.keystore, error) {
//	if !r.Authenticated() {
//		return keystore.keystore{}, ErrNotAuthenticated
//	}
//
//	k, err := r.keystores.keystore(id)
//	if err != nil {
//		return keystore.keystore{}, errors.WithMessage(err, "failed to get keystore")
//	}
//
//	keystoreKey, err := r.keystores.EncryptedKeystoreKey(k.Id)
//	if err != nil {
//		return keystore.keystore{}, errors.WithMessage(err, "failed to get keystore key")
//	}
//
//	jk := repository.KeystoreToJSON(k)
//
//	b, err := json.Marshal(jk)
//	if err != nil {
//		return keystore.keystore{}, err
//	}
//
//	b, err = crypto.AES256CBC_Encrypt(keystoreKey, b)
//	if err != nil {
//		return keystore.keystore{}, err
//	}
//
//	if !r.Authenticated() {
//		return keystore.keystore{}, ErrNotAuthenticated
//	}
//
//	res, err := r.client.CreateKeystoreWithResponse(
//		context.Background(), generated.CreateKeystoreJSONRequestBody{
//			Name:    k.Name,
//			Payload: b,
//		},
//	)
//	if err != nil {
//		return keystore.keystore{}, err
//	}
//
//	if res.JSON200 == nil {
//		return keystore.keystore{}, fmt.Errorf("invalid response")
//	}
//
//	k, err = KeystoreFromJSON(*res.JSON200)
//	if err != nil {
//		return keystore.keystore{}, errors.WithMessage(err, "failed to parse keystore")
//	}
//
//	return k, nil
// }

func (r *remote) Authenticate(username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := r.client.LoginWithResponse(
		ctx, generated.LoginJSONRequestBody{
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		return errors.Wrap(err, "authentication failed")
	}

	if res.JSON200 == nil {
		return ErrInvalidResponse
	}

	resp := *res.JSON200

	if resp.Token == "" {
		return errors.New("invalid token")
	}

	u := user.User{
		Id:        resp.User.Id,
		Username:  resp.User.Username,
		PublicKey: resp.User.PublicKey,
	}

	r.user = &u
	r.bearerToken = resp.Token

	fmt.Println("Signed in.")

	return nil
}

func (r *remote) Authenticated() bool {
	return r.bearerToken != ""
}

func (r *remote) CurrentUser() *user.User {
	return r.user
}

// func (r *remote) SignOut() error {
// 	fmt.Println("Signing out from remote...")
//
// 	r.bearerToken = ""
// 	r.user = nil
//
// 	fmt.Println("Signed out.")
//
// 	return nil
// }

func (r *remote) Register(username, password string, publicKey []byte) (user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := r.client.RegisterWithResponse(
		ctx, generated.RegisterJSONRequestBody{
			Username:  username,
			Password:  password,
			PublicKey: publicKey,
		},
	)
	if err != nil {
		return user.User{}, errors.Wrap(err, "failed to sign in")
	}

	if res.JSON201 == nil {
		return user.User{}, ErrInvalidResponse
	}

	resp := *res.JSON201

	if resp.Token == "" {
		return user.User{}, errors.New("invalid token")
	}

	u := user.User{
		Id:        resp.User.Id,
		Username:  resp.User.Username,
		PublicKey: publicKey,
	}

	r.user = &u
	r.bearerToken = resp.Token

	return u, nil
}

func (r *remote) authenticate() generated.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+r.bearerToken)
		return nil
	}
}
