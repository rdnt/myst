package repository

import (
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/keystore/entry"
)

type Enclave struct {
	Keystores map[string]Keystore `json:"keystores"`
	Keys      map[string][]byte   `json:"keys"`
	Remote    *Remote             `json:"remote,omitempty"`
}

type Remote struct {
	Address    string `json:"address"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PublicKey  []byte `json:"publicKey"`
	PrivateKey []byte `json:"privateKey"`
}

type Keystore struct {
	Id       string      `json:"id"`
	RemoteId string      `json:"remoteId"`
	Name     string      `json:"name"`
	Version  int         `json:"version"`
	Entries  []JSONEntry `json:"entries"`
}

type JSONEntry struct {
	Id       string `json:"id"`
	Website  string `json:"website"`
	Username string `json:"username"`
	Password string `json:"password"`
	Notes    string `json:"notes"`
}

func KeystoreToJSON(k keystore.Keystore) Keystore {
	entries := []JSONEntry{}

	for _, e := range k.Entries {
		entries = append(entries, JSONEntry{
			Id:       e.Id,
			Website:  e.Website,
			Username: e.Username,
			Password: e.Password,
			Notes:    e.Notes,
		})
	}

	return Keystore{
		Id:       k.Id,
		RemoteId: k.RemoteId,
		Name:     k.Name,
		Version:  k.Version,
		Entries:  entries,
	}
}

func KeystoreFromJSON(k Keystore) (keystore.Keystore, error) {
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
		keystore.WithName(k.Name),
		keystore.WithVersion(k.Version),
		keystore.WithEntries(entries),
	), nil
}
