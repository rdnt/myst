package mongorepo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"myst/src/server/application/domain/invitation"
)

type Invitation struct {
	Id                   string    `bson:"_id"`
	KeystoreId           string    `bson:"keystoreId"`
	InviterId            string    `bson:"inviterId"`
	InviteeId            string    `bson:"inviteeId"`
	EncryptedKeystoreKey []byte    `bson:"encryptedKeystoreKey"`
	Status               string    `bson:"status"`
	CreatedAt            time.Time `bson:"createdAt"`
	UpdatedAt            time.Time `bson:"updatedAt"`
	AcceptedAt           time.Time `bson:"acceptedAt"`
	DeclinedAt           time.Time `bson:"declinedAt"`
	DeletedAt            time.Time `bson:"deletedAt"`
}

func InvitationToBSON(inv invitation.Invitation) Invitation {
	return Invitation{
		Id:                   inv.Id,
		KeystoreId:           inv.KeystoreId,
		InviterId:            inv.InviterId,
		InviteeId:            inv.InviteeId,
		EncryptedKeystoreKey: inv.EncryptedKeystoreKey,
		Status:               inv.Status.String(),
		CreatedAt:            inv.CreatedAt,
		UpdatedAt:            inv.UpdatedAt,
		AcceptedAt:           inv.AcceptedAt,
		DeclinedAt:           inv.DeclinedAt,
		DeletedAt:            inv.DeletedAt,
	}
}

func InvitationFromBSON(inv Invitation) (invitation.Invitation, error) {
	stat, err := invitation.StatusFromString(inv.Status)
	if err != nil {
		return invitation.Invitation{}, err
	}

	return invitation.Invitation{
		Id:                   inv.Id,
		KeystoreId:           inv.KeystoreId,
		InviterId:            inv.InviterId,
		InviteeId:            inv.InviteeId,
		EncryptedKeystoreKey: inv.EncryptedKeystoreKey,
		Status:               stat,
		CreatedAt:            inv.CreatedAt,
		UpdatedAt:            inv.UpdatedAt,
		AcceptedAt:           inv.AcceptedAt,
		DeclinedAt:           inv.DeclinedAt,
		DeletedAt:            inv.DeletedAt,
	}, nil
}

func (r *Repository) CreateInvitation(inv invitation.Invitation) (invitation.Invitation, error) {
	collection := r.db.Database("myst").Collection("invitations")

	bsonInv := InvitationToBSON(inv)

	_, err := collection.InsertOne(context.Background(), bsonInv)
	if err != nil {
		return invitation.Invitation{}, err
	}

	inv, err = InvitationFromBSON(bsonInv)
	if err != nil {
		return invitation.Invitation{}, err
	}

	return inv, nil
}

func (r *Repository) Invitation(id string) (invitation.Invitation, error) {
	collection := r.db.Database("myst").Collection("invitations")

	res := collection.FindOne(context.Background(), bson.D{{"_id", id}})
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return invitation.Invitation{}, invitation.ErrNotFound
	} else if err != nil {
		return invitation.Invitation{}, err
	}

	var bsonInv Invitation
	err = res.Decode(&bsonInv)
	if err != nil {
		return invitation.Invitation{}, err
	}

	return InvitationFromBSON(bsonInv)
}

func (r *Repository) Invitations() ([]invitation.Invitation, error) {
	collection := r.db.Database("myst").Collection("invitations")

	ctx := context.Background()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var bsonInv Invitation
	invitations := []invitation.Invitation{}
	for cur.Next(ctx) {
		err := cur.Decode(&bsonInv)
		if err != nil {
			return nil, err
		}

		inv, err := InvitationFromBSON(bsonInv)
		if err != nil {
			return nil, err
		}

		invitations = append(invitations, inv)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return invitations, nil
}

func (r *Repository) UpdateInvitation(inv invitation.Invitation) (invitation.Invitation, error) {
	collection := r.db.Database("myst").Collection("invitations")

	ctx := context.Background()

	bsonInv := InvitationToBSON(inv)

	res := collection.FindOneAndReplace(ctx, bson.D{{"_id", bsonInv.Id}}, bsonInv)
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return invitation.Invitation{}, invitation.ErrNotFound
	} else if err != nil {
		return invitation.Invitation{}, err
	}

	return InvitationFromBSON(bsonInv)
}

func (r *Repository) UserInvitations(userId string) ([]invitation.Invitation, error) {
	allInvs, err := r.Invitations()
	if err != nil {
		return nil, err
	}

	invs := []invitation.Invitation{}

	for _, inv := range allInvs {
		if inv.InviterId == userId {
			invs = append(invs, inv)
		}

		if inv.InviteeId == userId {
			invs = append(invs, inv)
		}
	}

	return invs, nil
}

func (r *Repository) UserInvitation(userId, invitationId string) (invitation.Invitation, error) {
	allInvs, err := r.Invitations()
	if err != nil {
		return invitation.Invitation{}, err
	}

	for _, inv := range allInvs {
		if inv.InviterId == userId && inv.Id == invitationId {
			return inv, nil
		}

		if inv.InviteeId == userId && inv.Id == invitationId {
			return inv, nil
		}
	}

	return invitation.Invitation{}, invitation.ErrNotFound
}
