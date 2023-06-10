package mongorepo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"myst/src/server/application"
	"myst/src/server/application/domain/invitation"
)

func (r *Repository) CreateInvitation(inv invitation.Invitation) (invitation.Invitation, error) {
	collection := r.mdb.Database(r.database).Collection("invitations")
	ctx := context.Background()

	bsonInv := InvitationToBSON(inv)

	_, err := collection.InsertOne(ctx, bsonInv)
	if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to insert invitation")
	}

	inv, err = InvitationFromBSON(bsonInv)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to convert invitation")
	}

	return inv, nil
}

func (r *Repository) Invitation(id string) (invitation.Invitation, error) {
	collection := r.mdb.Database(r.database).Collection("invitations")
	ctx := context.Background()

	res := collection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}})
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return invitation.Invitation{}, application.ErrInvitationNotFound
	} else if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to find invitation")
	}

	var bsonInv Invitation
	err = res.Decode(&bsonInv)
	if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to decode invitation")
	}

	return InvitationFromBSON(bsonInv)
}

func (r *Repository) DeleteInvitation(id string) error {
	collection := r.mdb.Database(r.database).Collection("invitations")
	ctx := context.Background()

	res := collection.FindOneAndDelete(ctx, bson.D{primitive.E{Key: "_id", Value: id}})
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return application.ErrInvitationNotFound
	} else if err != nil {
		return errors.Wrap(err, "failed to delete invitation")
	}

	return nil
}

func (r *Repository) Invitations() ([]invitation.Invitation, error) {
	collection := r.mdb.Database(r.database).Collection("invitations")
	ctx := context.Background()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find invitations")
	}
	defer cur.Close(ctx)

	var bsonInv Invitation
	invitations := []invitation.Invitation{}
	for cur.Next(ctx) {
		err := cur.Decode(&bsonInv)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode invitation")
		}

		inv, err := InvitationFromBSON(bsonInv)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert invitation")
		}

		invitations = append(invitations, inv)
	}
	if err := cur.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate over invitations")
	}

	return invitations, nil
}

func (r *Repository) UpdateInvitation(inv invitation.Invitation) (invitation.Invitation, error) {
	collection := r.mdb.Database(r.database).Collection("invitations")
	ctx := context.Background()

	bsonInv := InvitationToBSON(inv)
	res := collection.FindOneAndReplace(ctx, bson.D{primitive.E{Key: "_id", Value: bsonInv.Id}}, bsonInv)
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return invitation.Invitation{}, application.ErrInvitationNotFound
	} else if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to update invitation")
	}

	return InvitationFromBSON(bsonInv)
}
