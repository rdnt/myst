package invitationrepo

import (
	"fmt"
	"sync"

	"myst/internal/server/core/domain/invitation"
)

type Repository struct {
	mux         sync.Mutex
	invitations map[string]invitation.Invitation
}

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

func (r *Repository) CreateInvitation(opts ...invitation.Option) (invitation.Invitation, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	i, err := invitation.New(opts...)
	if err != nil {
		return invitation.Invitation{}, err
	}

	_, ok := r.invitations[i.Id]
	if ok {
		return invitation.Invitation{}, fmt.Errorf("already exists")
	}

	r.invitations[i.Id] = i

	return i, nil
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

func (r *Repository) UpdateInvitation(i *invitation.Invitation) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, ok := r.invitations[i.Id]
	if !ok {
		return invitation.ErrNotFound
	}

	r.invitations[i.Id] = *i

	return nil
}

func (r *Repository) Delete(id string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	delete(r.invitations, id)
	return nil
}

func New() invitation.Repository {
	return &Repository{
		invitations: make(map[string]invitation.Invitation),
	}
}
