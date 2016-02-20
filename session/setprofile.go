package session

import (
	"github.com/getdiskette/diskette/tokens"
	"github.com/getdiskette/diskette/util"
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"labix.org/v2/mgo/bson"
)

// http POST localhost:5025/session/set-profile?st=<session_token> profile:='{"profession": "Software Developer"}'
func (service *serviceImpl) SetProfile(c *echo.Context) error {
	sessionToken := c.Get("sessionToken").(tokens.SessionToken)

	var request struct {
		Profile map[string]interface{} `json:"profile"`
	}
	c.Bind(&request)

	if request.Profile == nil {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'profile'")))
	}

	err := service.userCollection.UpdateId(
		bson.ObjectIdHex(sessionToken.UserId),
		bson.M{
			"$set": bson.M{
				"profile": request.Profile,
			},
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(nil))
}
