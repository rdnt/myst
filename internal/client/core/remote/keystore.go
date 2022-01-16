package remote

import (
	"context"
	"encoding/hex"
	"fmt"

	"myst/internal/client/core/domain/keystore"
	"myst/internal/server/api/http/generated"
)

func (r *remote) CreateKeystore(name string, payload []byte) error {
	fmt.Println("CreateKeystore", name, payload)

	if r.bearerToken == "" {
		return fmt.Errorf("not signed in")
	}

	_, err := r.client.CreateKeystoreWithResponse(
		context.Background(), generated.CreateKeystoreJSONRequestBody{
			Name:    name,
			Payload: hex.EncodeToString(payload),
		},
		r.authenticate(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *remote) Keystores() ([]*keystore.Keystore, error) {
	return nil, nil
}

func (r *remote) parseKeystore(gen *generated.Keystore) (*keystore.Keystore, error) {
	if gen == nil {
		return nil, ErrInvalidResponse
	}

	return keystore.New(
		keystore.WithId(gen.Id),
	)
}
