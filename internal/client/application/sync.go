package application

import (
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
		if inv.Accepted() && inv.InviterId == app.remote.CurrentUser().Id {
			log.Print("Finalizing invitation ", inv.Id, "...")

			_, err = app.FinalizeInvitation(inv.Id)
			if err != nil {
				log.Error(err)
				continue
			}

			log.Print("Invitation ", inv.Id, " finalized.")
		}
	}

	ks, err := app.remote.Keystores()
	if err != nil {
		return err
	}

	log.Debug(ks)

	return nil
}
