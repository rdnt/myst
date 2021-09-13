package userrepo

import (
	"go.mongodb.org/mongo-driver/mongo"

	"myst/app/server/domain/user"
)

type Repository struct {
	db *mongo.Database
}

func (r *Repository) CreateUser(opts ...user.Option) (*user.User, error) {
	panic("implement me")
}

func (r *Repository) User(id string) (*user.User, error) {
	panic("implement me")
}

func (r *Repository) UpdateUser(u *user.User) error {
	panic("implement me")
}

func (r *Repository) Users() ([]*user.User, error) {
	panic("implement me")
}

func (r *Repository) DeleteUser(id string) error {
	panic("implement me")
}

func New(db *mongo.Database) user.Repository {
	return &Repository{
		db: db,
	}
}
