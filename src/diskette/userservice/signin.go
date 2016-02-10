package userservice

import (
	"diskette/collections"
	"diskette/tokens"
	"diskette/util"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"labix.org/v2/mgo/bson"
)

// http POST localhost:5025/public/signin email=joe.doe@gmail.com password=abc
func (self impl) Signin(c *echo.Context) error {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.Bind(&request)

	if request.Email == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'email'")))
	}

	if request.Password == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'password'")))
	}

	var userDoc collections.UserDocument
	err := self.userCollection.Find(bson.M{"email": request.Email}).One(&userDoc)
	if err != nil {
		return c.JSON(http.StatusNotFound, util.CreateErrResponse(errors.New("The user doesn't exist.")))
	}

	if userDoc.ConfirmedAt.Before(userDoc.CreatedAt) {
		return c.JSON(http.StatusUnauthorized, util.CreateErrResponse(errors.New("The user has not confirmed the account.")))
	}

	if userDoc.IsSuspended {
		return c.JSON(http.StatusUnauthorized, util.CreateErrResponse(errors.New("The user is suspended.")))
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userDoc.HashedPass), []byte(request.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, util.CreateErrResponse(errors.New("The password didn't match.")))
	}

	token := tokens.SessionToken{
		UserId:    userDoc.Id.Hex(),
		CreatedAt: time.Now(),
	}

	tokenStr, err := token.ToString(self.jwtKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(bson.M{"sessionToken": tokenStr}))
}
