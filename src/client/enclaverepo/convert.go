package enclaverepo

import (
	"time"

	"myst/src/client/application/domain/entry"
	"myst/src/client/application/domain/keystore"
)

type EnclaveJSON struct {
	Keystores map[string]KeystoreJSON `json:"keystores"`
	Keys      map[string][]byte       `json:"keys"`
	Remote    *RemoteJSON             `json:"remote,omitempty"`
}

type RemoteJSON struct {
	Address    string `json:"address"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PublicKey  []byte `json:"publicKey"`
	PrivateKey []byte `json:"privateKey"`
}

type KeystoreJSON struct {
	Id        string      `json:"id"`
	RemoteId  string      `json:"remoteId"`
	Name      string      `json:"name"`
	Version   int         `json:"version"`
	Entries   []EntryJSON `json:"entries"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type EntryJSON struct {
	Id        string    `json:"id"`
	Website   string    `json:"website"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func KeystoreToJSON(k keystore.Keystore) KeystoreJSON {
	entries := []EntryJSON{}

	for _, e := range k.Entries {
		entries = append(entries, EntryJSON{
			Id:        e.Id,
			Website:   e.Website,
			Username:  e.Username,
			Password:  e.Password,
			Notes:     e.Notes,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		})
	}

	return KeystoreJSON{
		Id:        k.Id,
		RemoteId:  k.RemoteId,
		Name:      k.Name,
		Version:   k.Version,
		Entries:   entries,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

func KeystoreFromJSON(k KeystoreJSON) (keystore.Keystore, error) {
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
	}, nil
}
