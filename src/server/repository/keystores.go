package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"myst/src/server/application/domain/keystore"
)

func (r *Repository) CreateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.keystores[k.Id] = k

	{ // debug
		b, _ := json.Marshal(r.keystores)
		_ = os.WriteFile("keystores.json", b, 0666)
	}

	return k, nil
}

func (r *Repository) Keystore(id string) (keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	k, ok := r.keystores[id]
	if !ok {
		return keystore.Keystore{}, keystore.ErrNotFound
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
		return keystore.Keystore{}, fmt.Errorf("not found")
	}

	r.keystores[k.Id] = k

	{ // debug
		b, _ := json.Marshal(r.keystores)
		_ = os.WriteFile("keystores.json", b, 0666)
	}

	return k, nil
}

func (r *Repository) DeleteKeystore(id string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	delete(r.keystores, id)

	{ // debug
		b, _ := json.Marshal(r.keystores)
		_ = os.WriteFile("keystores.json", b, 0666)
	}

	return nil
}

func (r *Repository) UserKeystores(userId string) ([]keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	ks := []keystore.Keystore{}

	for _, k := range r.keystores {
		if k.OwnerId == userId {
			ks = append(ks, k)
		} else {
			for _, uid := range k.ViewerIds {
				if uid == userId {
					ks = append(ks, k)
				}
			}
		}
	}

	return ks, nil
}

func (r *Repository) UserKeystore(userId, keystoreId string) (keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, k := range r.keystores {
		if k.Id == keystoreId {
			if k.OwnerId == userId {
				return k, nil
			}

			for _, uid := range k.ViewerIds {
				if uid == userId {
					return k, nil
				}
			}
		}
	}

	return keystore.Keystore{}, keystore.ErrNotFound
}
