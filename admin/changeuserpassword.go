package admin

import (
	"github.com/getdiskette/diskette/collections"
	"github.com/getdiskette/diskette/util"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"labix.org/v2/mgo/bson"

	"github.com/labstack/echo"
)

// http POST localhost:5025/admin/change-user-password userId=<user_id> newPassword="123" X-Diskette-Session-Token:<session_token>
func (service *serviceImpl) ChangeUserPassword(c *echo.Context) error {
	var request struct {
		UserId      string `json:"userId"`
		NewPassword string `json:"newPassword"`
	}
	c.Bind(&request)

	if request.UserId == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'userId'")))
	}

	if request.NewPassword == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'newPassword'")))
	}

	var userDoc collections.UserDocument
	err := service.userCollection.FindId(bson.ObjectIdHex(request.UserId)).One(&userDoc)
	if err != nil {
		err = errors.New("No such user.")
		c.JSON(http.StatusNotFound, util.CreateErrResponse(err))
		return err
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	err = service.userCollection.UpdateId(
		bson.ObjectIdHex(request.UserId),
		bson.M{
			"$set": bson.M{
				"hashedPass": hashedPass,
			},
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(nil))
}
