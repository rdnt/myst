package keystorerepo

import (
	"fmt"
	"sync"

	"myst/internal/server/core/domain/keystore"
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

func (r *Repository) UserKeystores(userId string) ([]*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	ks := []*keystore.Keystore{}

	for _, k := range r.keystores {
		if k.OwnerId() == userId {
			ks = append(ks, &k)
		} else {
			for _, uid := range k.ViewerIds() {
				if uid == userId {
					ks = append(ks, &k)
				}
			}
		}
	}

	return ks, nil
}

func (r *Repository) UserKeystore(userId, keystoreId string) (*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, k := range r.keystores {
		if k.Id() == keystoreId {
			if k.OwnerId() == userId {
				return &k, nil
			}

			for _, uid := range k.ViewerIds() {
				if uid == userId {
					return &k, nil
				}
			}
		}
	}

	return nil, keystore.ErrNotFound
}

func New() keystore.Repository {
	return &Repository{
		keystores: make(map[string]keystore.Keystore),
	}
}
