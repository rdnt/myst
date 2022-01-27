package enclaverepo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"myst/internal/client/core/domain/keystore"
	"myst/internal/client/core/enclaverepo"
)

func TestRepository(t *testing.T) {
	repo := enclaverepo.Repository{}

	err := repo.Initialize("12345678")
	assert.NoError(t, err)

	err = repo.Authenticate("12345678")
	assert.NoError(t, err)

	e, err := repo.Enclave()
	assert.NoError(t, err)
	assert.Len(t, e.Keystores(), 0)

	k, err := repo.Create(keystore.WithName("test"))
	assert.NoError(t, err)

	e, err = repo.Enclave()
	assert.NoError(t, err)
	assert.Len(t, e.Keystores(), 1)

	_, err = e.Keystore(k.Id())
	assert.NoError(t, err)
}
