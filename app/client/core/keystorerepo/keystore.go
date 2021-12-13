package keystorerepo

import (
	"encoding/json"

	"myst/app/client/core/domain/keystore/entry"

	"myst/app/client/core/domain/keystore"
)

type Keystore struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Version int     `json:"version"`
	Entries []Entry `json:"entries"`
}

type Entry struct {
	Id       string `json:"id"`
	Label    string `json:"label"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func KeystoreToJSON(k *keystore.Keystore) ([]byte, error) {
	entries := make([]Entry, len(k.Entries()))

	for i, e := range k.Entries() {
		entries[i] = Entry{
			Id:       e.Id(),
			Label:    e.Label(),
			Username: e.Username(),
			Password: e.Password(),
		}
	}

	return json.Marshal(
		Keystore{
			Id:      k.Id(),
			Name:    k.Name(),
			Version: k.Version(),
			Entries: entries,
		},
	)
}

func KeystoreFromJSON(b []byte) (*keystore.Keystore, error) {
	var k Keystore

	if err := json.Unmarshal(b, &k); err != nil {
		return nil, err
	}

	entries := make([]*entry.Entry, len(k.Entries))

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

		entries[i] = e
	}

	return keystore.New(
		keystore.WithId(k.Id),
		keystore.WithName(k.Name),
		keystore.WithVersion(k.Version),
		keystore.WithEntries(entries),
	)
}
