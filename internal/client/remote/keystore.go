package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"myst/internal/client/enclaverepo"
	"myst/pkg/crypto"

	"myst/internal/client/application/domain/keystore"
	"myst/internal/server/api/http/generated"
)

func (r *remote) CreateKeystore(name string, key []byte, k *keystore.Keystore) (*generated.Keystore, error) {
	if r.bearerToken == "" {
		return nil, fmt.Errorf("not signed in")
	}

	jk := enclaverepo.KeystoreToJSON(k)

	b, err := json.Marshal(jk)
	if err != nil {
		return nil, err
	}

	b, err = crypto.AES256CBC_Encrypt(key, b)
	if err != nil {
		return nil, err
	}

	res, err := r.client.CreateKeystoreWithResponse(
		context.Background(), generated.CreateKeystoreJSONRequestBody{
			Name:    name,
			Payload: b,
		},
		r.authenticate(),
	)
	if err != nil {
		return nil, err
	}

	if res.JSON200 == nil {
		return nil, fmt.Errorf("invalid response")
	}

	return res.JSON200, nil
}

func (r *remote) Keystore(id string) (*generated.Keystore, error) {
	if r.bearerToken == "" {
		return nil, fmt.Errorf("not signed in")
	}

	res, err := r.client.KeystoreWithResponse(
		context.Background(), id,
		r.authenticate(),
	)
	if err != nil {
		return nil, err
	}

	return res.JSON200, nil
}

func (r *remote) Keystores() ([]*generated.Keystore, error) {
	if r.bearerToken == "" {
		return nil, fmt.Errorf("not signed in")
	}

	res, err := r.client.KeystoresWithResponse(
		context.Background(),
		r.authenticate(),
	)
	if err != nil {
		return nil, err
	}

	if res.JSON200 == nil {
		return nil, ErrInvalidResponse
	}

	ks := []*generated.Keystore{}

	for _, k := range *res.JSON200 {
		ks = append(ks, &k)
	}

	return ks, nil
}

//
//func (r *remote) parseKeystore(gen *generated.Keystore) (*keystore.Keystore, error) {
//	if gen == nil {
//		return nil, ErrInvalidResponse
//	}
//
//	return keystore.New(
//		keystore.WithId(gen.Id),
//	), nil
//}
