package session

import (
	"errors"
	"github.com/getdiskette/diskette/collections"
	"github.com/getdiskette/diskette/tokens"
	"github.com/getdiskette/diskette/util"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

// http POST localhost:5025/session/change-email?st=<session_token> password=<password> newEmail=<newEmail>
func (service *serviceImpl) ChangeEmail(c echo.Context) error {
	sessionToken := c.Get("sessionToken").(tokens.SessionToken)
	userDoc := c.Get("userDoc").(collections.UserDocument)

	var request struct {
		Password string `json:"password"`
		NewEmail string `json:"newEmail"`
	}
	c.Bind(&request)

	if request.Password == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'password'")))
	}

	if request.NewEmail == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'newEmail'")))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDoc.HashedPass), []byte(request.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, util.CreateErrResponse(errors.New("The password didn't match.")))
	}

	err := service.userCollection.UpdateId(
		bson.ObjectIdHex(sessionToken.UserId),
		bson.M{
			"$set": bson.M{
				"email": request.NewEmail,
			},
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(nil))
}
