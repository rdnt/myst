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
	Label    string `json:"label"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func ToJSONKeystore(k *keystore.Keystore) JSONKeystore {
	entries := make([]JSONEntry, len(k.Entries()))

	for i, e := range k.Entries() {
		entries[i] = JSONEntry{
			Id:       e.Id(),
			Label:    e.Label(),
			Username: e.Username(),
			Password: e.Password(),
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
		e, err := entry.New(
			entry.WithId(e.Id),
			entry.WithUsername(e.Username),
			entry.WithPassword(e.Password),
			entry.WithLabel(e.Label),
		)
		if err != nil {
			return nil, err
		}

		entries[i] = *e
	}

	return keystore.New(
		keystore.WithId(k.Id),
		keystore.WithName(k.Name),
		keystore.WithVersion(k.Version),
		keystore.WithEntries(entries),
	), nil
}
