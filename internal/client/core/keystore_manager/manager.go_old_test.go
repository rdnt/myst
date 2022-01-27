package keystore_manager_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"myst/internal/client/core/domain/keystore"
	"myst/internal/client/core/keystore_manager"
)

func TestCreate(t *testing.T) {
	password := "12345678"

	m := keystore_manager.New()

	k, err := m.Create(
		keystore.WithName("passwords"),
		keystore.WithPassword(password),
	)
	assert.Nil(t, err)

	fmt.Println(k)

	err = m.Authenticate(password)
	assert.Nil(t, err)
}
