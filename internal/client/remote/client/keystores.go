package client

import (
	"context"
	"fmt"

	"myst/internal/server/api/http/generated"
)

func (c *client) UploadKeystore(name string, payload []byte) (generated.Keystore, error) {
	if !c.SignedIn() {
		return generated.Keystore{}, ErrNotSignedIn
	}

	res, err := c.client.CreateKeystoreWithResponse(
		context.Background(), generated.CreateKeystoreJSONRequestBody{
			Name:    name,
			Payload: payload,
		},
	)
	if err != nil {
		return generated.Keystore{}, err
	}

	if res.JSON200 == nil {
		return generated.Keystore{}, fmt.Errorf("invalid response")
	}

	return *res.JSON200, nil
}

func (c *client) Keystore(id string) (generated.Keystore, error) {
	if !c.SignedIn() {
		return generated.Keystore{}, ErrNotSignedIn
	}

	res, err := c.client.KeystoreWithResponse(
		context.Background(), id,
	)
	if err != nil {
		return generated.Keystore{}, err
	}

	return *res.JSON200, nil
}

func (c *client) Keystores() ([]generated.Keystore, error) {
	if !c.SignedIn() {
		return nil, ErrNotSignedIn
	}

	res, err := c.client.KeystoresWithResponse(
		context.Background(),
	)
	if err != nil {
		return nil, err
	}

	if res.JSON200 == nil {
		return nil, ErrInvalidResponse
	}

	ks := []generated.Keystore{}

	for _, k := range *res.JSON200 {
		ks = append(ks, k)
	}

	return ks, nil
}
