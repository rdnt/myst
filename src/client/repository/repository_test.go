package repository_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"myst/src/client/application/domain/keystore"
	"myst/src/client/repository"
)

func TestRepository(t *testing.T) {
	repo, err := repository.New(t.TempDir())
	assert.NilError(t, err)

	err = repo.CreateEnclave("12345678")
	assert.NilError(t, err)

	err = repo.Authenticate("12345678")
	assert.NilError(t, err)

	k, err := repo.CreateKeystore(keystore.New(keystore.WithName("test")))
	assert.NilError(t, err)

	_, err = repo.Keystore(k.Id)
	assert.NilError(t, err)
}
