package inmemrepo

import (
	"fmt"

	"myst/src/server/application/domain/user"
)

func (r *Repository) CreateUser(u user.User) (user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, ok := r.users[u.Id]
	if ok {
		return user.User{}, fmt.Errorf("already exists")
	}

	for _, u2 := range r.users {
		if u2.Username == u.Username {
			return user.User{}, fmt.Errorf("already exists")
		}
	}

	r.users[u.Id] = u

	return u, nil
}

func (r *Repository) User(id string) (user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	u, ok := r.users[id]
	if !ok {
		return user.User{}, user.ErrNotFound
	}

	return u, nil
}

func (r *Repository) UserByUsername(username string) (user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, u := range r.users {
		if u.Username == username {
			return u, nil
		}
	}

	return user.User{}, user.ErrNotFound
}

func (r *Repository) Users() ([]user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	users := make([]user.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}

	return users, nil
}

func (r *Repository) UpdateUser(u user.User) (user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, ok := r.users[u.Id]
	if !ok {
		return user.User{}, fmt.Errorf("not found")
	}

	r.users[u.Id] = u

	return u, nil
}

// func (r *Repository) DeleteUser(id string) error {
// 	r.mux.Lock()
// 	defer r.mux.Unlock()
//
// 	delete(r.users, id)
// 	return nil
// }
