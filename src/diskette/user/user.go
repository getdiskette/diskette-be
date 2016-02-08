package user

import (
	"time"

	"diskette/vendor/labix.org/v2/mgo"
)

type UserDoc struct {
	Id               string    `json:"_id"`
	Roles            []string  `json:"roles"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	HashedPass       string    `json:"hashedPass"`
	Language         string    `json:"lang"`
	CreatedAt        time.Time `json:"createdAt"`
	ConfirmedAt      time.Time `json:"confirmedAt"`
	ResetKey         string    `json:"resetKey"`
	RequestedResetAt string    `json:"requestedResetAt"`
	IsSuspended      bool      `json:"isSuspended"`
	Sessions         map[string]struct {
		CreatedAt time.Time `json:"createdAt"`
	} `json:"sessions"`
}

type UserService interface {
	Signup(email, password, lang string) (confirmationToken string, err error)
}

type impl struct {
	db *mgo.Database
}

func (self impl) Signup(email, password, language string) (confirmationToken string, err error) {
	return "", nil
}
