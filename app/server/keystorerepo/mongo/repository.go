package keystorerepo

import (
	"go.mongodb.org/mongo-driver/mongo"

	"myst/app/server/domain/keystore"
)

type Repository struct {
	db *mongo.Database
}

func (r *Repository) Create(k *keystore.Keystore) error {
	panic("implement me")
}

func (r *Repository) Keystore(id string) (*keystore.Keystore, error) {
	panic("implement me")
}

func (r *Repository) Update(k *keystore.Keystore) error {
	panic("implement me")
}

func (r *Repository) Keystores() ([]*keystore.Keystore, error) {
	panic("implement me")
}

func (r *Repository) Delete(id string) error {
	panic("implement me")
}

func New(db *mongo.Database) keystore.Repository {
	return &Repository{
		db: db,
	}
}
