package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"labix.org/v2/mgo"
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

type Service interface {
	Signup(email, password, lang string) (confirmationToken string, err error)
}

type serviceImpl struct {
	db *mgo.Database
}

func (self serviceImpl) Signup(email, password, language string) (confirmationToken string, err error) {
	confirmationKey, err := self.createUser(email, password, lang, false)
	if err != nil {
		return
	}

}

func (self authImpl) createUser(email, password, lang string, isConfirmed bool) (confirmationKey string, err error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	return
}
