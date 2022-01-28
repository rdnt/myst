package keystorerepo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/keystorerepo"
)

func TestRepository(t *testing.T) {
	repo, err := keystorerepo.New(t.TempDir())
	assert.NoError(t, err)

	err = repo.Initialize("12345678")
	assert.NoError(t, err)

	err = repo.Authenticate("12345678")
	assert.NoError(t, err)

	k, err := repo.Create(keystore.WithName("test"))
	assert.NoError(t, err)

	_, err = repo.Keystore(k.Id())
	assert.NoError(t, err)
}
