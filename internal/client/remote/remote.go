package remote

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"

	"myst/internal/client/application/domain/invitation"
	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/application/domain/user"
	"myst/internal/server/api/http/generated"
	"myst/pkg/crypto"
	"myst/pkg/enclave"
)

var (
	ErrAuthenticationFailed = enclave.ErrAuthenticationFailed
	ErrInvalidResponse      = fmt.Errorf("invalid response")
)

// Remote is a remote repository that holds upstream keystores/invitations
type Remote interface {
	Address() string
	CreateInvitation(inv invitation.Invitation) (invitation.Invitation, error)
	Invitation(id string) (invitation.Invitation, error)
	AcceptInvitation(id string) (invitation.Invitation, error)
	DeclineOrCancelInvitation(id string) (invitation.Invitation, error)
	FinalizeInvitation(invitationId string, keystoreKey, privateKey []byte) (invitation.Invitation, error)
	UpdateInvitation(inv invitation.Invitation) error
	Invitations() (map[string]invitation.Invitation, error)
	DeleteInvitation(id string) error

	CreateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	Keystore(id string, privateKey []byte) (keystore.Keystore, error)
	UpdateKeystore(k keystore.Keystore) error
	Keystores(privateKey []byte) (map[string]keystore.Keystore, error)
	DeleteKeystore(id string) error

	SignIn(username, password string) error
	Register(username, password string, publicKey []byte) (user.User, error)
	SignOut() error
	SignedIn() bool
	CurrentUser() *user.User
	UserByUsername(username string) (user.User, error)

	// keystores
	// UploadKeystore(id string) (*generated.Keystore, error)
	// Keystores() ([]keystore.Keystore, error)
	//
	// // invitations
	// CreateInvitation(opts ...invitation.Option) (invitation.Invitation, error)
	// AcceptInvitation(keystoreId, invitationId string) (invitation.Invitation, error)
	//
	// Invitation(id string) (invitation.Invitation, error)
	// UpdateInvitation(k invitation.Invitation) error
	// DeleteInvitation(id string) error
	//
	// Invitations() ([]invitation.Invitation, error)
	//
	// // TODO: refine to dynamically find local keystoreId
	// FinalizeInvitation(localKeystoreId, keystoreId, invitationId string) (*generated.Invitation, error)
}

type remote struct {
	address string

	client *generated.ClientWithResponses

	bearerToken string
	user        *user.User
}

func New(opts ...Option) (Remote, error) {
	r := &remote{
		client:      nil,
		bearerToken: "",
		user:        nil,
	}

	for _, opt := range opts {
		opt(r)
	}

	var err error
	r.client, err = generated.NewClientWithResponses(
		r.address+"/api",
		generated.WithRequestEditorFn(r.authenticate()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create remote client")
	}

	return r, nil
}

// func (r *remote) UploadKeystore(id string) (keystore.Keystore, error) {
//	if !r.SignedIn() {
//		return keystore.Keystore{}, ErrSignedOut
//	}
//
//	k, err := r.keystores.Keystore(id)
//	if err != nil {
//		return keystore.Keystore{}, errors.WithMessage(err, "failed to get keystore")
//	}
//
//	keystoreKey, err := r.keystores.EncryptedKeystoreKey(k.Id)
//	if err != nil {
//		return keystore.Keystore{}, errors.WithMessage(err, "failed to get keystore key")
//	}
//
//	jk := keystorerepo.KeystoreToJSON(k)
//
//	b, err := json.Marshal(jk)
//	if err != nil {
//		return keystore.Keystore{}, err
//	}
//
//	b, err = crypto.AES256CBC_Encrypt(keystoreKey, b)
//	if err != nil {
//		return keystore.Keystore{}, err
//	}
//
//	if !r.SignedIn() {
//		return keystore.Keystore{}, ErrSignedOut
//	}
//
//	res, err := r.client.CreateKeystoreWithResponse(
//		context.Background(), generated.CreateKeystoreJSONRequestBody{
//			Name:    k.Name,
//			Payload: b,
//		},
//	)
//	if err != nil {
//		return keystore.Keystore{}, err
//	}
//
//	if res.JSON200 == nil {
//		return keystore.Keystore{}, fmt.Errorf("invalid response")
//	}
//
//	k, err = KeystoreFromJSON(*res.JSON200)
//	if err != nil {
//		return keystore.Keystore{}, errors.WithMessage(err, "failed to parse keystore")
//	}
//
//	return k, nil
// }

func (r *remote) SignIn(username, password string) error {
	fmt.Println("Signing in to remote...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := r.client.LoginWithResponse(
		ctx, generated.LoginJSONRequestBody{
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

	resp := *res.JSON200

	if resp.Token == "" {
		return errors.New("invalid token")
	}

	u := user.User{
		Id:       resp.User.Id,
		Username: resp.User.Username,
	}

	r.user = &u
	r.bearerToken = resp.Token

	fmt.Println("Signed in.")

	return nil
}

func (r *remote) SignedIn() bool {
	return r.bearerToken != ""
}

func (r *remote) CurrentUser() *user.User {
	return r.user
}

func (r *remote) UserByUsername(username string) (user.User, error) {
	if !r.SignedIn() {
		return user.User{}, ErrSignedOut
	}

	res, err := r.client.UserByUsernameWithResponse(
		context.Background(), &generated.UserByUsernameParams{Username: &username},
	)
	if err != nil {
		return user.User{}, err
	}

	if res.JSON200 == nil {
		return user.User{}, fmt.Errorf("invalid response")
	}

	u := UserFromJSON(*res.JSON200)
	if err != nil {
		return user.User{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return u, nil
}

func (r *remote) SignOut() error {
	fmt.Println("Signing out from remote...")

	r.bearerToken = ""
	r.user = nil

	fmt.Println("Signed out.")

	return nil
}

func (r *remote) Register(username, password string, publicKey []byte) (user.User, error) {
	fmt.Println("Registering to remote...")
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

	// TODO: properly sign in after registration
	if resp.Token == "" {
		return user.User{}, errors.New("invalid token")
	}

	u := user.User{
		Id:       resp.User.Id,
		Username: resp.User.Username,
	}

	r.user = &u
	r.bearerToken = resp.Token

	fmt.Println("Registered.")

	return u, nil
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

func (r *remote) authenticate() generated.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+r.bearerToken)
		return nil
	}
}
