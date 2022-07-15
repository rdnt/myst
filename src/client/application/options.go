package application

type Option func(app *application) error

func WithRemote(remote Remote) Option {
	return func(app *application) error {
		app.remote = remote
		return nil
	}
}

func WithKeystoreRepository(repo KeystoreRepository) Option {
	return func(app *application) error {
		app.keystores = repo
		return nil
	}
}
