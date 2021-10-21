package keystore

import (
	"fmt"

	"myst/pkg/crypto"

	"github.com/sanity-io/litter"
)

type Keystore struct {
	ID      string           `json:"id"`
	Version string           `json:"version"`
	Entries map[string]Entry `json:"entries"`
}

type Entry struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EncryptedKeystore struct {
	Keystore string `json:"domain"`
	Key      string `jsom:"key"`
}

// NewEncrypted creates and saves an encrypted domain from the given payload
func NewEncrypted(store, password string) (*EncryptedKeystore, error) {
	key, err := crypto.GenerateRandomBytes(32)
	if err != nil {
		return nil, err
	}
	fmt.Println(key)
	return &EncryptedKeystore{
		Keystore: store,
	}, nil
}

func (e *EncryptedKeystore) Save() {
	fmt.Println("save", litter.Sdump(e))
}
