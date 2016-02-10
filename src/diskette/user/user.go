package user

import (
	"github.com/labstack/echo"
	"labix.org/v2/mgo"
)

type AuthenticationService interface {
	// Publicly Available
	Signup(c *echo.Context) error
	ConfirmSignup(c *echo.Context) error
	Signin(c *echo.Context) error
	ForgotPassword(c *echo.Context) error
	ResetPassword(c *echo.Context) error
	// Authorized User
	Signout(c *echo.Context) error
}

type impl struct {
	userCollection *mgo.Collection
	jwtKey         []byte
}

func NewAuthenticationService(userCollection *mgo.Collection, jwtKey []byte) AuthenticationService {
	return &impl{userCollection, jwtKey}
}
