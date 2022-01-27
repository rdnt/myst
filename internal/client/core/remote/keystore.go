package remote

import (
	"context"
	"encoding/hex"
	"fmt"

	"myst/internal/client/core/domain/keystore"
	"myst/internal/server/api/http/generated"
)

func (r *remote) CreateKeystore(name string, payload []byte) (*keystore.Keystore, error) {
	fmt.Printf("CreateKeystore %s %x\n", name, payload)

	if r.bearerToken == "" {
		return nil, fmt.Errorf("not signed in")
	}

	res, err := r.client.CreateKeystoreWithResponse(
		context.Background(), generated.CreateKeystoreJSONRequestBody{
			Name:    name,
			Payload: hex.EncodeToString(payload),
		},
		r.authenticate(),
	)
	if err != nil {
		return nil, err
	}

	return r.parseKeystore(res.JSON200)
}

func (r *remote) Keystore(id string) (*keystore.Keystore, error) {
	fmt.Println("Keystore", id)

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

	return r.parseKeystore(res.JSON200)
}

func (r *remote) Keystores() ([]*keystore.Keystore, error) {
	fmt.Println("Keystores")

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

	ks := []*keystore.Keystore{}

	for _, k := range *res.JSON200 {
		gen, err := r.parseKeystore(&k)
		if err != nil {
			return nil, err
		}

		ks = append(ks, gen)
	}

	return ks, nil
}

func (r *remote) parseKeystore(gen *generated.Keystore) (*keystore.Keystore, error) {
	if gen == nil {
		return nil, ErrInvalidResponse
	}

	return keystore.New(
		keystore.WithId(gen.Id),
	), nil
}
