package main

import (
	application "myst/internal/server"
	"myst/internal/server/api/http"
	invitationrepo "myst/internal/server/core/invitationrepo/memory"
	"myst/internal/server/core/invitationservice"
	keystorerepo "myst/internal/server/core/keystorerepo/memory"
	"myst/internal/server/core/keystoreservice"
	userrepo "myst/internal/server/core/userrepo/memory"
	"myst/internal/server/core/userservice"
	"myst/pkg/logger"
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

	app, err := application.New(
		application.WithKeystoreRepository(keystoreRepo),
		application.WithUserRepository(userRepo),
		application.WithInvitationRepository(invitationRepo),
		application.WithUserService(userService),
		application.WithKeystoreService(keystoreService),
		application.WithInvitationService(invitationService),
	)
	if err != nil {
		panic(err)
	}

	api := http.New(app)

	err = api.Run(":8081")
	if err != nil {
		panic(err)
	}
}
