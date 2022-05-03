package remote

import (
	"myst/internal/server/api/http/generated"

	"github.com/pkg/errors"
)

// TODO: keystore sync

func (r *remote) Keystores() ([]*generated.Keystore, error) {
	ks, err := r.client.Keystores()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get keystores")
	}

	return ks, nil
}
