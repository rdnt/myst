package application

type Option func(app *application) error

func WithRemote(remote Remote) Option {
	return func(app *application) error {
		app.remote = remote
		return nil
	}
}

func WithEnclaveRepository(enclave EnclaveRepository) Option {
	return func(app *application) error {
		app.enclave = enclave
		return nil
	}
}
