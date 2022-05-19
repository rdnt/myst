package remote

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/keystorerepo"
	"myst/internal/server/api/http/generated"
	"myst/pkg/crypto"
)

func (r *remote) CreateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	if !r.SignedIn() {
		return keystore.Keystore{}, ErrSignedOut
	}

	keystoreKey, err := r.keystores.KeystoreKey(k.Id)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get keystore key")
	}

	jk := keystorerepo.KeystoreToJSON(k)

	b, err := json.Marshal(jk)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to marshal keystore to json")
	}

	b, err = crypto.AES256CBC_Encrypt(keystoreKey, b)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "aes256cbc encrypt failed")
	}

	res, err := r.client.CreateKeystoreWithResponse(
		context.Background(), generated.CreateKeystoreJSONRequestBody{
			Name:    k.Name,
			Payload: b,
		},
	)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to create keystore")
	}

	if res.JSON200 == nil {
		return keystore.Keystore{}, fmt.Errorf("invalid response")
	}

	k, err = KeystoreFromJSON(*res.JSON200)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to parse keystore")
	}

	return k, nil
}

func (r *remote) Keystore(id string) (keystore.Keystore, error) {
	if !r.SignedIn() {
		return keystore.Keystore{}, ErrSignedOut
	}

	k, err := r.keystores.Keystore(id)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get keystore key")
	}

	res, err := r.client.KeystoreWithResponse(
		context.Background(), k.Id,
	)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to create keystore")
	}

	if res.JSON200 == nil {
		return keystore.Keystore{}, fmt.Errorf("invalid response")
	}

	k, err = KeystoreFromJSON(*res.JSON200)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to parse keystore")
	}

	return k, nil
}

func (r *remote) UpdateKeystore(k keystore.Keystore) error {
	//TODO implement me
	panic("implement me")
}

func (r *remote) Keystores() (map[string]keystore.Keystore, error) {
	//TODO implement me
	panic("implement me")
}

func (r *remote) DeleteKeystore(id string) error {
	//TODO implement me
	panic("implement me")
}
