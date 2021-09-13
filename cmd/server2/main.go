package main

import (
	"myst/app/server"
	invitationrepo "myst/app/server/invitationrepo/memory"
	"myst/app/server/invitationservice"
	keystorerepo "myst/app/server/keystorerepo/memory"
	"myst/app/server/keystoreservice"
	userrepo "myst/app/server/userrepo/memory"
	"myst/app/server/userservice"
)

func main() {
	keystoreRepo := keystorerepo.New()
	userRepo := userrepo.New()
	invitationRepo := invitationrepo.New()

	userService, err := userservice.New(
		userservice.WithUserRepository(userRepo),
		userservice.WithKeystoreRepository(keystoreRepo),
	)
	if err != nil {
		panic(err)
	}

	keystoreService, err := keystoreservice.New(
		keystoreservice.WithUserRepository(userRepo),
		keystoreservice.WithKeystoreRepository(keystoreRepo),
	)
	if err != nil {
		panic(err)
	}

	invitationService, err := invitationservice.New(
		invitationservice.WithUserRepository(userRepo),
		invitationservice.WithKeystoreRepository(keystoreRepo),
		invitationservice.WithInvitationRepository(invitationRepo),
	)
	if err != nil {
		panic(err)
	}

	app, err := server.New(
		server.WithKeystoreRepository(keystoreRepo),
		server.WithUserRepository(userRepo),
		server.WithInvitationRepository(invitationRepo),
		server.WithUserService(userService),
		server.WithKeystoreService(keystoreService),
		server.WithInvitationService(invitationService),
	)
	if err != nil {
		panic(err)
	}

	app.Start()
}
