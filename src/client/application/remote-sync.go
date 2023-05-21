package application

import (
	"github.com/pkg/errors"
)

func (app *application) Sync() error {
	log.Println("sync: started")
	defer log.Print("sync: finished")

	if !app.remote.Authenticated() {
		log.Print("sync: not signed in")
		return nil
	}

	rem, err := app.enclave.Credentials()
	if err != nil {
		return err
	}

	keystores, err := app.enclave.Keystores()
	if err != nil {
		return errors.WithMessage(err, "failed to get enclave")
	}

	remoteKeystores, err := app.remote.Keystores(rem.PrivateKey)
	if err != nil {
		return errors.WithMessage(err, "failed to query remote keystores")
	}

	for _, k := range keystores {
		if k.RemoteId == "" {
			// no need to sync keystore
			continue
		}

		rk, ok := remoteKeystores[k.RemoteId]
		if !ok {
			// apparently this keystore is no longer at the remote?
			// TODO: should we sync it again?
			continue
		}

		if rk.Version > k.Version {
			err = app.enclave.UpdateKeystore(rk)
			if err != nil {
				return errors.WithMessage(err, "failed to update local keystore")
			}

			log.Println("Local keystore updated", k.Id, k.RemoteId)
		} else if rk.Version < k.Version {
			_, err = app.remote.UpdateKeystore(k)
			if err != nil {
				return errors.WithMessage(err, "failed to update remote keystore")
			}

			log.Println("Credentials keystore updated", k.Id, k.RemoteId)
		} else {
			log.Println("No change", k.Id, k.RemoteId)
		}
	}

	return nil
}
