package user

import "time"

type UserDoc struct {
	Name             string    `bson:"name"`
	Email            string    `bson:"email"`
	HashedPass       []byte    `bson:"hashedPass"`
	Language         string    `bson:"lang"`
	Roles            []string  `bson:"roles"`
	CreatedAt        time.Time `bson:"createdAt"`
	ConfirmationKey  string    `bson:"confirmationKey"`
	ConfirmedAt      time.Time `bson:"confirmedAt"`
	ResetKey         string    `bson:"resetKey"`
	RequestedResetAt time.Time `bson:"requestedResetAt"`
	IsSuspended      bool      `bson:"isSuspended"`
	Sessions         map[string]struct {
		CreatedAt time.Time `bson:"createdAt"`
	} `bson:"sessions"`
}
