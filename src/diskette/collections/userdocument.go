package collections

import (
	"time"

	"labix.org/v2/mgo/bson"
)

const USER_COLLECTION_NAME = "user"

type UserDocument struct {
	Id                   bson.ObjectId `bson:"_id"`
	Name                 string        `bson:"name"`
	Email                string        `bson:"email"`
	HashedPass           []byte        `bson:"hashedPass"`
	Language             string        `bson:"language"`
	Roles                []string      `bson:"roles"`
	CreatedAt            time.Time     `bson:"createdAt"`
	ConfirmationKey      string        `bson:"confirmationKey"`
	ConfirmedAt          time.Time     `bson:"confirmedAt"`
	ResetKey             string        `bson:"resetKey"`
	RequestedResetAt     time.Time     `bson:"requestedResetAt"`
	IsSuspended          bool          `bson:"isSuspended"`
	RejectSessionsBefore time.Time     `bson:"rejectSessionsBefore"`
}
