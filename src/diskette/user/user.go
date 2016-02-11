package user

import (
	"github.com/labstack/echo"
	"labix.org/v2/mgo"
)

type UserService interface {
	// Publicly Available
	Signup(c *echo.Context) error
	ConfirmSignup(c *echo.Context) error
	Signin(c *echo.Context) error
	ForgotPassword(c *echo.Context) error
	ResetPassword(c *echo.Context) error
	// Authorized User
	Signout(c *echo.Context) error
}

type serviceImpl struct {
	userCollection *mgo.Collection
	jwtKey         []byte
}

func NewUserService(userCollection *mgo.Collection, jwtKey []byte) UserService {
	return &serviceImpl{userCollection, jwtKey}
}
