package repository

import (
	"myst/src/server/application/domain/user"
)

func UserToJSON(u user.User) User {
	return User{
		Id:           u.Id,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		PublicKey:    u.PublicKey,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func UserFromJSON(u User) user.User {
	return user.User{
		Id:           u.Id,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		PublicKey:    u.PublicKey,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
