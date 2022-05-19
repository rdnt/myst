package keystorerepo

import (
	"myst/internal/client/application/domain/entry"
	"myst/internal/client/application/domain/keystore"
)

type JSONEnclave struct {
	Keystores map[string]JSONKeystore `json:"keystores"`
	Keys      map[string][]byte       `json:"keys"`
}

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

func KeystoreToJSON(k keystore.Keystore) JSONKeystore {
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

	return JSONKeystore{
		Id:      k.Id,
		Name:    k.Name,
		Version: k.Version,
		Entries: entries,
	}
}

func KeystoreFromJSON(k JSONKeystore) (keystore.Keystore, error) {
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
		keystore.WithName(k.Name),
		keystore.WithVersion(k.Version),
		keystore.WithEntries(entries),
	), nil
}
