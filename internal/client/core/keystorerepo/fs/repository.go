package keystorerepo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"myst/internal/client/core/domain/keystore"
	jsonkeystore "myst/internal/client/core/keystorerepo/keystore"
	"myst/internal/client/core/sessionrepo"

	"myst/pkg/crypto"
	"myst/pkg/enclave"
)

var (
	ErrAuthenticationFailed = enclave.ErrAuthenticationFailed
)

type repository struct {
	mux     sync.Mutex
	keyRepo *sessionrepo.Repository
}

func (r *repository) Create(opts ...keystore.Option) (*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	k, err := keystore.New(opts...)
	if err != nil {
		return nil, err
	}

	kpath := "data/keystores/" + k.Id() + ".mst"

	if _, err := os.Stat(kpath); err == nil {
		return nil, fmt.Errorf("already exists")

	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	b, err := jsonkeystore.Marshal(k)
	if err != nil {
		return nil, err
	}

	p := crypto.DefaultArgon2IdParams

	salt, err := crypto.GenerateRandomBytes(uint(p.SaltLength))
	if err != nil {
		return nil, err
	}

	key := crypto.Argon2Id([]byte(k.Passphrase()), salt)

	b, err = enclave.Encrypt(b, key, salt)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(kpath, b, 0600)
	if err != nil {
		return nil, err
	}

	r.keyRepo.Set(k.Id(), key)

	return k, nil
}

func (r *repository) Unlock(id string, passphrase string) (*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	kpath := "data/keystores/" + id + ".mst"

	if _, err := os.Stat(kpath); errors.Is(err, os.ErrNotExist) {
		return nil, keystore.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(kpath)
	if err != nil {
		return nil, err
	}

	salt, err := enclave.GetSaltFromData(b)
	if err != nil {
		return nil, err
	}

	key := crypto.Argon2Id([]byte(passphrase), salt)

	r.keyRepo.Set(id, key)

	return r.keystore(id)
}

func (r *repository) HealthCheck() {
	r.keyRepo.HealthCheck()
}

func (r *repository) Keystore(id string) (*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.keystore(id)
}

func (r *repository) keystore(id string) (*keystore.Keystore, error) {
	kpath := "data/keystores/" + id + ".mst"

	if _, err := os.Stat(kpath); errors.Is(err, os.ErrNotExist) {
		return nil, keystore.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(kpath)
	if err != nil {
		return nil, err
	}

	key, err := r.keyRepo.Key(id)
	if err != nil {
		return nil, keystore.ErrAuthenticationRequired
	}

	b, err = enclave.Decrypt(b, key)
	if errors.Is(err, enclave.ErrAuthenticationFailed) {
		return nil, ErrAuthenticationFailed
	} else if err != nil {
		return nil, err
	}

	k, err := jsonkeystore.Unmarshal(b)
	if err != nil {
		return nil, err
	}

	return k, nil
}

func (r *repository) Keystores() ([]*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	fis, err := os.ReadDir("data/keystores")
	if err != nil {
		return nil, err
	}

	ks := []*keystore.Keystore{}

	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}

		if filepath.Ext(fi.Name()) != ".mst" {
			continue
		}

		id := strings.TrimSuffix(filepath.Base(fi.Name()), filepath.Ext(".mst"))

		k, err := r.keystore(id)
		if err != nil {
			return nil, err
		}

		ks = append(ks, k)
	}

	return ks, nil
}

func (r *repository) Update(k *keystore.Keystore) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	key, err := r.keyRepo.Key(k.Id())
	if err != nil {
		return keystore.ErrAuthenticationRequired
	}

	kpath := "data/keystores/" + k.Id() + ".mst"

	if _, err := os.Stat(kpath); errors.Is(err, os.ErrNotExist) {
		return keystore.ErrNotFound
	} else if err != nil {
		return err
	}

	// read the existing keystore and get the salt; we will reuse the salt so that we don't have to re-generate the
	// encryption key
	b, err := os.ReadFile(kpath)
	if err != nil {
		return err
	}

	salt, err := enclave.GetSaltFromData(b)
	if err != nil {
		return err
	}

	b, err = jsonkeystore.Marshal(k)
	if err != nil {
		return err
	}

	b, err = enclave.Encrypt(b, key, salt)
	if err != nil {
		return err
	}

	err = os.WriteFile(kpath, b, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(id string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	kpath := "data/keystores/" + id + ".mst"

	if _, err := os.Stat(kpath); errors.Is(err, os.ErrNotExist) {
		return keystore.ErrNotFound
	} else if err != nil {
		return err
	}

	err := os.Remove(kpath)
	if err != nil {
		return err
	}

	r.keyRepo.Delete(id)

	return nil
}

func New() *repository {
	return &repository{
		keyRepo: sessionrepo.New(),
	}
}
