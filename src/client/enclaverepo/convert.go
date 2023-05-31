package enclaverepo

import (
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/keystore/entry"
)

type EnclaveJSON struct {
	Keystores map[string]KeystoreJSON `json:"keystores"`
	Keys      map[string][]byte       `json:"keys"`
	Remote    *RemoteJSON             `json:"remote,omitempty"`
}

type RemoteJSON struct {
	Address    string            `json:"address"`
	Username   string            `json:"username"`
	Password   string            `json:"password"`
	PublicKey  []byte            `json:"publicKey"`
	PrivateKey []byte            `json:"privateKey"`
	UserKeys   map[string][]byte `json:"userKeys"`
}

type KeystoreJSON struct {
	Id       string      `json:"id"`
	RemoteId string      `json:"remoteId"`
	ReadOnly bool        `json:"readOnly"`
	Name     string      `json:"name"`
	Version  int         `json:"version"`
	Entries  []EntryJSON `json:"entries"`
}

type EntryJSON struct {
	Id       string `json:"id"`
	Website  string `json:"website"`
	Username string `json:"username"`
	Password string `json:"password"`
	Notes    string `json:"notes"`
}

func KeystoreToJSON(k keystore.Keystore) KeystoreJSON {
	entries := []EntryJSON{}

	for _, e := range k.Entries {
		entries = append(entries, EntryJSON{
			Id:       e.Id,
			Website:  e.Website,
			Username: e.Username,
			Password: e.Password,
			Notes:    e.Notes,
		})
	}

	return KeystoreJSON{
		Id:       k.Id,
		RemoteId: k.RemoteId,
		ReadOnly: k.ReadOnly,
		Name:     k.Name,
		Version:  k.Version,
		Entries:  entries,
	}
}

func KeystoreFromJSON(k KeystoreJSON) (keystore.Keystore, error) {
	entries := make(map[string]entry.Entry, len(k.Entries))

	for _, e := range k.Entries {
		e := entry.New(
			entry.WithId(e.Id),
			entry.WithWebsite(e.Website),
			entry.WithUsername(e.Username),
			entry.WithPassword(e.Password),
			entry.WithNotes(e.Notes),
		)

		entries[e.Id] = e
	}

	return keystore.New(
		keystore.WithId(k.Id),
		keystore.WithRemoteId(k.RemoteId),
		keystore.WithReadOnly(k.ReadOnly),
		keystore.WithName(k.Name),
		keystore.WithVersion(k.Version),
		keystore.WithEntries(entries),
	), nil
}
