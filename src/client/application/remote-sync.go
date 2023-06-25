package application

import (
	"github.com/pkg/errors"
)

func (app *application) Sync() error {
	app.mux.Lock()
	defer app.mux.Unlock()

	if app.key == nil || !app.remote.Authenticated() {
		log.Print("sync: not signed in")
		return nil
	}

	log.Println("sync: started")
	defer log.Print("sync: finished")

	rem, err := app.enclave.Credentials(app.key)
	if err != nil {
		return errors.WithMessage(err, "failed to get credentials")
	}

	keystores, err := app.enclave.Keystores(app.key)
	if err != nil {
		return errors.WithMessage(err, "failed to get local keystores")
	}

	remoteKeystores, err := app.remote.Keystores(rem.PrivateKey)
	if err != nil {
		return errors.WithMessage(err, "failed to get remote keystores")
	}

	for _, k := range keystores {
		if k.RemoteId == "" {
			continue
		}

		rk, ok := remoteKeystores[k.Id]
		if !ok {
			log.Println("sync: skipped (not upstream)", k.Id, k.RemoteId)

			continue
		}

		if rk.Version > k.Version {
			_, err = app.enclave.UpdateKeystore(app.key, rk)
			if err != nil {
				return errors.WithMessage(err, "failed to update local keystore")
			}
			log.Println("sync: local keystore updated", k.Id, k.RemoteId)
		} else if rk.Version < k.Version {
			_, err = app.remote.UpdateKeystore(k)
			if err != nil {
				return errors.WithMessage(err, "failed to update remote keystore")
			}

			log.Println("sync: credentials keystore updated", k.Id, k.RemoteId)
		} else {
			log.Println("sync: skipped (no change)", k.Id, k.RemoteId)
		}
	}

	return nil
}
