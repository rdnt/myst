package invitation

type Repository interface {
	CreateInvitation(inv Invitation) (Invitation, error)
	Invitation(id string) (Invitation, error)
	UpdateInvitation(inv Invitation) error
	Invitations() (map[string]Invitation, error)
	DeleteInvitation(id string) error
}

type Service interface {
	CreateKeystoreInvitation(keystoreId, inviterId, inviteeId string, inviterKey []byte) (Invitation, error)
}
