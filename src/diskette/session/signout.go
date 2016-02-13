package session

import (
	"diskette/tokens"
	"diskette/util"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"labix.org/v2/mgo/bson"
)

// http POST localhost:5025/session/signout?st=<session_token>
func (service *serviceImpl) Signout(c *echo.Context) error {
	sessionToken := c.Get("sessionToken").(tokens.SessionToken)

	err := service.userCollection.UpdateId(
		bson.ObjectIdHex(sessionToken.UserId),
		bson.M{
			"$set": bson.M{
				"signedOutAt": time.Now(),
			},
		},
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(nil))
}
