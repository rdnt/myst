package enclaverepo

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"myst/src/client/application/domain/credentials"
	"myst/src/client/application/domain/entry"

	"myst/src/client/application/domain/keystore"
)

type enclaveJSON struct {
	Keystores   map[string]keystoreJSON `json:"keystores"`
	Keys        map[string][]byte       `json:"keys"`
	Credentials credentialsJSON         `json:"credentials"`
}

type credentialsJSON struct {
	Address    string `json:"address"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PublicKey  []byte `json:"publicKey"`
	PrivateKey []byte `json:"privateKey"`
}

type keystoreJSON struct {
	Id        string      `json:"id"`
	RemoteId  string      `json:"remoteId"`
	Name      string      `json:"name"`
	Version   int         `json:"version"`
	Entries   []entryJSON `json:"entries"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type entryJSON struct {
	Id        string    `json:"id"`
	Website   string    `json:"website"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func enclaveToJSON(e *enclave) ([]byte, error) {
	ks := map[string]keystoreJSON{}

	eks, err := e.keystoresWithKeys()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get keystores")
	}
	for _, k := range eks {
		ks[k.Id] = keystoreToJSON(k)
	}

	jrem := credentialsJSON{
		Address:    e.creds.Address,
		Username:   e.creds.Username,
		Password:   e.creds.Password,
		PublicKey:  e.creds.PublicKey,
		PrivateKey: e.creds.PrivateKey,
	}

	b, err := json.Marshal(enclaveJSON{
		Keystores:   ks,
		Keys:        e.keys,
		Credentials: jrem,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal enclave")
	}

	return b, nil
}

func enclaveFromJSON(b, encSalt, signSalt []byte) (*enclave, error) {
	e := &enclaveJSON{}

	err := json.Unmarshal(b, e)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal enclave")
	}

	ks := map[string]keystore.Keystore{}

	for _, k := range e.Keystores {
		ks[k.Id] = keystoreFromJSON(k)
	}

	creds := credentials.Credentials{
		Address:    e.Credentials.Address,
		Username:   e.Credentials.Username,
		Password:   e.Credentials.Password,
		PublicKey:  e.Credentials.PublicKey,
		PrivateKey: e.Credentials.PrivateKey,
	}

	return &enclave{
		keystores: ks,
		creds:     creds,
		encSalt:   encSalt,
		signSalt:  signSalt,
		keys:      e.Keys,
	}, nil
}

func keystoreToJSON(k keystore.Keystore) keystoreJSON {
	entries := make([]entryJSON, len(k.Entries))

	for i, e := range lo.Values(k.Entries) {
		entries[i] = entryJSON{
			Id:        e.Id,
			Website:   e.Website,
			Username:  e.Username,
			Password:  e.Password,
			Notes:     e.Notes,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		}
	}

	return keystoreJSON{
		Id:        k.Id,
		RemoteId:  k.RemoteId,
		Name:      k.Name,
		Version:   k.Version,
		Entries:   entries,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

func keystoreFromJSON(k keystoreJSON) keystore.Keystore {
	entries := make(map[string]entry.Entry, len(k.Entries))

	for _, e := range k.Entries {
		e := entry.Entry{
			Id:        e.Id,
			Website:   e.Website,
			Username:  e.Username,
			Password:  e.Password,
			Notes:     e.Notes,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		}

		entries[e.Id] = e
	}

	return keystore.Keystore{
		Id:        k.Id,
		RemoteId:  k.RemoteId,
		Name:      k.Name,
		Version:   k.Version,
		Entries:   entries,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}
