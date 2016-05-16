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

// http POST localhost:5025/session/change-password?st=<session_token> oldPassword=<old_password> newPassword=<new_password>
func (service *serviceImpl) ChangePassword(c echo.Context) error {
	sessionToken := c.Get("sessionToken").(tokens.SessionToken)
	userDoc := c.Get("userDoc").(collections.UserDocument)

	var request struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	c.Bind(&request)

	if request.OldPassword == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'oldPassword'")))
	}

	if request.NewPassword == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'newPassword'")))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDoc.HashedPass), []byte(request.OldPassword)); err != nil {
		return c.JSON(http.StatusUnauthorized, util.CreateErrResponse(errors.New("The old password didn't match.")))
	}

	newHashedPass, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	err = service.userCollection.UpdateId(
		bson.ObjectIdHex(sessionToken.UserId),
		bson.M{
			"$set": bson.M{
				"hashedPass": newHashedPass,
			},
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(nil))
}
