package application

import (
	"time"

	"github.com/pkg/errors"

	"myst/pkg/rand"
)

func (app *application) newSession() (sessionId []byte, err error) {
	// check is active are by the same user of the enclave.
	sessionKey, err := rand.Bytes(16)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate session key")
	}

	app.sessions[string(sessionKey)] = time.Now()

	return sessionKey, nil
}

func (app *application) sessionActive(sessionId []byte) bool {
	if sessionId == nil {
		return false
	}
	_, ok := app.sessions[string(sessionId)]
	return ok
}

func (app *application) startHealthCheck() {
	ticker := time.NewTicker(20 * time.Second) // FIXME @rdnt: make this configurable
	defer ticker.Stop()

	for range ticker.C {
		app.mux.Lock()

		for sessionId, lastHealthCheck := range app.sessions {
			elapsed := time.Since(lastHealthCheck)

			if elapsed < time.Minute {
				continue
			}

			delete(app.sessions, sessionId)

			if len(app.sessions) == 0 && app.key != nil {
				app.key = nil
			}
		}

		app.mux.Unlock()
	}
}
