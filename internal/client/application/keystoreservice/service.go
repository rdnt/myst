package keystoreservice

import (
	"strings"

	"myst/internal/client/application/domain/enclave"
	"myst/internal/client/application/domain/entry"
	"myst/internal/client/application/domain/keystore"
	"myst/pkg/logger"

	"github.com/pkg/errors"
)

var (
	ErrInvalidKeystoreRepository = errors.New("invalid keystore repository")
	ErrAuthenticationRequired    = errors.New("authentication required")
	ErrAuthenticationFailed      = errors.New("authentication failed")
	ErrInitializationRequired    = errors.New("initialization required")
	ErrInvalidPassword           = errors.New("invalid password")
	ErrEntryNotFound             = errors.New("entry not found")
)

type KeystoreRepository interface {
	keystore.Repository
	Authenticate(password string) error
	HealthCheck()
	CreateEnclave(password string) error
	Enclave() error
	SetRemote(address, username, password string, publicKey, privateKey []byte) error
	Remote() (enclave.Remote, error)
}

type service struct {
	keystores KeystoreRepository
}

func New(opts ...Option) (enclave.KeystoreService, error) {
	s := &service{}

	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	if s.keystores == nil {
		return nil, ErrInvalidKeystoreRepository
	}

	return s, nil
}

func (s *service) KeystoreEntries(id string) (map[string]entry.Entry, error) {
	k, err := s.keystores.Keystore(id)
	if err != nil {
		return nil, err
	}

	return k.Entries, nil
}

func (s *service) UpdateKeystoreEntry(keystoreId, entryId string, password, notes *string) (entry.Entry, error) {
	// do not allow empty password
	if password != nil && strings.TrimSpace(*password) == "" {
		return entry.Entry{}, ErrInvalidPassword
	}

	k, err := s.keystores.Keystore(keystoreId)
	if err != nil {
		return entry.Entry{}, err
	}

	entries := k.Entries

	e, ok := entries[entryId]
	if !ok {
		return entry.Entry{}, ErrEntryNotFound
	}

	if password != nil {
		e.SetPassword(*password)
	}

	if notes != nil {
		e.SetNotes(*notes)
	}

	entries[e.Id] = e

	k.Entries = entries

	return e, s.UpdateKeystore(k)
}

func (s *service) DeleteKeystoreEntry(keystoreId, entryId string) error {
	k, err := s.keystores.Keystore(keystoreId)
	if err != nil {
		return err
	}

	entries := k.Entries

	if _, ok := entries[entryId]; !ok {
		return ErrEntryNotFound
	}

	delete(entries, entryId)
	k.Entries = entries

	return s.UpdateKeystore(k)
}

func (s *service) CreateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	return s.keystores.CreateKeystore(k)
}

func (s *service) DeleteKeystore(id string) error {
	return s.keystores.DeleteKeystore(id)
}

func (s *service) CreateKeystoreEntry(keystoreId string, opts ...entry.Option) (entry.Entry, error) {
	k, err := s.keystores.Keystore(keystoreId)
	if err != nil {
		return entry.Entry{}, err
	}

	e := entry.New(opts...)

	entries := k.Entries
	entries[e.Id] = e
	k.Entries = entries

	return e, s.keystores.UpdateKeystore(k)
}

func (s *service) CreateFirstKeystore(k keystore.Keystore, password string) (keystore.Keystore, error) {
	err := s.keystores.CreateEnclave(password)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to initialize enclave")
	}

	k, err = s.keystores.CreateKeystore(k)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to create keystore")
	}

	logger.Printf("@@@@@@@@@@############### %#v\n", k)
	return k, nil
}

func (s *service) CreateEnclave(password string) error {
	return s.keystores.CreateEnclave(password)
}

func (s *service) Enclave() error {
	return s.keystores.Enclave()
}

func (s *service) Keystore(id string) (keystore.Keystore, error) {
	k, err := s.keystores.Keystore(id)
	//if errors.Is(err, keystore.ErrAuthenticationRequired) {
	//	return nil, ErrAuthenticationRequired
	//}

	return k, err
}

func (s *service) KeystoreByRemoteId(id string) (keystore.Keystore, error) {
	ks, err := s.keystores.Keystores()
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get keystores")
	}

	for _, k2 := range ks {
		if k2.RemoteId == id {
			return k2, nil
		}
	}

	return keystore.Keystore{}, errors.New("keystore not found")
}

func (s *service) Keystores() (map[string]keystore.Keystore, error) {
	ks, err := s.keystores.Keystores()
	if err != nil {
		return nil, err
	}

	//	return nil, ErrAuthenticationRequired
	//} else if err == keystore.ErrInitializationRequired {
	//	return nil, ErrInitializationRequired
	//}

	return ks, err
}

func (s *service) UpdateKeystore(k keystore.Keystore) error {
	err := s.keystores.UpdateKeystore(k)
	//if errors.Is(err, keystore.ErrAuthenticationRequired) {
	//	return ErrAuthenticationRequired
	//}

	return err
}

func (s *service) Authenticate(password string) error {
	return s.keystores.Authenticate(password)
}

func (s *service) HealthCheck() {
	s.keystores.HealthCheck()
}

func (s *service) SetRemote(address, username, password string, publicKey, privateKey []byte) error {
	return s.keystores.SetRemote(address, username, password, publicKey, privateKey)
}

func (s *service) Remote() (enclave.Remote, error) {
	return s.keystores.Remote()
}
