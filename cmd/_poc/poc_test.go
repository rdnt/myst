package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestKeystore(t *testing.T) {
	// shared encryption key between users
	keystoreKey := []byte("uKR6dCFtqBbj22mMCCTr1LioReGeBq6W") // 32 bytes

	key := []byte("{\"id\":\"keATfB8JP2ggT7U9JZrpV9\"}")
	b := key

	enc, mac, err := encryptKeystore(b, keystoreKey)
	if err != nil {
		panic(err)
	}

	// [enc, mac is sent to the server]
	btest, _ := json.Marshal(append(enc, mac...))
	fmt.Println("store", string(btest))

	b, err = decryptKeystore(enc, keystoreKey, mac)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, string(key), string(b))
}

// encoding/decoding for this is run on the client
func TestKey(t *testing.T) {
	// shared encryption key between users
	masterPassword := []byte("uNoSp5bRncRK3m1ElSnRxQoLKDWfAvCV")
	salt := []byte("7dY4fnrxUfpZM6LH")

	// generated on the client
	masterKey := crypto2.Argon2Id(masterPassword, salt)

	// generated first time on the client then stored on server
	key := []byte("uKR6dCFtqBbj22mMCCTr1LioReGeBq6W") // 32 bytes

	enc, mac, err := encryptKey(key, masterKey)
	if err != nil {
		panic(err)
	}

	// [enc,mac] is sent to server (encrypted key)
	btest, _ := json.Marshal(append(enc, mac...))
	fmt.Println("key", string(btest))

	b, err := decryptKey(enc, masterKey, mac)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, string(key), string(b))
}
