package user

import (
	"diskette/tokens"
	"diskette/util"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"labix.org/v2/mgo/bson"
)

// http POST localhost:5025/private/signout?st=<session_token>
func (service *serviceImpl) Signout(c *echo.Context) error {
	sessionToken := c.Get("sessionToken").(tokens.SessionToken)

	err := service.userCollection.UpdateId(
		sessionToken.UserId,
		bson.M{
			"$set": bson.M{
				"signedOutAt": time.Now(),
			},
		},
	)

	if err != nil {
		return c.JSON(http.StatusNotFound, util.CreateErrResponse(errors.New("The user doesn't exist.")))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(nil))
}
