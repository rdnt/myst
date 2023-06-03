package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"

	"myst/pkg/crypto"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/enclaverepo"
	"myst/src/server/rest/generated"
)

func (r *remote) CreateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	if !r.Authenticated() {
		return keystore.Keystore{}, ErrNotAuthenticated
	}

	jk := enclaverepo.KeystoreToJSON(k)

	b, err := json.Marshal(jk)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to marshal keystore to json")
	}

	b, err = crypto.AES256CBC_Encrypt(k.Key, b)
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

	if res.JSON201 == nil {
		return keystore.Keystore{}, fmt.Errorf("invalid response: %s", string(res.Body))
	}

	// err = r.keystores.UpdateKeystore(k)
	// if err != nil {
	//	return keystore.KeystoreJSON{}, errors.Wrap(err, "failed to update keystore with remote id")
	// }

	// k, err = KeystoreFromJSON(*res.JSON200, k.Key)
	// if err != nil {
	//	return keystore.KeystoreJSON{}, errors.WithMessage(err, "failed to parse keystore")
	// }

	k.RemoteId = (*res.JSON201).Id

	return k, nil
}

func (r *remote) UpdateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	if !r.Authenticated() {
		return keystore.Keystore{}, ErrNotAuthenticated
	}

	jk := enclaverepo.KeystoreToJSON(k)

	b, err := json.Marshal(jk)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to marshal keystore to json")
	}

	b, err = crypto.AES256CBC_Encrypt(k.Key, b)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "aes256cbc encrypt failed")
	}

	res, err := r.client.UpdateKeystoreWithResponse(context.Background(), k.RemoteId, generated.UpdateKeystoreRequest{
		Name:    &k.Name,
		Payload: &b,
	})
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to get keystores")
	}

	if res.JSON200 == nil {
		return keystore.Keystore{}, fmt.Errorf("invalid response")
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
		return errors.Wrap(err, "failed to get keystores")
	}

	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("request failed with status code %d", res.StatusCode())
	}

	return nil
}

func (r *remote) Keystores(privateKey []byte) (map[string]keystore.Keystore, error) {
	if !r.Authenticated() {
		return map[string]keystore.Keystore{}, ErrNotAuthenticated
	}

	res, err := r.client.KeystoresWithResponse(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get keystores")
	}

	if res.JSON200 == nil {
		return nil, fmt.Errorf("invalid response")
	}

	restKeystores := *res.JSON200
	keystores := make(map[string]keystore.Keystore)

	invs, err := r.Invitations()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get invitations")
	}

	for _, restKeystore := range restKeystores {
		var keystoreKey []byte
		for _, inv := range invs {
			// find the finalized invitation for this keystore and decrypt its payload
			if inv.Keystore.RemoteId == restKeystore.Id && inv.Finalized() {

				var pub []byte
				if inv.Inviter.Id == r.CurrentUser().Id {
					pub = inv.Invitee.PublicKey
				} else {
					pub = inv.Inviter.PublicKey
				}

				symKey, err := curve25519.X25519(privateKey, pub)
				if err != nil {
					return nil, errors.Wrap(err, "failed to create asymmetric key")
				}

				keystoreKey, err = crypto.AES256CBC_Decrypt(symKey, inv.EncryptedKeystoreKey)
				if err != nil {
					return nil, errors.Wrap(err, "failed to decrypt keystore key")
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

	// ks := map[string]keystore.KeystoreJSON{}
	// for _, k := range *res.JSON200 {
	//	k2, err := KeystoreFromJSON(k)
	//	if err != nil {
	//		return nil, errors.WithMessage(err, "failed to parse keystore")
	//	}
	//
	//	ks[k2.Id] = k2
	// }

	return keystores, nil
}
