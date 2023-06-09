package inmemrepo

import (
	"myst/src/server/application"
	"myst/src/server/application/domain/user"
)

func (r *Repository) CreateUser(u user.User) (user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.users[u.Id] = u

	return u, nil
}

func (r *Repository) User(id string) (user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	u, ok := r.users[id]
	if !ok {
		return user.User{}, application.ErrKeystoreNotFound
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

	return user.User{}, application.ErrKeystoreNotFound
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
		return user.User{}, application.ErrUserNotFound
	}

	r.users[u.Id] = u

	return u, nil
}
