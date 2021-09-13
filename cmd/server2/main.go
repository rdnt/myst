package main

import (
	"myst/app/server"
	keystorerepo "myst/app/server/keystorerepo/memory"
	"myst/app/server/keystoreservice"
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

	keystoreService, err := keystoreservice.New(
		keystoreservice.WithUserRepository(userRepo),
		keystoreservice.WithKeystoreRepository(keystoreRepo),
	)
	if err != nil {
		panic(err)
	}

	app, err := server.New(
		server.WithKeystoreRepository(keystoreRepo),
		server.WithUserRepository(userRepo),
		server.WithUserService(userService),
		server.WithKeystoreService(keystoreService),
	)
	if err != nil {
		panic(err)
	}

	app.Start()
}
