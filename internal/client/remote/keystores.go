package remote

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"

	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/keystorerepo"
	"myst/internal/server/api/http/generated"
	"myst/pkg/crypto"
	"myst/pkg/logger"
)

func (r *remote) CreateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	if !r.SignedIn() {
		return keystore.Keystore{}, ErrSignedOut
	}

	//id := k.Id

	jk := keystorerepo.KeystoreToJSON(k)

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

	if res.JSON200 == nil {
		return keystore.Keystore{}, fmt.Errorf("invalid response")
	}

	//err = r.keystores.UpdateKeystore(k)
	//if err != nil {
	//	return keystore.Keystore{}, errors.Wrap(err, "failed to update keystore with remote id")
	//}

	//k, err = KeystoreFromJSON(*res.JSON200, k.Key)
	//if err != nil {
	//	return keystore.Keystore{}, errors.WithMessage(err, "failed to parse keystore")
	//}

	k.RemoteId = (*res.JSON200).Id

	return k, nil
}

func (r *remote) Keystore(id string) (keystore.Keystore, error) {
	if !r.SignedIn() {
		return keystore.Keystore{}, ErrSignedOut
	}

	res, err := r.client.KeystoreWithResponse(
		context.Background(), id,
	)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to create keystore")
	}

	if res.JSON200 == nil {
		return keystore.Keystore{}, fmt.Errorf("invalid response")
	}

	invs, err := r.Invitations()
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to get invitations")
	}

	var keystoreKey []byte
	for _, inv := range invs {
		if inv.KeystoreId == id && inv.Finalized() {
			symKey, err := curve25519.X25519(r.privateKey, inv.InviterKey)
			if err != nil {
				return keystore.Keystore{}, errors.Wrap(err, "failed to create asymmetric key")
			}

			logger.Error("@@@ ###################### SYMMETRIC WHEN SYNC", string(symKey))

			keystoreKey, err = crypto.AES256CBC_Decrypt(symKey, inv.KeystoreKey)
			if err != nil {
				return keystore.Keystore{}, errors.Wrap(err, "failed to decrypt keystore key")
			}

			break
		}
	}

	if keystoreKey == nil {
		panic(err)
	}

	k, err := KeystoreFromJSON(*res.JSON200, keystoreKey)
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
	res, err := r.client.KeystoresWithResponse(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get keystores")
	}

	if res.JSON200 == nil {
		return nil, fmt.Errorf("invalid response")
	}

	restKeystores := *res.JSON200
	keystores := make(map[string]keystore.Keystore)

	for _, restKeystore := range restKeystores {
		invs, err := r.Invitations()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get invitations")
		}

		var keystoreKey []byte
		for _, inv := range invs {
			if inv.KeystoreId == restKeystore.Id && inv.Finalized() {
				symKey, err := curve25519.X25519(r.privateKey, inv.InviterKey)
				if err != nil {
					return nil, errors.Wrap(err, "failed to create asymmetric key")
				}

				logger.Error("@@@ ###################### SYMMETRIC WHEN SYNC", string(symKey))

				keystoreKey, err = crypto.AES256CBC_Decrypt(symKey, inv.KeystoreKey)
				if err != nil {
					return nil, errors.Wrap(err, "failed to decrypt keystore key")
				}

				break
			}
		}

		if keystoreKey == nil {
			panic(err)
		}

		k, err := KeystoreFromJSON(restKeystore, keystoreKey)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to parse keystore")
		}

		keystores[k.Id] = k
	}

	//ks := map[string]keystore.Keystore{}
	//for _, k := range *res.JSON200 {
	//	k2, err := KeystoreFromJSON(k)
	//	if err != nil {
	//		return nil, errors.WithMessage(err, "failed to parse keystore")
	//	}
	//
	//	ks[k2.Id] = k2
	//}

	return keystores, nil
}

func (r *remote) DeleteKeystore(id string) error {
	//TODO implement me
	panic("implement me")
}
