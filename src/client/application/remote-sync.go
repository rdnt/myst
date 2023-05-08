package application

import (
	"github.com/pkg/errors"
)

func (app *application) Sync() error {
	log.Println("sync: started")
	defer log.Print("sync: finished")

	if !app.remote.SignedIn() {
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
		return err
	}

	log.Debugln(remoteKeystores)

	for _, k := range keystores {
		if k.RemoteId == "" {
			// no need to sync keystore
			continue
		}

		rk, err := app.remote.Keystore(k.RemoteId, k.Key)
		if err != nil {
			return err
		}

		if rk.Version > k.Version {
			err = app.enclave.UpdateKeystore(rk)
			if err != nil {
				return err
			}

			log.Println("Local keystore updated", k.Id, k.RemoteId)
		} else if rk.Version < k.Version {
			_, err = app.remote.UpdateKeystore(k)
			if err != nil {
				return err
			}

			log.Println("Credentials keystore updated", k.Id, k.RemoteId)
		} else {
			log.Println("No change", k.Id, k.RemoteId)
		}
	}

	return nil
}

// func (app *applicationrefactor) syncKeystoreKeys() error {
//	invs, err := app.remote.Invitations()
//	if err != nil {
//		return err
//	}
//
//	ks, err := app.enclave.Keystores()
//	if err != nil {
//		return err
//	}
//
// }
