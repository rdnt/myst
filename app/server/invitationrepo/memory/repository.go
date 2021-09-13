package invitationrepo

import (
	"fmt"
	"sync"

	"myst/app/server/domain/invitation"
)

type Repository struct {
	mux         sync.Mutex
	invitations map[string]invitation.Invitation
}

func (r *Repository) Create(opts ...invitation.Option) (*invitation.Invitation, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	i, err := invitation.New(opts...)
	if err != nil {
		return nil, err
	}

	_, ok := r.invitations[i.Id()]
	if ok {
		return nil, fmt.Errorf("already exists")
	}

	r.invitations[i.Id()] = *i

	return i, nil
}

func (r *Repository) Invitation(id string) (*invitation.Invitation, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	k, ok := r.invitations[id]
	if !ok {
		return nil, invitation.ErrNotFound
	}

	return &k, nil
}

func (r *Repository) Invitations() ([]*invitation.Invitation, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	invitations := make([]*invitation.Invitation, 0, len(r.invitations))
	for _, k := range r.invitations {
		invitations = append(invitations, &k)
	}

	return invitations, nil
}

func (r *Repository) Update(s *invitation.Invitation) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, ok := r.invitations[s.Id()]
	if !ok {
		return fmt.Errorf("not found")
	}

	r.invitations[s.Id()] = *s
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
