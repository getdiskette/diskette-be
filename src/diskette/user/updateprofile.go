package user

import (
	"diskette/tokens"
	"diskette/util"
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"labix.org/v2/mgo/bson"
)

// http POST localhost:5025/private/update-profile?st=<session_token> profile:='{"profession": "Software Developer"}'
func (service *serviceImpl) UpdateProfile(c *echo.Context) error {
	sessionToken := c.Get("sessionToken").(tokens.SessionToken)

	var request struct {
		Profile map[string]interface{} `json:"profile"`
	}
	c.Bind(&request)

	if request.Profile == nil {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'profile'")))
	}

	err := service.userCollection.UpdateId(
		sessionToken.UserId,
		bson.M{
			"$set": bson.M{
				"profile": request.Profile,
			},
		},
	)
	if err != nil {
		return c.JSON(http.StatusNotFound, util.CreateErrResponse(errors.New("The user doesn't exist.")))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(nil))
}
