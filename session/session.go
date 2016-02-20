package session

import (
	"github.com/labstack/echo"
	"labix.org/v2/mgo"
)

type Service interface {
	Signout(c *echo.Context) error
	ChangePassword(c *echo.Context) error
	ChangeEmail(c *echo.Context) error
	SetProfile(c *echo.Context) error
}

type serviceImpl struct {
	userCollection *mgo.Collection
	jwtKey         []byte
}

func NewService(userCollection *mgo.Collection, jwtKey []byte) Service {
	return &serviceImpl{userCollection, jwtKey}
}
