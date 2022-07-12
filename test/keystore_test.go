package test

import (
	"context"
	"time"

	"myst/pkg/optional"
	"myst/src/client/rest/generated"
)

func (s *IntegrationTestSuite) TestKeystores() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.Run("Keystores are empty", func() {
		res, err := s.client1.KeystoresWithResponse(ctx)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON200)
		s.Require().Len(*res.JSON200, 0)
	})

	keystoreName := s.rand()
	var keystoreId string

	s.Run("Keystore created", func() {
		res, err := s.client1.CreateKeystoreWithResponse(ctx,
			generated.CreateKeystoreJSONRequestBody{Name: keystoreName},
		)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON201)
		s.Require().Equal(keystoreName, (*res.JSON201).Name)

		keystoreId = (*res.JSON201).Id
	})

	s.Run("Keystore can be queried", func() {
		res, err := s.client1.KeystoreWithResponse(ctx, keystoreId)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON200)
		s.Require().Equal(keystoreId, (*res.JSON200).Id)
	})

	s.Run("Keystores contain created keystore", func() {
		res, err := s.client1.KeystoresWithResponse(ctx)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON200)
		s.Require().Len(*res.JSON200, 1)
		s.Require().Equal((*res.JSON200)[0].Id, keystoreId)
	})

	website := s.rand()
	username := s.rand()
	password := s.rand()
	notes := s.rand()
	var entryId string

	s.Run("Entry created", func() {
		res, err := s.client1.CreateEntryWithResponse(ctx, keystoreId,
			generated.CreateEntryJSONRequestBody{
				Website:  website,
				Username: username,
				Password: password,
				Notes:    notes,
			},
		)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON201)

		s.Require().Equal((*res.JSON201).Website, website)
		s.Require().Equal((*res.JSON201).Username, username)
		s.Require().Equal((*res.JSON201).Password, password)
		s.Require().Equal((*res.JSON201).Notes, notes)

		entryId = (*res.JSON201).Id
	})

	s.Run("Entry can be queried", func() {
		res, err := s.client1.KeystoreWithResponse(ctx, keystoreId)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON200)
		s.Require().Len((*res.JSON200).Entries, 1)
		s.Require().Equal((*res.JSON200).Entries[0].Id, entryId)
	})

	newPassword := s.rand()
	newNotes := s.rand()

	s.Run("Entry can be updated", func() {
		res, err := s.client1.UpdateEntryWithResponse(ctx, keystoreId, entryId,
			generated.UpdateEntryJSONRequestBody{
				Password: optional.Ref(newPassword),
				Notes:    optional.Ref(newNotes),
			},
		)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON200)
		s.Require().Equal((*res.JSON200).Password, newPassword)
		s.Require().Equal((*res.JSON200).Notes, newNotes)

		res2, err := s.client1.KeystoreWithResponse(ctx, keystoreId)
		s.Require().Nil(err)
		s.Require().NotNil(res2.JSON200)
		s.Require().Len((*res2.JSON200).Entries, 1)
		s.Require().Equal((*res2.JSON200).Entries[0].Password, newPassword)
		s.Require().Equal((*res2.JSON200).Entries[0].Notes, newNotes)
	})

	s.Run("Entry can be deleted", func() {
		_, err := s.client1.DeleteEntryWithResponse(ctx, keystoreId, entryId)
		s.Require().Nil(err)

		res, err := s.client1.KeystoreWithResponse(ctx, keystoreId)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON200)
		s.Require().Len((*res.JSON200).Entries, 0)
	})

	s.Run("Keystore can be deleted", func() {
		_, err := s.client1.DeleteKeystore(ctx, keystoreId)
		s.Require().Nil(err)

		res, err := s.client1.KeystoresWithResponse(ctx)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON200)
		s.Require().Len(*res.JSON200, 0)
	})
}
