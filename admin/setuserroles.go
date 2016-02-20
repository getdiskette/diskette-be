package admin

import (
	"github.com/getdiskette/diskette/collections"
	"github.com/getdiskette/diskette/util"
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"labix.org/v2/mgo/bson"
)

// http POST localhost:5025/admin/set-user-roles userId=<user_id> newRoles:='["admin"]' X-Diskette-Session-Token:<session_token>
func (service *serviceImpl) SetUserRoles(c *echo.Context) error {
	var request struct {
		UserId   string   `json:"userId"`
		NewRoles []string `json:"newRoles"`
	}
	c.Bind(&request)

	if request.UserId == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'userId'")))
	}

	if request.NewRoles == nil {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'newRoles'")))
	}

	var userDoc collections.UserDocument
	err := service.userCollection.FindId(bson.ObjectIdHex(request.UserId)).One(&userDoc)
	if err != nil {
		err = errors.New("No such user.")
		c.JSON(http.StatusNotFound, util.CreateErrResponse(err))
		return err
	}

	err = service.userCollection.UpdateId(
		bson.ObjectIdHex(request.UserId),
		bson.M{
			"$set": bson.M{
				"roles": request.NewRoles,
			},
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(nil))
}
