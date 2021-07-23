package keystore

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/natefinch/atomic"
	"myst/server/crypto"
	"os"
	"path/filepath"
)

var (
	ErrAuthFailed  = fmt.Errorf("authentication failed")
	argon2idParams = crypto.DefaultArgon2IdParams
)

// Keystore holds
type Keystore struct {
	ID      string           `json:"id"`
	Version string           `json:"version"`
	Entries map[string]Entry `json:"entries"`
}

// Entry represents a website and its saved data
type Entry struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func New(path, pass string) (store *Keystore, key []byte, err error) {
	salt, err := crypto.GenerateRandomBytes(uint(argon2idParams.SaltLength))
	if err != nil {
		return nil, nil, err
	}

	key = crypto.Argon2Id(pass, salt)

	store = &Keystore{
		ID:      uuid.New().String(),
		Entries: make(map[string]Entry, 0),
	}

	err = store.Save(path, key)
	if err != nil {
		return nil, nil, err
	}

	return store, key, nil
}

func Load(path, pass string) (store *Keystore, key []byte, err error) {
	b, err := os.ReadFile(path)

	if os.IsNotExist(err) {
		return New(path, pass)
	} else if err != nil {
		return nil, nil, err
	}

	salt := b[:argon2idParams.SaltLength]
	mac := b[argon2idParams.SaltLength : sha256.Size+argon2idParams.SaltLength]
	b = b[argon2idParams.SaltLength+sha256.Size:]

	key = crypto.Argon2Id(pass, salt)

	valid := crypto.VerifyHMAC_SHA256(key, mac, b)
	if !valid {
		return nil, nil, ErrAuthFailed
	}

	// Decrypt keystore
	b, err = crypto.AES256CBC_Decrypt(key, b)
	if err != nil {
		return nil, nil, err
	}
	// Decode decrypted keystore from json
	err = json.Unmarshal(b, &store)
	if err != nil {
		return nil, nil, err
	}
	//store.salt = salt
	// Return the keystore
	return store, key, nil

}

func (store *Keystore) Save(path string, key []byte) (err error) {

	// Encode to json
	b, err := json.Marshal(store)
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

type AddSiteOptions struct {
	Label    string
	Email    string
	Password string
}

func (store *Keystore) AddSite(opts AddSiteOptions) {
	//store.Lock()
	//defer store.Unlock()
	//store.Entries[uuid.NewString()] = &Entry{
	//	Label:    opts.Label,
	//	Email:    opts.Email,
	//	Password: opts.Password,
	//}
}
