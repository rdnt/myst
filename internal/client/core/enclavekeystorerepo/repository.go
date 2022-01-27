package enclavekeystorerepo

import (
	"encoding/json"
	"fmt"

	"myst/internal/client/core/domain/keystore"
	"myst/internal/client/core/enclaverepo"
	"myst/internal/client/core/keystorerepo/jsonkeystore"
	"myst/pkg/crypto"
)

type Repository struct {
	enclaveRepo *enclaverepo.Repository
}

func (r *Repository) Initialize(password, keystoreName string) (*keystore.Keystore, error) {
	k := keystore.New(keystore.WithName(keystoreName))

	jk := jsonkeystore.Marshal(k)

	p := crypto.DefaultArgon2IdParams

	key, err := crypto.GenerateRandomBytes(uint(p.KeyLength))
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(jk)
	if err != nil {
		return nil, err
	}

	b, err = crypto.AES256CBC_Encrypt(key, b)
	if err != nil {
		return nil, err
	}

	err = r.enclaveRepo.Create(password)
	if err != nil {
		return nil, err
	}

	e, err := r.enclaveRepo.Enclave()
	if err != nil {
		return nil, err
	}

	e.AddKeystore(k)

	return k, nil
}

func (r *Repository) Create(opts ...keystore.Option) (*keystore.Keystore, error) {
	k := keystore.New(opts...)

	jk := jsonkeystore.Marshal(k)

	p := crypto.DefaultArgon2IdParams

	key, err := crypto.GenerateRandomBytes(uint(p.KeyLength))
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(jk)
	if err != nil {
		return nil, err
	}

	b, err = crypto.AES256CBC_Encrypt(key, b)
	if err != nil {
		return nil, err
	}

	enclaveExists := false

	if !enclaveExists {
		err = r.enclaveRepo.Create(password)
		if err != nil {
			return nil, err
		}
	} else {
		err = r.enclaveRepo.Create(password)
		if err != nil {
			return nil, err
		}
	}

	return k, nil
}

func (r *Repository) Authenticate(password string) error {
	return r.enclaveRepo.Authenticate(password)
}

func (r *Repository) Keystore(id string) (*keystore.Keystore, error) {
	e, err := r.enclaveRepo.Enclave()
	if err != nil {
		return nil, err
	}

	ks := e.Keystores()

	k, ok := ks[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	return k, nil
}

func (r *Repository) Update(k *keystore.Keystore) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) KeystoreIds() ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) Keystores() ([]*keystore.Keystore, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}

func New() (keystore.Repository, error) {
	erepo, err := enclaverepo.New("data")
	if err != nil {
		return nil, err
	}

	return &Repository{
		enclaveRepo: erepo,
	}, nil
}
