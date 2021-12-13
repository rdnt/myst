package keystorerepo

import (
	"encoding/json"
	"fmt"
	"sync"

	"myst/pkg/crypto"

	"myst/app/client/core/domain/keystore"
)

type Repository struct {
	mux       sync.Mutex
	keystores map[string][]byte
	keys      map[string][]byte
}

func (r *Repository) Create(opts ...keystore.Option) (*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	k, err := keystore.New(opts...)
	if err != nil {
		return nil, err
	}

	_, ok := r.keystores[k.Id()]
	if ok {
		return nil, fmt.Errorf("already exists")
	}

	b, err := json.Marshal(k)
	if err != nil {
		return nil, err
	}

	p := crypto.DefaultArgon2IdParams

	salt, err := crypto.GenerateRandomBytes(uint(p.SaltLength))
	if err != nil {
		return nil, err
	}

	key := crypto.Argon2Id(k.Passphrase(), salt)

	b, err = Encrypt(b, key)
	if err != nil {
		return nil, err
	}

	r.keystores[k.Id()] = b
	r.keys[k.Id()] = key

	return k, nil
}

func (r *Repository) Authenticate(id string, passphrase []byte) error {
	b, ok := r.keystores[id]
	if !ok {
		return keystore.ErrNotFound
	}

	salt, err := GetSaltFromData(b)
	if err != nil {
		return err
	}

	key := crypto.Argon2Id(passphrase, salt)

	r.keys[id] = key
	return nil
}

func (r *Repository) Keystore(id string) (*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	b, ok := r.keystores[id]
	if !ok {
		return nil, keystore.ErrNotFound
	}

	key, ok := r.keys[id]
	if !ok {
		return nil, fmt.Errorf("authentication required")
	}

	b, err := Decrypt(b, key)
	if err != nil {
		return nil, err
	}

	k := &keystore.Keystore{}

	err = json.Unmarshal(b, &k)
	if err != nil {
		return nil, err
	}

	return k, nil
}

func (r *Repository) Keystores() ([]*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	return nil, nil

	//keystores := make([]*keystore.Keystore, 0, len(r.keystores))
	//for _, k := range r.keystores {
	//	keystores = append(keystores, &k)
	//}
	//
	//return keystores, nil
}

func (r *Repository) Update(s *keystore.Keystore) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	//_, ok := r.keystores[s.Id()]
	//if !ok {
	//	return fmt.Errorf("not found")
	//}
	//
	//r.keystores[s.Id()] = *s
	return nil
}

func (r *Repository) Delete(id string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	delete(r.keystores, id)
	delete(r.keys, id)

	return nil
}

func New() keystore.Repository {
	return &Repository{
		keystores: make(map[string][]byte),
		keys:      make(map[string][]byte),
	}
}
