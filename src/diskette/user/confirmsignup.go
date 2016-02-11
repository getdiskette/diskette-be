package user

import (
	"diskette/util"
	"errors"
	"net/http"
	"time"

	"diskette/tokens"

	"github.com/labstack/echo"
	"labix.org/v2/mgo/bson"
)

// http POST localhost:5025/user/confirm token=<confirmation_token>
func (service *impl) ConfirmSignup(c *echo.Context) error {
	var request struct {
		Token string `json:"token"`
	}
	c.Bind(&request)

	if request.Token == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'token'")))
	}

	token, err := tokens.ParseConfirmationToken(service.jwtKey, request.Token)
	if err != nil || token.Key == "" {
		return c.JSON(http.StatusForbidden, util.CreateErrResponse(err))
	}

	err = service.userCollection.Update(
		bson.M{"confirmationKey": token.Key},
		bson.M{
			"$set": bson.M{
				"confirmedAt": time.Now(),
			},
		},
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(nil))
}
