package admin

import (
	"github.com/getdiskette/diskette/collections"
	"github.com/getdiskette/diskette/util"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"labix.org/v2/mgo/bson"
)

// http DELETE localhost:5025/admin/remove-unconfirmed-users X-Diskette-Session-Token:<session_token>
func (service *serviceImpl) RemoveUnconfirmedUsers(c *echo.Context) error {

	d, err := time.ParseDuration("12h")
	if err != nil {
		return err
	}

	objectIds := []bson.ObjectId{}
	iter := service.userCollection.Find(nil).Iter()
	var userDoc collections.UserDocument
	for iter.Next(&userDoc) {
		isAccountOld := userDoc.CreatedAt.Before(time.Now().Add(-1 * d))
		if isAccountOld && !userDoc.IsConfirmed {
			objectIds = append(objectIds, userDoc.Id)
		}
	}
	if err := iter.Close(); err != nil {
		return err
	}

	info, err := service.userCollection.RemoveAll(bson.M{
		"_id": bson.M{
			"$in": objectIds,
		},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(info.Updated))
}
