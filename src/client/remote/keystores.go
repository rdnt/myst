package remote

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"

	"myst/pkg/crypto"
	"myst/src/client/application/domain/keystore"
	// TODO: @rdnt @@@: this is clearly a shared package, re-setup shared pkg
	//  dir and move it there
	"myst/src/server/rest/generated"
)

func (r *remote) CreateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	if !r.Authenticated() {
		return keystore.Keystore{}, ErrNotAuthenticated
	}

	jk := keystoreToJSON(k)

	b, err := json.Marshal(jk)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to marshal keystore")
	}

	b, err = crypto.AES256CBC_Encrypt(k.Key, b)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "aes256cbc encrypt failed")
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

	if res.StatusCode() != http.StatusCreated || res.JSON201 == nil {
		return keystore.Keystore{}, errors.New("invalid response")
	}

	k.RemoteId = (*res.JSON201).Id

	return k, nil
}

func (r *remote) UpdateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	if !r.Authenticated() {
		return keystore.Keystore{}, ErrNotAuthenticated
	}

	jk := keystoreToJSON(k)

	b, err := json.Marshal(jk)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to marshal keystore")
	}

	b, err = crypto.AES256CBC_Encrypt(k.Key, b)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "aes256cbc encrypt failed")
	}

	res, err := r.client.UpdateKeystoreWithResponse(context.Background(), k.RemoteId, generated.UpdateKeystoreRequest{
		Name:    &k.Name,
		Payload: &b,
	})
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to update keystore")
	}

	if res.StatusCode() != http.StatusOK || res.JSON200 == nil {
		return keystore.Keystore{}, errors.New("invalid response")
	}

	k2, err := KeystoreFromJSON(*res.JSON200, k.Key)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to parse keystore")
	}

	k2.Key = k.Key

	return k2, nil
}

func (r *remote) DeleteKeystore(id string) error {
	if !r.Authenticated() {
		return ErrNotAuthenticated
	}

	res, err := r.client.DeleteKeystoreWithResponse(context.Background(), id)
	if err != nil {
		return errors.Wrap(err, "failed to delete keystore")
	}

	if res.StatusCode() != http.StatusOK {
		return errors.New("invalid response")
	}

	return nil
}

func (r *remote) Keystores(privateKey []byte) (map[string]keystore.Keystore, error) {
	if !r.Authenticated() {
		return nil, ErrNotAuthenticated
	}

	res, err := r.client.KeystoresWithResponse(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get keystores")
	}

	if res.StatusCode() != http.StatusOK || res.JSON200 == nil {
		return nil, errors.New("invalid response")
	}

	restKeystores := *res.JSON200

	invs, err := r.Invitations()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get invitations")
	}

	keystores := make(map[string]keystore.Keystore)

	for _, restKeystore := range restKeystores {
		var keystoreKey []byte

		// find the finalized invitation for this keystore and decrypt its payload
		for _, inv := range invs {
			if inv.RemoteKeystoreId == restKeystore.Id && inv.Finalized() {
				var pub []byte
				if inv.Inviter.Id == r.CurrentUser().Id {
					pub = inv.Invitee.PublicKey
				} else {
					pub = inv.Inviter.PublicKey
				}

				sharedSecret, err := curve25519.X25519(privateKey, pub)
				if err != nil {
					return nil, errors.Wrap(err, "failed to derive shared secret")
				}

				keystoreKey, err = crypto.AES256CBC_Decrypt(sharedSecret, inv.EncryptedKeystoreKey)
				if err != nil {
					return nil, errors.WithMessage(err, "failed to decrypt keystore key")
				}

				break
			}
		}

		if keystoreKey == nil {
			continue
		}

		k, err := KeystoreFromJSON(restKeystore, keystoreKey)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to parse keystore")
		}

		keystores[k.Id] = k
	}

	return keystores, nil
}
