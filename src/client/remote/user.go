package remote

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"myst/src/client/application/domain/user"
	"myst/src/server/rest/generated"
)

func (r *remote) Authenticate(username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := r.client.LoginWithResponse(
		ctx, generated.LoginJSONRequestBody{
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		return errors.Wrap(err, "login failed")
	}

	if res.StatusCode() != http.StatusOK || res.JSON200 == nil {
		return errors.New("invalid response")
	}

	if res.JSON200.Token == "" {
		return errors.New("invalid token")
	}

	u := user.User{
		Id:        res.JSON200.User.Id,
		Username:  res.JSON200.User.Username,
		PublicKey: res.JSON200.User.PublicKey,
	}

	r.user = &u
	r.bearerToken = res.JSON200.Token

	return nil
}

func (r *remote) Authenticated() bool {
	return r.bearerToken != ""
}

func (r *remote) CurrentUser() *user.User {
	return r.user
}

func (r *remote) Register(username, password string, publicKey []byte) (user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := r.client.RegisterWithResponse(
		ctx, generated.RegisterJSONRequestBody{
			Username:  username,
			Password:  password,
			PublicKey: publicKey,
		},
	)
	if err != nil {
		return user.User{}, errors.Wrap(err, "failed to register")
	}

	if res.StatusCode() != http.StatusCreated || res.JSON201 == nil {
		return user.User{}, errors.New("invalid response")
	}

	if res.JSON201.Token == "" {
		return user.User{}, errors.New("invalid token")
	}

	u := user.User{
		Id:        res.JSON201.User.Id,
		Username:  res.JSON201.User.Username,
		PublicKey: publicKey,
	}

	r.user = &u
	r.bearerToken = res.JSON201.Token

	return u, nil
}
