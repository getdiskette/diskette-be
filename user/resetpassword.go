package user

import (
	"errors"
	"github.com/getdiskette/diskette/tokens"
	"github.com/getdiskette/diskette/util"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// http POST localhost:5025/user/reset-password token=<reset_token> password=123
func (service *serviceImpl) ResetPassword(c echo.Context) error {
	var request struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	c.Bind(&request)

	if request.Token == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'token'")))
	}

	if request.Password == "" {
		return c.JSON(http.StatusBadRequest, util.CreateErrResponse(errors.New("Missing parameter 'password'")))
	}

	token, err := tokens.ParseResetToken(service.jwtKey, request.Token)
	if err != nil || token.Key == "" {
		return c.JSON(http.StatusForbidden, util.CreateErrResponse(err))
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.CreateErrResponse(err))
	}

	err = service.userCollection.Update(
		bson.M{"resetKey": token.Key},
		bson.M{
			"$set": bson.M{
				"resetKey":         "",
				"requestedResetAt": time.Unix(0, 0),
				"hashedPass":       hashedPass,
			},
		},
	)
	if err != nil {
		return c.JSON(http.StatusNotFound, util.CreateErrResponse(errors.New("The token doesn't exist.")))
	}

	return c.JSON(http.StatusOK, util.CreateOkResponse(nil))
}
