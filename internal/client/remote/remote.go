package remote

import (
	"encoding/json"
	"fmt"
	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/keystorerepo"
	"myst/internal/client/remote/client"
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
	//invitation.Repository
	//keystore.Repository

	SignIn(username, password string) error
	SignOut() error

	// keystores
	UploadKeystore(id string) (*generated.Keystore, error)
	Keystores() ([]*generated.Keystore, error)

	// invitations
	CreateInvitation(keystoreId, inviteeId string) (*generated.Invitation, error)
	AcceptInvitation(keystoreId, invitationId string) (*generated.Invitation, error)

	// TODO: refine to dynamically find local keystoreId
	FinalizeInvitation(localKeystoreId, keystoreId, invitationId string) (*generated.Invitation, error)
}

type remote struct {
	client    client.Client
	keystores keystore.Service

	bearerToken string

	publicKey  []byte
	privateKey []byte
}

func New(keystores keystore.Service) (Remote, error) {
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

	rc, err := client.New()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create remote client")
	}

	r.client = rc
	r.publicKey = pub
	r.privateKey = key

	return r, nil
}

func (r *remote) UploadKeystore(id string) (*generated.Keystore, error) {
	k, err := r.keystores.Keystore(id)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get keystore")
	}

	keystoreKey, err := r.keystores.KeystoreKey(k.Id())
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get keystore key")
	}

	jk := keystorerepo.KeystoreToJSON(k)

	b, err := json.Marshal(jk)
	if err != nil {
		return nil, err
	}

	b, err = crypto.AES256CBC_Encrypt(keystoreKey, b)
	if err != nil {
		return nil, err
	}

	rk, err := r.client.UploadKeystore(k.Name(), b)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to upload keystore")
	}

	return rk, nil
}

func (r *remote) SignIn(username, password string) error {
	fmt.Println("Signing in to remote...")

	err := r.client.SignIn(username, password)
	if err != nil {
		return err
	}

	fmt.Println("Signed in.")

	return nil
}

func (r *remote) SignOut() error {
	fmt.Println("Signing out from remote...")

	err := r.client.SignOut()
	if err != nil {
		return err
	}

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
