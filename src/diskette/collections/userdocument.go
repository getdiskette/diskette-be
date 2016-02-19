package collections

import (
	"time"

	"labix.org/v2/mgo/bson"
)

const UserCollectionName = "user"

type UserDocument struct {
	Id               bson.ObjectId          `bson:"_id"`
	Email            string                 `bson:"email"`
	HashedPass       []byte                 `bson:"hashedPass"`
	Roles            []string               `bson:"roles"`
	CreatedAt        time.Time              `bson:"createdAt"`
	ConfirmationKey  string                 `bson:"confirmationKey"`
	IsConfirmed      bool                   `bson:"isConfirmed"`
	ResetKey         string                 `bson:"resetKey"`
	RequestedResetAt time.Time              `bson:"requestedResetAt"`
	IsSuspended      bool                   `bson:"isSuspended"`
	SignedOutAt      time.Time              `bson:"signedOutAt"`
	Profile          map[string]interface{} `bson:"profile"`
}
