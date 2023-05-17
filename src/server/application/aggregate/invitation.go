package aggregate

import (
	"myst/src/server/application/domain/invitation"
	"myst/src/server/application/domain/user"
)

type Invitation struct {
	invitation.Invitation
	inviter user.User
	invitee user.User
}
