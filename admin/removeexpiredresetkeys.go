package admin

import (
	"diskette/collections"
	"diskette/util"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"labix.org/v2/mgo/bson"
)

// http POST localhost:5025/admin/remove-expired-reset-keys X-Diskette-Session-Token:<session_token>
func (service *serviceImpl) RemoveExpiredResetKeys(c *echo.Context) error {

	d, err := time.ParseDuration("12h")
	if err != nil {
		return err
	}

	objectIds := []bson.ObjectId{}
	iter := service.userCollection.Find(nil).Iter()
	var userDoc collections.UserDocument
	for iter.Next(&userDoc) {
		hasResetKey := userDoc.ResetKey != ""
		hasResetKeyExpired := userDoc.RequestedResetAt.Before(time.Now().Add(-1 * d))
		if hasResetKey && hasResetKeyExpired {
			objectIds = append(objectIds, userDoc.Id)
		}
	}
	if err := iter.Close(); err != nil {
		return err
	}

	info, err := service.userCollection.UpdateAll(
		bson.M{
			"_id": bson.M{
				"$in": objectIds,
			},
		},
		bson.M{
			"$set": bson.M{
				"resetKey": "",
			},
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(info.Updated))
}
