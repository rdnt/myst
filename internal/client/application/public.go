package application

func (app *application) SignIn(username, password string) error {
	return app.repositories.remote.SignIn(username, password)
}

func (app *application) SignOut() error {
	return app.repositories.remote.SignOut()
}
