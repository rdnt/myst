package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"myst/src/client/application/domain/keystore"
	"myst/src/client/repository"
)

func TestRepository(t *testing.T) {
	repo, err := repository.New(t.TempDir())
	assert.NoError(t, err)

	err = repo.CreateEnclave("12345678")
	assert.NoError(t, err)

	err = repo.Authenticate("12345678")
	assert.NoError(t, err)

	k, err := repo.CreateKeystore(keystore.New(keystore.WithName("test")))
	assert.NoError(t, err)

	_, err = repo.Keystore(k.Id)
	assert.NoError(t, err)
}
