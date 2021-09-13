package keystorerepo

import (
	"fmt"
	"sync"

	"myst/app/server/domain/keystore"
	"myst/app/server/domain/user"
)

type Repository struct {
	mux   sync.Mutex
	users map[string]user.User
}

func (r *Repository) CreateUser(opts ...user.Option) (*user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	u, err := user.New(opts...)
	if err != nil {
		return nil, err
	}

	_, ok := r.users[u.Id()]
	if ok {
		return nil, fmt.Errorf("already exists")
	}

	r.users[u.Id()] = *u
	return u, nil
}

func (r *Repository) User(id string) (*user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	s, ok := r.users[id]
	if !ok {
		return nil, keystore.ErrNotFound
	}

	return &s, nil
}

func (r *Repository) Users() ([]*user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	users := make([]*user.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, &u)
	}

	return users, nil
}

func (r *Repository) UpdateUser(u *user.User) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, ok := r.users[u.Id()]
	if !ok {
		return fmt.Errorf("not found")
	}

	r.users[u.Id()] = *u
	return nil
}

func (r *Repository) DeleteUser(id string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	delete(r.users, id)
	return nil
}

func New() user.Repository {
	return &Repository{
		users: make(map[string]user.User),
	}
}
