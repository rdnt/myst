package keystorerepo

import (
	"fmt"
	"sync"

	"myst/app/client/core/keyrepo"

	"myst/pkg/crypto"

	"myst/app/client/core/domain/keystore"
)

type repository struct {
	mux       sync.Mutex
	keystores map[string][]byte
	keyRepo   *keyrepo.Repository
}

func (r *repository) Create(opts ...keystore.Option) (*keystore.Keystore, error) {
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

	b, err := KeystoreToJSON(k)
	if err != nil {
		return nil, err
	}

	p := crypto.DefaultArgon2IdParams

	salt, err := crypto.GenerateRandomBytes(uint(p.SaltLength))
	if err != nil {
		return nil, err
	}

	key := crypto.Argon2Id([]byte(k.Passphrase()), salt)

	b, err = Encrypt(b, key, salt)
	if err != nil {
		return nil, err
	}

	r.keystores[k.Id()] = b

	r.keyRepo.Set(k.Id(), key)

	return k, nil
}

func (r *repository) Unlock(id string, passphrase string) (*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	b, ok := r.keystores[id]
	if !ok {
		return nil, keystore.ErrNotFound
	}

	salt, err := GetSaltFromData(b)
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
	b, ok := r.keystores[id]
	if !ok {
		return nil, keystore.ErrNotFound
	}

	key, err := r.keyRepo.Key(id)
	if err != nil {
		return nil, fmt.Errorf("authentication required")
	}

	b, err = Decrypt(b, key)
	if err != nil {
		return nil, err
	}

	k, err := KeystoreFromJSON(b)
	if err != nil {
		return nil, err
	}

	return k, nil
}

func (r *repository) Keystores() ([]*keystore.Keystore, error) {
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

func (r *repository) Update(s *keystore.Keystore) error {
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

func (r *repository) Delete(id string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	delete(r.keystores, id)

	r.keyRepo.Delete(id)

	return nil
}

func New() *repository {
	return &repository{
		keystores: map[string][]byte{},
		keyRepo:   keyrepo.New(),
	}
}
