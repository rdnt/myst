package client

import (
	"myst/app/client/core/domain/keystore"
)

func (app *application) CreateKeystore(name string, passphrase []byte) (*keystore.Keystore, error) {
	return app.keystoreService.Create(
		keystore.WithName(name),
		keystore.WithPassphrase(passphrase),
	)
}
