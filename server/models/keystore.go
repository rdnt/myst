package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/natefinch/atomic"
	"myst/server/crypto"
	"myst/server/database"
	"myst/server/util"
	"os"
	"path/filepath"
	"time"
)

var (
	ErrAuthenticationFailed = fmt.Errorf("authentication failed")
	argon2IdParams          = crypto.DefaultArgon2IdParams
)

// Keystore holds
type Keystore struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Version   string           `json:"version"`
	Entries   map[string]Entry `json:"entries"`
	CreatedAt Timestamp        `json:"created_at"`
	UpdatedAt Timestamp        `json:"updated_at"`
}

// Entry represents a site entry and its associated data
type Entry struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// NewKeystore creates a new keystore and saves it
func NewKeystore(name, passphrase string) (*Keystore, error) {
	k := &Keystore{
		Name:    name,
		Version: "1",
	}
	err := k.Save(passphrase)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Save saves the user on the database
func (k *Keystore) Save(passphrase string) error {
	now := Timestamp{time.Now()}
	if k.ID == "" {
		k.ID = util.NewUUID()
		k.CreatedAt = now
	}
	if k.Entries == nil {
		k.Entries = make(map[string]Entry)
	}
	k.UpdatedAt = now

	salt, err := crypto.GenerateRandomSalt()
	if err != nil {
		return err
	}

	key := crypto.Argon2Id(passphrase, salt)

	k.save(path, key)

	//b, err := json.Marshal(k)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (k *Keystore) save(path string, key []byte) (err error) {
	err = database.Save(fmt.Sprintf("data/keystores/%s.json", k.ID), b)
	if err != nil {
		return err
	}

	// Encode to json
	b, err := json.Marshal(k)
	if err != nil {
		return err
	}

	// Encrypt keystore
	b, err = crypto.AES256CBC_Encrypt(key, b)
	if err != nil {
		return err
	}

	// Authenticate ciphertext
	mac := crypto.HMAC_SHA256(key, b)
	b = append(mac, b...)
	//b = append(store.salt, b...)

	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(absPath), os.ModePerm)
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)
	err = atomic.WriteFile(absPath, r)
	if err != nil {
		return err
	}

	return nil
}
