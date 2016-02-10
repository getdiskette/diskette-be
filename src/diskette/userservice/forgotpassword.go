package userservice

import (
	"diskette/tokens"
	"diskette/util"
	"errors"
	"net/http"
	"time"

	"github.com/satori/go.uuid"

	"github.com/labstack/echo"
	"labix.org/v2/mgo/bson"
)

// http POST localhost:5025/public/forgot-password email=joe.doe@gmail.com
func (self impl) ForgotPassword(c *echo.Context) error {
	var request struct {
		Email string `json:"email"`
	}
	c.Bind(&request)

	if request.Email == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'email'")))
	}

	resetKey := uuid.NewV4().String()

	err := self.userCollection.Update(
		bson.M{"email": request.Email},
		bson.M{
			"$set": bson.M{
				"resetKey":         resetKey,
				"requestedResetAt": time.Now(),
			},
		},
	)
	if err != nil {
		return c.JSON(http.StatusNotFound, util.CreateErrResponse(errors.New("The user doesn't exist.")))
	}

	token := tokens.ResetToken{Key: resetKey}
	tokenStr, err := token.ToString(self.jwtKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(bson.M{"ResetToken": tokenStr}))
}
