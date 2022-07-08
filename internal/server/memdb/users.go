package memdb

import (
	"fmt"

	"myst/internal/server/application/domain/user"
)

type User struct {
	user.User
	passwordHash string
}

func (r *Repository) CreateUser(opts ...user.Option) (user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	u, err := user.New(opts...)
	if err != nil {
		return user.User{}, err
	}

	_, ok := r.users[u.Id]
	if ok {
		return user.User{}, fmt.Errorf("already exists")
	}

	for _, u2 := range r.users {
		if u2.User.Username == u.Username {
			return user.User{}, fmt.Errorf("already exists")
		}
	}

	ru := User{
		User:         u,
		passwordHash: "asd", // todo: generate hash
	}

	r.users[ru.Id] = ru

	return ru.User, nil
}

func (r *Repository) User(id string) (user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	u, ok := r.users[id]
	if !ok {
		return user.User{}, user.ErrNotFound
	}

	return u.User, nil
}

func (r *Repository) UserByUsername(username string) (user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, u := range r.users {
		if u.Username == username {
			return u.User, nil
		}
	}

	return user.User{}, user.ErrNotFound
}

func (r *Repository) Users() ([]user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	users := make([]user.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u.User)
	}

	return users, nil
}

func (r *Repository) UpdateUser(u *user.User) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	ru, ok := r.users[u.Id]
	if !ok {
		return fmt.Errorf("not found")
	}

	ru2 := User{
		User:         *u,
		passwordHash: ru.passwordHash,
	}

	r.users[ru2.Id] = ru2

	return nil
}

func (r *Repository) DeleteUser(id string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	delete(r.users, id)
	return nil
}
