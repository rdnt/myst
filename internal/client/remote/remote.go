package remote

import (
	"context"
	"fmt"
	"net/http"

	"myst/internal/client/application/domain/invitation"
	"myst/internal/client/application/domain/keystore"
	"myst/internal/server/api/http/generated"
	"myst/pkg/crypto"

	"github.com/pkg/errors"

	"golang.org/x/crypto/curve25519"

	"myst/pkg/enclave"
)

var (
	ErrAuthenticationFailed = enclave.ErrAuthenticationFailed
	ErrInvalidResponse      = fmt.Errorf("invalid response")
)

// Remote is a remote repository that holds upstream keystores/invitations
type Remote interface {
	CreateInvitation(inv invitation.Invitation) (invitation.Invitation, error)
	Invitation(id string) (invitation.Invitation, error)
	UpdateInvitation(inv invitation.Invitation) error
	Invitations() (map[string]invitation.Invitation, error)
	DeleteInvitation(id string) error

	CreateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	Keystore(id string) (keystore.Keystore, error)
	UpdateKeystore(k keystore.Keystore) error
	Keystores() (map[string]keystore.Keystore, error)
	DeleteKeystore(id string) error

	SignIn(username, password string) error
	SignedIn() bool
	SignOut() error

	// keystores
	//UploadKeystore(id string) (*generated.Keystore, error)
	//Keystores() ([]keystore.Keystore, error)
	//
	//// invitations
	//CreateInvitation(opts ...invitation.Option) (invitation.Invitation, error)
	//AcceptInvitation(keystoreId, invitationId string) (invitation.Invitation, error)
	//
	//Invitation(id string) (invitation.Invitation, error)
	//UpdateInvitation(k invitation.Invitation) error
	//DeleteInvitation(id string) error
	//
	//Invitations() ([]invitation.Invitation, error)
	//
	//// TODO: refine to dynamically find local keystoreId
	//FinalizeInvitation(localKeystoreId, keystoreId, invitationId string) (*generated.Invitation, error)
}

type remote struct {
	client *generated.ClientWithResponses

	keystores keystore.Service

	bearerToken string

	publicKey  []byte
	privateKey []byte
}

func New(keystores keystore.Service, address string) (Remote, error) {
	r := &remote{
		client:      nil,
		bearerToken: "",
		publicKey:   nil,
		privateKey:  nil,
		keystores:   keystores,
	}

	pub, key, err := newKeypair()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate public/private keypair")
	}

	r.client, err = generated.NewClientWithResponses(
		address,
		generated.WithRequestEditorFn(r.authenticate()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create remote client")
	}

	r.publicKey = pub
	r.privateKey = key

	return r, nil
}

//func (r *remote) UploadKeystore(id string) (keystore.Keystore, error) {
//	if !r.SignedIn() {
//		return keystore.Keystore{}, ErrSignedOut
//	}
//
//	k, err := r.keystores.Keystore(id)
//	if err != nil {
//		return keystore.Keystore{}, errors.WithMessage(err, "failed to get keystore")
//	}
//
//	keystoreKey, err := r.keystores.KeystoreKey(k.Id)
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
//}

func (r *remote) SignIn(username, password string) error {
	fmt.Println("Signing in to remote...")

	res, err := r.client.LoginWithResponse(
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

	r.bearerToken = string(token)

	fmt.Println("Signed in.")

	return nil
}

func (r *remote) SignedIn() bool {
	return r.bearerToken != ""
}

func (r *remote) SignOut() error {
	fmt.Println("Signing out from remote...")

	r.bearerToken = ""

	fmt.Println("Signed out.")

	return nil
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
