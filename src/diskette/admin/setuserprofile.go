package admin

import (
	"diskette/collections"
	"diskette/util"
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"labix.org/v2/mgo/bson"
)

// http POST localhost:5025/admin/set-user-profile userId=<user_id> newProfile:='{"name": "Joe Doe Jr."}' X-Diskette-Session-Token:<session_token>
func (service *serviceImpl) SetUserProfile(c *echo.Context) error {
	var request struct {
		UserId     string                 `json:"userId"`
		NewProfile map[string]interface{} `json:"newProfile"`
	}
	c.Bind(&request)

	if request.UserId == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'userId'")))
	}

	if request.NewProfile == nil {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'newProfile'")))
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
				"profile": request.NewProfile,
			},
		},
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(nil))
}
