package remote_test

import (
	"fmt"
	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/keystorerepo"
	"myst/internal/client/remote"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemote(t *testing.T) {
	dir, err := os.MkdirTemp("", "test-remote-*")
	assert.NoError(t, err)

	defer func() {
		err := os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	repo, err := keystorerepo.New(dir)
	assert.NoError(t, err)

	r1, err := remote.New()
	assert.NoError(t, err)

	r2, err := remote.New()
	assert.NoError(t, err)

	user1, pass1 := "rdnt", "1234"
	err = r1.SignIn(user1, pass1)
	assert.NoError(t, err)

	user2, pass2 := "abcd", "5678"
	err = r2.SignIn(user2, pass2)
	assert.NoError(t, err)

	err = repo.Initialize("12345678")
	assert.NoError(t, err)

	k, err := repo.Create(keystore.WithName("test-keystore"))
	assert.NoError(t, err)

	keystoreKey, err := repo.KeystoreKey(k.Id())
	assert.NoError(t, err)

	genk, err := r1.CreateKeystore(
		k.Name(), keystoreKey, k,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	inv, err := r1.CreateInvitation(genk.Id, user2)
	assert.NoError(t, err)

	inv, err = r2.AcceptInvitation(genk.Id, inv.Id)
	assert.NoError(t, err)

	inv, err = r1.FinalizeInvitation(genk.Id, inv.Id, keystoreKey)
	assert.NoError(t, err)

	ks, err := r2.Keystores()
	assert.NoError(t, err)
	assert.Len(t, ks, 1)
}