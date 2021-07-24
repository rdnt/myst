package keystore

import "fmt"

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
	Keystore string `json:"keystore"`
	Key      string `jsom:"key"`
}

// NewEncrypted creates and saves an encrypted keystore from the given payload
func NewEncrypted(payload, password string) *EncryptedKeystore {
	enc := EncryptedKeystore(payload)
	return &enc
}

func (e *EncryptedKeystore) Save() {
	fmt.Println("save", e)
}
