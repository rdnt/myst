package enclaverepo_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"myst/pkg/rand"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/enclaverepo"
)

func TestRepository(t *testing.T) {
	repo, err := enclaverepo.New(t.TempDir())
	assert.NilError(t, err)

	pass, err := rand.String(16)
	assert.NilError(t, err)

	err = repo.Initialize(pass)
	assert.NilError(t, err)

	err = repo.Authenticate(pass)
	assert.NilError(t, err)

	k, err := repo.CreateKeystore(keystore.New(keystore.WithName("test")))
	assert.NilError(t, err)

	k2, err := repo.Keystore(k.Id)
	assert.NilError(t, err)

	assert.Equal(t, k.Id, k2.Id)
	assert.Equal(t, k.Name, k2.Name)
}
