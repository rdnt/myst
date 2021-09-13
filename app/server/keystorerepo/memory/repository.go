package keystorerepo

import (
	"fmt"
	"sync"

	"myst/app/server/domain/keystore"
)

type Repository struct {
	mux       sync.Mutex
	keystores map[string]keystore.Keystore
}

func (r *Repository) CreateKeystore(opts ...keystore.Option) (*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	s, err := keystore.New(opts...)
	if err != nil {
		return nil, err
	}

	_, ok := r.keystores[s.Id()]
	if ok {
		return nil, fmt.Errorf("already exists")
	}

	r.keystores[s.Id()] = *s
	return s, nil
}

func (r *Repository) Keystore(id string) (*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	s, ok := r.keystores[id]
	if !ok {
		return nil, keystore.ErrNotFound
	}

	return &s, nil
}

func (r *Repository) Keystores() ([]*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	sessions := make([]*keystore.Keystore, 0, len(r.keystores))
	for _, s := range r.keystores {
		sessions = append(sessions, &s)
	}

	return sessions, nil
}

func (r *Repository) UpdateKeystore(s *keystore.Keystore) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, ok := r.keystores[s.Id()]
	if !ok {
		return fmt.Errorf("not found")
	}

	r.keystores[s.Id()] = *s
	return nil
}

func (r *Repository) DeleteKeystore(id string) error {
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
