package database

import (
	"bytes"
	"fmt"
	"github.com/natefinch/atomic"
	"os"
)

//store *keystore.Keystore
//key    []byte
//loaded bool
//mux sync.Mutex
//defaultKeystorePath = fmt.Sprintf("data/keystores/%s.mst", uuid.Nil.String())

var (
	ErrNotFound = fmt.Errorf("not found")
)

func Init() error {
	err := os.MkdirAll("data/users", os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll("data/keystores", os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func Load(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return b, nil
}

func Save(path string, b []byte) error {
	return atomic.WriteFile(path, bytes.NewReader(b))
}

//func Debug() keystore.Keystore {
//	mux.Lock()
//	defer mux.Unlock()
//	return *store
//}

//func LoadDefaultKeystore(pass string) error {
//
//	var err error
//	store, err = keystore.NewFromFile(defaultKeystorePath, pass)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func SaveDefaultKeystore(pass string) error {
//	if store == nil {
//		return ErrNotLoaded
//	}
//	return store.Save(defaultKeystorePath, pass)
//}
//

//func Load(pass string) (err error) {
//	mux.Lock()
//	defer mux.Unlock()
//	store, key, err = keystore.Load(defaultKeystorePath, pass)
//	if err != nil {
//		return err
//	}
//	loaded = true
//	return nil
//}
//
//func Save() error {
//	mux.Lock()
//	defer mux.Unlock()
//	if !loaded {
//		return ErrAccessDenied
//	}
//	return store.Save(defaultKeystorePath, key)
//}
//
//func AddSite(opts keystore.AddSiteOptions) error {
//	mux.Lock()
//	defer mux.Unlock()
//	if !loaded {
//		return ErrAccessDenied
//	}
//	store.AddSite(opts)
//	return store.Save(defaultKeystorePath, key)
//}
