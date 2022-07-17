package invitation

type Repository interface {
	CreateInvitation(inv Invitation) (Invitation, error)
	Invitation(id string) (Invitation, error)
	UpdateInvitation(inv Invitation) error
	Invitations() (map[string]Invitation, error)
	DeleteInvitation(id string) error
}

type Service interface {
	CreateInvitation(keystoreId string, inviteeUsername string) (Invitation, error)
	AcceptInvitation(id string) (Invitation, error)
	DeclineOrCancelInvitation(id string) (Invitation, error)
	FinalizeInvitation(id string) (Invitation, error)
	Invitations() (map[string]Invitation, error)
	Invitation(id string) (Invitation, error)
}
