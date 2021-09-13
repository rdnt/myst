package user

type Repository interface {
	CreateUser(opts ...Option) (*User, error)
	User(id string) (*User, error)
	UpdateUser(*User) error
	Users() ([]*User, error)
	DeleteUser(id string) error
}

type Service interface {
	RegisterUser(u *User, password string) error
	CreateKeystore(u *User, keystore []byte) error
}
