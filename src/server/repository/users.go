package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"myst/src/server/application/domain/user"
)

type User struct {
	Id           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"passwordHash"`
	PublicKey    []byte    `json:"publicKey"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

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

	r.users[u.Id] = UserToJSON(u)

	{ // debug
		b, _ := json.Marshal(r.users)
		_ = os.WriteFile("users.json", b, 0666)
	}

	return u, nil
}

func (r *Repository) User(id string) (user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	u, ok := r.users[id]
	if !ok {
		return user.User{}, user.ErrNotFound
	}

	return UserFromJSON(u), nil
}

func (r *Repository) UserByUsername(username string) (user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, u := range r.users {
		if u.Username == username {
			return UserFromJSON(u), nil
		}
	}

	return user.User{}, user.ErrNotFound
}

func (r *Repository) Users() ([]user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	users := make([]user.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, UserFromJSON(u))
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

	u2 := UserToJSON(u)
	r.users[u.Id] = u2

	{ // debug
		b, _ := json.Marshal(r.users)
		_ = os.WriteFile("users.json", b, 0666)
	}

	return UserFromJSON(u2), nil
}

// func (r *Repository) DeleteUser(id string) error {
// 	r.mux.Lock()
// 	defer r.mux.Unlock()
//
// 	delete(r.users, id)
// 	return nil
// }
