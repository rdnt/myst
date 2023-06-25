package enclaverepo_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"myst/pkg/rand"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/enclaverepo"
)

func TestRepository(t *testing.T) {
	t.Parallel()

	repo := enclaverepo.New(t.TempDir())

	pass, err := rand.String(16)
	assert.NilError(t, err)

	key, err := repo.Initialize(pass)
	assert.NilError(t, err)

	key2, err := repo.Authenticate(pass)
	assert.NilError(t, err)

	assert.DeepEqual(t, key, key2)

	k, err := repo.CreateKeystore(key, keystore.New(keystore.WithName("test")))
	assert.NilError(t, err)

	k2, err := repo.Keystore(key, k.Id)
	assert.NilError(t, err)

	assert.Equal(t, k.Id, k2.Id)
	assert.Equal(t, k.Name, k2.Name)
}
