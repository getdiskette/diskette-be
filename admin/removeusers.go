package admin

import (
	"errors"
	"github.com/getdiskette/diskette/util"
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

// http DELETE localhost:5025/admin/remove-users userIds:='["56bf19d65a1d18b704000001", "56be731d5a1d18accd000001"]' X-Diskette-Session-Token:<session_token>
func (service *serviceImpl) RemoveUsers(c echo.Context) error {
	var request struct {
		UserIds []string `json:"userIds"`
	}
	c.Bind(&request)

	if request.UserIds == nil {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'userIds'")))
	}

	objectIds := []bson.ObjectId{}
	for _, userId := range request.UserIds {
		objectIds = append(objectIds, bson.ObjectIdHex(userId))
	}

	info, err := service.userCollection.RemoveAll(bson.M{
		"_id": bson.M{
			"$in": objectIds,
		},
	})
	if err != nil {
		err = errors.New("No such user.")
		c.JSON(http.StatusNotFound, util.CreateErrResponse(err))
		return err
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(info.Removed))
}
