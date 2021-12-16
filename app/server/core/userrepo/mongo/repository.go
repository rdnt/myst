package userrepo

import (
	"myst/app/server/core/domain/user"
	"myst/pkg/mongo"
)

type Repository struct {
	db *mongo.Database
}

func (r *Repository) Create(opts ...user.Option) (*user.User, error) {
	panic("implement me")
}

func (r *Repository) User(id string) (*user.User, error) {
	panic("implement me")
}

func (r *Repository) Update(u *user.User) error {
	panic("implement me")
}

func (r *Repository) Users() ([]*user.User, error) {
	panic("implement me")
}

func (r *Repository) Delete(id string) error {
	panic("implement me")
}

func New(db *mongo.Database) user.Repository {
	return &Repository{
		db: db,
	}
}
