package userrepo

import (
	"fmt"
	"sync"

	"myst/app/server/domain/user"
)

type Repository struct {
	mux   sync.Mutex
	users map[string]User
}

type User struct {
	user.User
	passwordHash string
}

func (r *Repository) Create(opts ...user.Option) (*user.User, error) {
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

	ru := User{
		User:         *u,
		passwordHash: "asd", // todo: generate hash
	}

	r.users[ru.Id()] = ru

	return &ru.User, nil
}

func (r *Repository) User(id string) (*user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	u, ok := r.users[id]
	if !ok {
		return nil, user.ErrNotFound
	}

	return &u.User, nil
}

func (r *Repository) Users() ([]*user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	users := make([]*user.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, &u.User)
	}

	return users, nil
}

func (r *Repository) Update(u *user.User) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	ru, ok := r.users[u.Id()]
	if !ok {
		return fmt.Errorf("not found")
	}

	ru2 := User{
		User:         *u,
		passwordHash: ru.passwordHash,
	}

	r.users[ru2.Id()] = ru2

	return nil
}

func (r *Repository) Delete(id string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	delete(r.users, id)
	return nil
}

func New() user.Repository {
	return &Repository{
		users: make(map[string]User),
	}
}
