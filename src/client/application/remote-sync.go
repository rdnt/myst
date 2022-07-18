package application

import (
	"fmt"

	"github.com/pkg/errors"
)

func (app *application) sync() error {
	log.Println("Sync started.")
	defer log.Print("Sync finished.")

	if !app.remote.SignedIn() {
		return errors.New("signed out")
	}

	invs, err := app.remote.Invitations()
	if err != nil {
		return errors.WithMessage(err, "failed to get invitations from remote")
	}

	for _, inv := range invs {
		if inv.Accepted() && inv.Inviter.Id == app.remote.CurrentUser().Id {
			log.Print("Finalizing invitation ", inv.Id, "...")

			_, err = app.FinalizeInvitation(inv.Id)
			if err != nil {
				log.Error(err)
				continue
			}

			log.Print("Invitation ", inv.Id, " finalized.")
		}
	}

	keystores, err := app.keystores.Keystores()
	if err != nil {
		return errors.WithMessage(err, "failed to get keystores")
	}

	rem, err := app.credentials.Remote()
	if err != nil {
		return errors.WithMessage(err, "failed to get remote config")
	}

	remoteKeystores, err := app.remote.Keystores(rem.PrivateKey)
	if err != nil {
		return errors.WithMessage(err, "failed to get remote keystores")
	}

	for _, k := range keystores {
		if k.RemoteId == "" {
			k, err = app.remote.CreateKeystore(k)
			if err != nil {
				fmt.Println("failed to create keystore")
				return err
			}

			err = app.keystores.UpdateKeystore(k)
			if err != nil {
				return err
			}

			fmt.Println("keystore uploaded")

			continue
		}

		fmt.Println("syncing keystore", k)
		rk, err := app.remote.Keystore(k.RemoteId, k.Key)
		if err != nil {
			fmt.Println("remote invalid response")
			return err
		}

		if rk.Version > k.Version {
			err = app.keystores.UpdateKeystore(rk)
			if err != nil {
				fmt.Println("local update failed")
				return err
			}
			fmt.Println("local updated", k.Id, k.RemoteId)
		} else if rk.Version < k.Version {
			_, err = app.remote.UpdateKeystore(k)
			if err != nil {
				fmt.Println("remote update failed")
				return err
			}
			fmt.Println("remote updated", k.Id, k.RemoteId)
		} else {
			fmt.Println("no change", k.Id, k.RemoteId)
		}

		fmt.Println(k.Version)

		// if _, ok := remoteKeystores[k.Id]; !ok {
		//	// sync from local to remote
		//
		//
		// }
	}

	for _, rk := range remoteKeystores {
		fmt.Println("REMOTE KEYSTORE", rk)
		if k, ok := keystores[rk.Id]; !ok {
			// sync from remote to local

			_, err = app.keystores.CreateKeystore(rk)
			if err != nil {
				fmt.Println("local update failed")
				return err
			}
			fmt.Println("local updated", rk.Id, rk.RemoteId)
		} else {
			if rk.Version > k.Version {
				err = app.keystores.UpdateKeystore(rk)
				if err != nil {
					fmt.Println("local update failed")
					return err
				}
				fmt.Println("local updated", k.Id, k.RemoteId)
			} else if rk.Version < k.Version {
				_, err = app.remote.UpdateKeystore(k)
				if err != nil {
					fmt.Println("remote update failed")
					return err
				}
				fmt.Println("remote updated", k.Id, k.RemoteId)
			} else {
				fmt.Println("no change", k.Id, k.RemoteId)
			}
		}
	}

	// for _, k := range ks {
	//	if k.RemoteId == "" {
	//		k, err = app.remote.CreateKeystore(k)
	//		if err != nil {
	//			fmt.Println("failed to create keystore")
	//			return err
	//		}
	//
	//		err = app.keystores.UpdateKeystore(k)
	//		if err != nil {
	//			return err
	//		}
	//
	//		fmt.Println("keystore uploaded")
	//
	//		continue
	//	}
	//
	//	fmt.Println("syncing keystore", k)
	//	rk, err := app.remote.Keystore(k.RemoteId, k.Key)
	//	if err != nil {
	//		fmt.Println("remote invalid response")
	//		return err
	//	}
	//
	//	if rk.Version > k.Version {
	//		err = app.keystores.UpdateKeystore(rk)
	//		if err != nil {
	//			fmt.Println("local update failed")
	//			return err
	//		}
	//		fmt.Println("local updated", k.Id, k.RemoteId)
	//	} else if rk.Version < k.Version {
	//		err = app.remote.UpdateKeystore(k)
	//		if err != nil {
	//			fmt.Println("remote update failed")
	//			return err
	//		}
	//		fmt.Println("remote updated", k.Id, k.RemoteId)
	//	} else {
	//		fmt.Println("no change", k.Id, k.RemoteId)
	//	}
	//
	//	fmt.Println(k.Version)
	//
	// }

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
