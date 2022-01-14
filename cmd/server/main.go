package main

import (
	application "myst/internal/server"
	"myst/internal/server/api/http"
	invitationrepo "myst/internal/server/core/invitationrepo/memory"
	keystorerepo "myst/internal/server/core/keystorerepo/memory"
	userrepo "myst/internal/server/core/userrepo/memory"
	"myst/pkg/logger"
)

var log = logger.New("app", logger.Red)

func main() {
	keystoreRepo := keystorerepo.New()
	userRepo := userrepo.New()
	invitationRepo := invitationrepo.New()

	app, err := application.New(
		application.WithKeystoreRepository(keystoreRepo),
		application.WithUserRepository(userRepo),
		application.WithInvitationRepository(invitationRepo),
	)
	if err != nil {
		panic(err)
	}

	api := http.New(app)

	err = api.Run(":8080")
	if err != nil {
		panic(err)
	}
}
