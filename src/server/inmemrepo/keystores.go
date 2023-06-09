package inmemrepo

import (
	"myst/src/server/application"
	"myst/src/server/application/domain/keystore"
)

func (r *Repository) CreateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.keystores[k.Id] = k

	return k, nil
}

func (r *Repository) Keystore(id string) (keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	k, ok := r.keystores[id]
	if !ok {
		return keystore.Keystore{}, application.ErrKeystoreNotFound
	}

	return k, nil
}

func (r *Repository) Keystores() ([]keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	keystores := make([]keystore.Keystore, 0, len(r.keystores))
	for _, k := range r.keystores {
		keystores = append(keystores, k)
	}

	return keystores, nil
}

func (r *Repository) UpdateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, ok := r.keystores[k.Id]
	if !ok {
		return keystore.Keystore{}, application.ErrKeystoreNotFound
	}

	r.keystores[k.Id] = k

	return k, nil
}

func (r *Repository) DeleteKeystore(id string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	delete(r.keystores, id)

	return nil
}
