package main

import (
	"myst/app/server"
	keystorerepo "myst/app/server/keystorerepo/memory"
	userrepo "myst/app/server/userrepo/memory"
	"myst/app/server/userservice"
)

func main() {
	keystoreRepo := keystorerepo.New()
	userRepo := userrepo.New()

	userService, err := userservice.New(
		userservice.WithUserRepository(userRepo),
		userservice.WithKeystoreRepository(keystoreRepo),
	)
	if err != nil {
		panic(err)
	}

	app, err := server.New(
		server.WithKeystoreRepository(keystoreRepo),
		server.WithUserRepository(userRepo),
		server.WithUserService(userService),
	)
	if err != nil {
		panic(err)
	}

	app.Start()
}
