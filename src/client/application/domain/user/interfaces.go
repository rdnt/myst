package user

type Service interface {
	SignIn(username, password string) (User, error)
	Register(username, password string) (User, error)
	CurrentUser() (*User, error)
}
