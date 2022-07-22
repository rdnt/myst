package application

import (
	"github.com/pkg/errors"
)

func (app *application) Sync() error {
	log.Println("Sync started.")
	defer log.Print("Sync finished.")

	if !app.remote.SignedIn() {
		return errors.New("signed out")
	}

	keystores, err := app.keystores.Keystores()
	if err != nil {
		return errors.WithMessage(err, "failed to get keystores")
	}

	for _, k := range keystores {
		if k.RemoteId == "" {
			k, err = app.remote.CreateKeystore(k)
			if err != nil {
				return err
			}

			err = app.keystores.UpdateKeystore(k)
			if err != nil {
				return err
			}

			log.Println("Keystore uploaded", k.Id, k.RemoteId)

			continue
		}

		rk, err := app.remote.Keystore(k.RemoteId, k.Key)
		if err != nil {
			return err
		}

		if rk.Version > k.Version {
			err = app.keystores.UpdateKeystore(rk)
			if err != nil {
				return err
			}

			log.Println("Local keystore updated", k.Id, k.RemoteId)
		} else if rk.Version < k.Version {
			_, err = app.remote.UpdateKeystore(k)
			if err != nil {
				return err
			}

			log.Println("Remote keystore updated", k.Id, k.RemoteId)
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
//	ks, err := app.keystores.Keystores()
//	if err != nil {
//		return err
//	}
//
// }
