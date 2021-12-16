package main

import (
	"myst/app/server"
	invitationrepo "myst/app/server/core/invitationrepo/memory"
	"myst/app/server/core/invitationservice"
	keystorerepo "myst/app/server/core/keystorerepo/memory"
	"myst/app/server/core/keystoreservice"
	userrepo "myst/app/server/core/userrepo/memory"
	"myst/app/server/core/userservice"
	"myst/pkg/logger"

	"myst/app/server/rest"
)

var log = logger.New("app", logger.Red)

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

	api := rest.New(app)

	err = api.Run(":8080")
	if err != nil {
		panic(err)
	}
}
