package remote

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"myst/src/client/application/domain/invitation"
	"myst/src/server/rest/generated"
)

var (
	ErrNotAuthenticated = errors.New("not authenticated")
)

func (r *remote) Invitation(id string) (invitation.Invitation, error) {
	if !r.Authenticated() {
		return invitation.Invitation{}, ErrNotAuthenticated
	}

	res, err := r.client.GetInvitationWithResponse(context.Background(), id)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get invitation")
	}

	if res.JSON200 == nil {
		return invitation.Invitation{}, fmt.Errorf("invalid response")
	}

	return InvitationFromJSON(*res.JSON200)
}

func (r *remote) CreateInvitation(keystoreRemoteId, inviteeUsername string) (invitation.Invitation, error) {
	if !r.Authenticated() {
		return invitation.Invitation{}, ErrNotAuthenticated
	}

	res, err := r.client.CreateInvitationWithResponse(
		context.Background(), keystoreRemoteId, generated.CreateInvitationJSONRequestBody{
			Invitee: inviteeUsername,
		},
	)
	if err != nil {
		return invitation.Invitation{}, err
	}

	if res.JSON201 == nil {
		return invitation.Invitation{}, fmt.Errorf("invalid response: %s", string(res.Body))
	}

	inv, err := InvitationFromJSON(*res.JSON201)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}

func (r *remote) Address() string {
	return r.address
}

func (r *remote) Invitations() (map[string]invitation.Invitation, error) {
	if !r.Authenticated() {
		return nil, ErrNotAuthenticated
	}

	res, err := r.client.InvitationsWithResponse(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get invitations")
	}

	if res.JSON200 == nil {
		return nil, fmt.Errorf("invalid response")
	}

	invs := map[string]invitation.Invitation{}
	for _, inv := range *res.JSON200 {
		inv, err := InvitationFromJSON(inv)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to parse invitation")
		}

		invs[inv.Id] = inv
	}

	return invs, nil
}

func (r *remote) AcceptInvitation(invitationId string) (invitation.Invitation, error) {
	res, err := r.client.AcceptInvitationWithResponse(
		context.Background(), invitationId,
	)
	if err != nil {
		return invitation.Invitation{}, err
	}

	if res.JSON200 == nil {
		return invitation.Invitation{}, fmt.Errorf("invalid response")
	}

	inv, err := InvitationFromJSON(*res.JSON200)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}

func (r *remote) DeleteInvitation(id string) (invitation.Invitation, error) {
	if !r.Authenticated() {
		return invitation.Invitation{}, ErrNotAuthenticated
	}

	res, err := r.client.DeleteInvitationWithResponse(
		context.Background(), id,
	)
	if err != nil {
		return invitation.Invitation{}, err
	}

	if res.JSON200 == nil {
		return invitation.Invitation{}, fmt.Errorf("invalid response")
	}

	inv, err := InvitationFromJSON(*res.JSON200)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}

func (r *remote) FinalizeInvitation(invitationId string, encryptedKeystoreKey []byte) (invitation.Invitation, error) {
	res, err := r.client.FinalizeInvitationWithResponse(context.Background(), invitationId, generated.FinalizeInvitationJSONRequestBody{
		KeystoreKey: encryptedKeystoreKey,
	})
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to finalize invitation")
	}

	if res.JSON200 == nil {
		return invitation.Invitation{}, fmt.Errorf("invalid response")
	}

	return InvitationFromJSON(*res.JSON200)
}
