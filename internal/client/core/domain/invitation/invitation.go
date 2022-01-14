package invitation

import (
	"errors"
	"fmt"

	"myst/pkg/logger"
	"myst/pkg/uuid"
)

var (
	ErrAccepted    = errors.New("invitation already accepted")
	ErrNotAccepted = errors.New("invitation not accepted")
	ErrFinalized   = errors.New("invitation already finalized")
)

type Invitation struct {
	id string
	//inviter     *user.User
	//keystore    *keystore.Keystore
	//invitee     *user.User
	//inviterKey  []byte
	//inviteeKey  []byte
	//keystoreKey []byte // encrypted
	//accepted    bool
	//finalized   bool
	//createdAt timestamp.Timestamp
	//updatedAt timestamp.Timestamp
}

func (i *Invitation) Id() string {
	return i.id
}

//func (i *Invitation) Inviter() *user.User {
//	return i.inviter
//}
//
//func (i *Invitation) Keystore() *keystore.Keystore {
//	return i.keystore
//}
//
//func (i *Invitation) Invitee() *user.User {
//	return i.invitee
//}
//
//func (i *Invitation) InviterKey() []byte {
//	return i.inviterKey
//}
//
//func (i *Invitation) InviteeKey() []byte {
//	return i.inviteeKey
//}
//
//func (i *Invitation) KeystoreKey() []byte {
//	return i.keystoreKey
//}
//
//func (i *Invitation) Accepted() bool {
//	return i.accepted
//}
//
//func (i *Invitation) Finalized() bool {
//	return i.finalized
//}
//
//func (i *Invitation) CreatedAt() timestamp.Timestamp {
//	return i.createdAt
//}
//
//func (i *Invitation) UpdatedAt() timestamp.Timestamp {
//	return i.updatedAt
//}

func (i *Invitation) String() string {
	//return fmt.Sprintln(i.id, i.inviter.Id(), i.keystore.Id(), i.invitee.Id(), i.accepted, i.finalized)
	return fmt.Sprintln("invitation", i.id)
}

//func (i *Invitation) Accept(inviteeKey []byte) error {
//	if i.accepted {
//		return ErrAccepted
//	}
//
//	i.inviteeKey = inviteeKey
//	i.accepted = true
//
//	return nil
//}
//
//func (i *Invitation) Finalize(keystoreKey []byte) error {
//	if !i.accepted {
//		return ErrNotAccepted
//	}
//
//	if i.finalized {
//		return ErrFinalized
//	}
//
//	i.keystoreKey = keystoreKey
//	i.finalized = true
//
//	return nil
//}

func New(opts ...Option) (*Invitation, error) {
	k := &Invitation{
		id: uuid.New().String(),
		//createdAt: timestamp.New(),
		//updatedAt: timestamp.New(),
	}

	for _, opt := range opts {
		err := opt(k)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	return k, nil
}
