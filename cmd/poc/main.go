package main

import (
	"fmt"
	"myst/pkg/crypto"
)

var (
	ErrAuthenticationFailed = fmt.Errorf("authentication failed")
)

func encryptKeystore(b, key []byte) ([]byte, []byte, error) {
	// Encrypt keystore
	b, err := crypto.AES256CBC_Encrypt(key, b)
	if err != nil {
		return nil, nil, err
	}

	// Authenticate ciphertext
	mac := crypto.HMAC_SHA256(key, b)

	return b, mac, nil
}

func decryptKeystore(b, key, mac []byte) ([]byte, error) {
	valid := crypto.VerifyHMAC_SHA256(key, mac, b)
	if !valid {
		return nil, ErrAuthenticationFailed
	}

	// Decrypt keystore
	b, err := crypto.AES256CBC_Decrypt(key, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func encryptKey(b, key []byte) ([]byte, []byte, error) {
	// Encrypt keystore
	b, err := crypto.AES256CBC_Encrypt(key, b)
	if err != nil {
		return nil, nil, err
	}

	// Authenticate ciphertext
	mac := crypto.HMAC_SHA256(key, b)

	return b, mac, nil
}

func decryptKey(b, key, mac []byte) ([]byte, error) {
	valid := crypto.VerifyHMAC_SHA256(key, mac, b)
	if !valid {
		return nil, ErrAuthenticationFailed
	}

	// Decrypt keystore
	b, err := crypto.AES256CBC_Decrypt(key, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
