package remote

import (
	"myst/internal/client/application/domain/invitation"
	"myst/internal/server/api/http/generated"
	"myst/pkg/crypto"
	"myst/pkg/timestamp"

	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"
)

func (r *remote) CreateInvitation(keystoreId, inviteeId string) (*generated.Invitation, error) {
	inv, err := r.client.CreateInvitation(keystoreId, inviteeId, r.publicKey)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create invitation")
	}

	return inv, nil
}

func (r *remote) Invitations() ([]invitation.Invitation, error) {
	restInvs, err := r.client.Invitations()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query invitations")
	}

	var invs []invitation.Invitation
	for _, inv := range restInvs {
		inv, err := InvitationFromJSON(inv)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to parse invitation")
		}

		invs = append(invs, inv)
	}

	return invs, nil
}

func (r *remote) AcceptInvitation(keystoreId, invitationId string) (*generated.Invitation, error) {
	inv, err := r.client.AcceptInvitation(keystoreId, invitationId, r.publicKey)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to accept invitation")
	}

	return inv, nil
}
func (r *remote) FinalizeInvitation(localKeystoreId, keystoreId, invitationId string) (*generated.Invitation, error) {
	inv, err := r.client.Invitation(keystoreId, invitationId)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get invitation")
	}

	// TODO: this should be LOCAL keystoreId, while the function only receives remote one
	keystoreKey, err := r.keystores.KeystoreKey(localKeystoreId)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get keystore key")
	}

	if inv.Status != "accepted" || inv.InviteeKey == nil {
		return nil, errors.New("invitation has not been accepted")
	}

	asymKey, err := curve25519.X25519(r.privateKey, inv.InviteeKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create asymmetric key")
	}

	// encrypt the keystore key with the asymmetric key
	encryptedKeystoreKey, err := crypto.AES256CBC_Encrypt(asymKey, keystoreKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt keystore key")
	}

	inv, err = r.client.FinalizeInvitation(keystoreId, invitationId, encryptedKeystoreKey)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to finalize invitation")
	}

	return inv, nil
}

func InvitationToJSON(inv invitation.Invitation) generated.Invitation {
	return generated.Invitation{
		Id:          inv.Id,
		InviterId:   inv.InviterId,
		KeystoreId:  inv.KeystoreId,
		InviteeId:   inv.InviteeId,
		InviterKey:  inv.InviterKey,
		InviteeKey:  inv.InviteeKey,
		KeystoreKey: inv.KeystoreKey,
		Status:      inv.Status.String(),
		CreatedAt:   inv.CreatedAt.Unix(),
		UpdatedAt:   inv.UpdatedAt.Unix(),
	}
}

func InvitationFromJSON(gen generated.Invitation) (invitation.Invitation, error) {
	status, err := invitation.StatusFromString(gen.Status)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "invalid status")
	}

	return invitation.Invitation{
		Id:          gen.Id,
		InviterId:   gen.InviterId,
		KeystoreId:  gen.KeystoreId,
		InviteeId:   gen.InviteeId,
		InviterKey:  gen.InviterKey,
		InviteeKey:  gen.InviteeKey,
		KeystoreKey: gen.KeystoreKey,
		Status:      status,
		CreatedAt:   timestamp.Parse(gen.CreatedAt),
		UpdatedAt:   timestamp.Parse(gen.UpdatedAt),
	}, nil
}
