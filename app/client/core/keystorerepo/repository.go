package keystorerepo

import (
	"fmt"
	"sync"

	"myst/app/client/core/domain/keystore"
)

type Repository struct {
	mux       sync.Mutex
	keystores map[string]keystore.Keystore
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

	r.keystores[k.Id()] = *k

	return k, nil
}

func (r *Repository) Keystore(id string) (*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	k, ok := r.keystores[id]
	if !ok {
		return nil, keystore.ErrNotFound
	}

	return &k, nil
}

func (r *Repository) Keystores() ([]*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	keystores := make([]*keystore.Keystore, 0, len(r.keystores))
	for _, k := range r.keystores {
		keystores = append(keystores, &k)
	}

	return keystores, nil
}

func (r *Repository) Update(s *keystore.Keystore) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, ok := r.keystores[s.Id()]
	if !ok {
		return fmt.Errorf("not found")
	}

	r.keystores[s.Id()] = *s
	return nil
}

func (r *Repository) Delete(id string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	delete(r.keystores, id)
	return nil
}

func New() keystore.Repository {
	return &Repository{
		keystores: make(map[string]keystore.Keystore),
	}
}
