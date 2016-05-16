package user

import (
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
)

type Service interface {
	Signup(c echo.Context) error
	ConfirmSignup(c echo.Context) error
	Signin(c echo.Context) error
	ForgotPassword(c echo.Context) error
	ResetPassword(c echo.Context) error
}

type serviceImpl struct {
	userCollection *mgo.Collection
	jwtKey         []byte
}

func NewService(userCollection *mgo.Collection, jwtKey []byte) Service {
	return &serviceImpl{userCollection, jwtKey}
}
