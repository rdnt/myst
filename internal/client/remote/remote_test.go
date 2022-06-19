package remote_test

//func TestRemote(t *testing.T) {
//	// FIXME: how to properly test this?
//	t.Skip("skip remote test")
//
//	dir, err := os.MkdirTemp("", "test-remote-*")
//	assert.NoError(t, err)
//
//	repo, err := keystorerepo.New(dir)
//	assert.NoError(t, err)
//
//	keystoreService, err := keystoreservice.New(keystoreservice.WithKeystoreRepository(repo))
//	assert.NoError(t, err)
//
//	defer func() {
//		err := os.RemoveAll(dir)
//		assert.NoError(t, err)
//	}()
//
//	ports, err := freeport.GetFreePorts(2)
//	assert.NoError(t, err)
//
//	r1, err := remote.New(keystoreService, fmt.Sprintf("http://localhost:%d", ports[0]))
//	assert.NoError(t, err)
//
//	r2, err := remote.New(keystoreService, fmt.Sprintf("http://localhost:%d", ports[1]))
//	assert.NoError(t, err)
//
//	user1, pass1 := "rdnt", "1234"
//	err = r1.SignIn(user1, pass1)
//	assert.NoError(t, err)
//
//	user2, pass2 := "abcd", "5678"
//	err = r2.SignIn(user2, pass2)
//	assert.NoError(t, err)
//
//	err = repo.Initialize("12345678")
//	assert.NoError(t, err)
//
//	k, err := repo.CreateKeystore(keystore.New(keystore.WithName("test-keystore")))
//	assert.NoError(t, err)
//
//	//keystoreKey, err := repo.KeystoreKey(k.Id())
//	//assert.NoError(t, err)
//
//	genk, err := r1.CreateKeystore(k)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	inv, err := r1.CreateInvitation(invitation.New(invitation.WithKeystoreId(genk.Id), invitation.WithInviteeId(user2)))
//	assert.NoError(t, err)
//
//	inv, err = r2.AcceptInvitation(inv.Id)
//	assert.NoError(t, err)
//
//	inv, err = r1.FinalizeInvitation(inv.Id)
//	assert.NoError(t, err)
//
//	ks, err := r2.Keystores()
//	assert.NoError(t, err)
//	assert.Len(t, ks, 1)
//}
