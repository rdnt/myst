package keystorerepo

import (
	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/application/domain/keystore/entry"
)

type JSONKeystore struct {
	Id      string      `json:"id"`
	Name    string      `json:"name"`
	Version int         `json:"version"`
	Entries []JSONEntry `json:"entries"`
}

type JSONEntry struct {
	Id       string `json:"id"`
	Website  string `json:"website"`
	Username string `json:"username"`
	Password string `json:"password"`
	Notes    string `json:"notes"`
}

func ToJSONKeystore(k *keystore.Keystore) JSONKeystore {
	entries := make([]JSONEntry, len(k.Entries()))

	for i, e := range k.Entries() {
		entries[i] = JSONEntry{
			Id:       e.Id(),
			Website:  e.Website(),
			Username: e.Username(),
			Password: e.Password(),
			Notes:    e.Notes(),
		}
	}

	return JSONKeystore{
		Id:      k.Id(),
		Name:    k.Name(),
		Version: k.Version(),
		Entries: entries,
	}
}

func ToKeystore(k JSONKeystore) (*keystore.Keystore, error) {
	entries := make([]entry.Entry, len(k.Entries))

	for i, e := range k.Entries {
		e := entry.New(
			entry.WithId(e.Id),
			entry.WithWebsite(e.Website),
			entry.WithUsername(e.Username),
			entry.WithPassword(e.Password),
			entry.WithNotes(e.Notes),
		)

		entries[i] = e
	}

	return keystore.New(
		keystore.WithId(k.Id),
		keystore.WithName(k.Name),
		keystore.WithVersion(k.Version),
		keystore.WithEntries(entries),
	), nil
}
