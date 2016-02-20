package admin

import (
	"github.com/getdiskette/diskette/util"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"labix.org/v2/mgo/bson"
)

// http POST localhost:5025/admin/signout-users userIds:='["56bf19d65a1d18b704000001", "56be731d5a1d18accd000001"]' X-Diskette-Session-Token:<session_token>
func (service *serviceImpl) SignoutUsers(c *echo.Context) error {
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

	info, err := service.userCollection.UpdateAll(
		bson.M{
			"_id": bson.M{
				"$in": objectIds,
			},
		},
		bson.M{
			"$set": bson.M{
				"signedOutAt": time.Now(),
			},
		},
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(info.Updated))
}
