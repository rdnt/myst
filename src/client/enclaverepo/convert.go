package enclaverepo

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"myst/src/client/application/domain/credentials"
	"myst/src/client/application/domain/entry"
	"time"

	"myst/src/client/application/domain/keystore"
)

type enclaveJSON struct {
	Keystores map[string]keystoreJSON `json:"keystores"`
	Keys      map[string][]byte       `json:"keys"`
	Remote    *remoteJSON             `json:"creds,omitempty"`
}

type remoteJSON struct {
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

	var jrem *remoteJSON
	rem := e.creds

	if rem != nil {
		jrem = &remoteJSON{
			Address:    rem.Address,
			Username:   rem.Username,
			Password:   rem.Password,
			PublicKey:  rem.PublicKey,
			PrivateKey: rem.PrivateKey,
		}
	}

	b, err := json.Marshal(enclaveJSON{
		Keystores: ks,
		Keys:      e.keys,
		Remote:    jrem,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal enclave")
	}

	return b, nil
}

func enclaveFromJSON(b, salt []byte) (*enclave, error) {
	e := &enclaveJSON{}

	err := json.Unmarshal(b, e)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal enclave")
	}

	ks := map[string]keystore.Keystore{}

	for _, k := range e.Keystores {
		ks[k.Id] = keystoreFromJSON(k)
	}

	var rem *credentials.Credentials
	jrem := e.Remote

	if jrem != nil {
		rem = &credentials.Credentials{
			Address:    jrem.Address,
			Username:   jrem.Username,
			Password:   jrem.Password,
			PublicKey:  jrem.PublicKey,
			PrivateKey: jrem.PrivateKey,
		}
	}

	return &enclave{
		keystores: ks,
		creds:     rem,
		salt:      salt,
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
