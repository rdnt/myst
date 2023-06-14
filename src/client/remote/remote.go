package remote

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"myst/src/client/application"
	"myst/src/client/application/domain/user"
	"myst/src/server/rest/generated"
)

type remote struct {
	address     string
	client      *generated.ClientWithResponses
	bearerToken string
	user        *user.User
}

func New(opts ...Option) (application.Remote, error) {
	r := &remote{}

	for _, opt := range opts {
		if opt != nil {
			opt(r)
		}
	}

	var err error
	r.client, err = generated.NewClientWithResponses(
		r.address,
		generated.WithRequestEditorFn(r.authenticate()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http client")
	}

	return r, nil
}

type Option func(*remote)

func WithAddress(address string) Option {
	return func(r *remote) {
		r.address = address
	}
}

func (r *remote) Address() string {
	return r.address
}

func (r *remote) authenticate() generated.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+r.bearerToken)

		return nil
	}
}
