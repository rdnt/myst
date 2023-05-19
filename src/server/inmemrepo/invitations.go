package inmemrepo

import (
	"fmt"

	"myst/src/server/application/domain/invitation"
)

func (r *Repository) UserInvitations(userId string) ([]invitation.Invitation, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	invs := []invitation.Invitation{}

	for _, inv := range r.invitations {
		if inv.InviterId == userId {
			invs = append(invs, inv)
		}

		if inv.InviteeId == userId {
			invs = append(invs, inv)
		}
	}

	return invs, nil
}

func (r *Repository) UserInvitation(userId, invitationId string) (invitation.Invitation, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, inv := range r.invitations {
		if inv.InviterId == userId && inv.Id == invitationId {
			return inv, nil
		}

		if inv.InviteeId == userId && inv.Id == invitationId {
			return inv, nil
		}
	}

	return invitation.Invitation{}, invitation.ErrNotFound
}

func (r *Repository) CreateInvitation(inv invitation.Invitation) (invitation.Invitation, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, ok := r.invitations[inv.Id]
	if ok {
		return invitation.Invitation{}, fmt.Errorf("already exists")
	}

	r.invitations[inv.Id] = inv

	return inv, nil
}

func (r *Repository) Invitation(id string) (invitation.Invitation, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	i, ok := r.invitations[id]
	if !ok {
		return invitation.Invitation{}, invitation.ErrNotFound
	}

	return i, nil
}

func (r *Repository) Invitations() ([]invitation.Invitation, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	invitations := make([]invitation.Invitation, 0, len(r.invitations))
	for _, i := range r.invitations {
		invitations = append(invitations, i)
	}

	return invitations, nil
}

func (r *Repository) UpdateInvitation(inv invitation.Invitation) (invitation.Invitation, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, ok := r.invitations[inv.Id]
	if !ok {
		return invitation.Invitation{}, invitation.ErrNotFound
	}

	r.invitations[inv.Id] = inv

	return inv, nil
}

// func (r *Repository) Delete(id string) error {
// 	r.mux.Lock()
// 	defer r.mux.Unlock()
//
// 	delete(r.invitations, id)
// 	return nil
// }
